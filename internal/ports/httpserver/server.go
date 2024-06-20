package httpserver

import (
    "context"
    "errors"
    "fmt"
    "github.com/boldyrevn/mod-example/internal/adapters/db"
    "github.com/boldyrevn/mod-example/internal/adapters/db/postgres"
    "github.com/boldyrevn/mod-example/internal/adapters/db/rediscache"
    "github.com/boldyrevn/mod-example/internal/adapters/mongodb"
    "github.com/boldyrevn/mod-example/internal/ports/httpserver/handlers"
    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/mongo/options"
    "log/slog"
    "net/http"
    "os"
    "strconv"
    "time"
)

type Server struct {
    http         *http.Server
    router       *gin.Engine
    db           db.DB
    logger       *slog.Logger
    mongo        *mongodb.DataBase
    redis        *rediscache.DataBase
    postgres     *postgres.DataBase
    userHandlers *handlers.UserHandlers
}

func (s *Server) InitPostgres() error {
    var (
        user     = os.Getenv("POSTGRES_USER")
        password = os.Getenv("POSTGRES_PASSWORD")
        host     = os.Getenv("POSTGRES_HOST")
        port     = os.Getenv("POSTGRES_PORT")
        dataBase = os.Getenv("POSTGRES_DB")
    )

    connString := fmt.Sprintf(
        "user=%s password=%s host=%s port=%s dbname=%s",
        user, password, host, port, dataBase,
    )
    conn, err := postgres.New(context.Background(), connString)
    if err != nil {
        return err
    }

    s.db = conn
    s.postgres = conn
    return nil
}

func (s *Server) InitMongo() error {
    var (
        mongoDBName = os.Getenv("MONGODB_NAME")
        mongoDBHost = os.Getenv("MONGODB_HOST")
        mongoDBPort = os.Getenv("MONGODB_PORT")
    )
    if mongoDBName == "" {
        return errors.New("mongoDB database is not specified")
    }

    conn, err := mongodb.New(
        context.Background(),
        mongoDBName,
        &options.ClientOptions{
            Hosts: []string{fmt.Sprintf("%s:%s", mongoDBHost, mongoDBPort)},
        },
    )
    if err != nil {
        return err
    }

    s.mongo = conn
    return nil
}

func (s *Server) InitRedis() error {
    var (
        host = os.Getenv("REDIS_HOST")
        port = os.Getenv("REDIS_PORT")
    )
    dbNum, err := strconv.Atoi(os.Getenv("REDIS_DB"))
    if err != nil {
        return err
    }

    if s.db == nil {
        return errors.New("db must be initialized before starting redis")
    }
    conn, err := rediscache.New(context.Background(), host, port, dbNum, s.db)
    if err != nil {
        return err
    }

    s.redis = conn
    s.db = conn
    return nil
}

func (s *Server) InitLogger() {
    logger := slog.New(slog.NewJSONHandler(s.mongo, &slog.HandlerOptions{
        Level: slog.LevelInfo,
    }))
    s.logger = logger
}

func (s *Server) InitHandlers() {
    s.userHandlers = handlers.NewUserHandlers(s.db, s.logger)
}

func (s *Server) InitRouter() {
    s.router = gin.New()
    s.router.Use(
        gin.Recovery(),
        RequestIDMiddleware(),
        LoggingMiddleware(s.logger),
    )

    userRoutes := s.router.Group("/user")
    {
        userRoutes.POST("/create", s.userHandlers.CreateUser)
        userRoutes.GET("/get", s.userHandlers.GetUser)
        userRoutes.PUT("/update", s.userHandlers.UpdateUser)
        userRoutes.DELETE("/delete", s.userHandlers.DeleteUser)
    }

    s.http.Handler = s.router.Handler()
}

func (s *Server) Init() error {
    s.http = &http.Server{
        ReadTimeout:       4 * time.Second,
        ReadHeaderTimeout: time.Second,
        WriteTimeout:      10 * time.Second,
    }

    if err := s.InitPostgres(); err != nil {
        return err
    }
    if err := s.InitMongo(); err != nil {
        return err
    }
    if err := s.InitRedis(); err != nil {
        return err
    }

    s.InitLogger()
    s.InitHandlers()
    s.InitRouter()
    return nil
}

func (s *Server) Listen() error {
    return s.http.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
    return s.http.Shutdown(ctx)
}

func (s *Server) CloseConnections(ctx context.Context) {
    _ = s.mongo.Close(ctx)
    _ = s.redis.Close()
    _ = s.postgres.Close(ctx)
}

func New() *Server {
    return &Server{}
}
