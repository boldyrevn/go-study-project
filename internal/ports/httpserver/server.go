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
    "os"
)

type Server struct {
    router       *gin.Engine
    db           db.DB
    logger       *slog.Logger
    mongo        *mongodb.DataBase
    redis        *rediscache.DataBase
    userHandlers *handlers.UserHandlers
}

func (s *Server) InitPostgres(ctx context.Context) error {
    var (
        user     = os.Getenv("POSTGRES_USER")
        password = os.Getenv("POSTGRES_PASSWORD")
        host     = os.Getenv("POSTGRES_HOST")
        port     = os.Getenv("POSTGRES_PORT")
        dataBase = os.Getenv("POSTGRES_DATABASE")
    )

    connString := fmt.Sprintf(
        "user=%s password=%s host=%s port=%s dbname=%s",
        user, password, host, port, dataBase,
    )
    conn, err := postgres.New(ctx, connString)
    if err != nil {
        return err
    }

    s.db = conn
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

func (s *Server) InitRedis(ctx context.Context) error {
    var (
        host = os.Getenv("REDIS_HOST")
        port = os.Getenv("REDIS_PORT")
    )

    conn, err := rediscache.New(ctx, host, port, s.db)
    if err != nil {
        return err
    }

    s.redis = conn
    s.db = conn
    return nil
}

func (s *Server) InitLogger() error {
    logger := slog.New(slog.NewJSONHandler(s.mongo, &slog.HandlerOptions{
        Level: slog.LevelInfo,
    }))
    s.logger = logger
    return nil
}

func (s *Server) InitHandlers() error {
    s.userHandlers = handlers.NewUserHandlers(s.db, s.logger)
    return nil
}

func (s *Server) InitRouter() error {
    s.router.Use(
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
    return nil
}

func New() *Server {
    return &Server{
        router: gin.New(),
    }
}
