package main

import (
    "context"
    "errors"
    "github.com/boldyrevn/mod-example/internal/ports/httpserver"
    "github.com/joho/godotenv"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func loadEnv() {
    if err := godotenv.Load(); err != nil {
        panic(err)
    }
}

func main() {
    server := httpserver.New()

    if err := server.Init(); err != nil {
        panic(err)
    }
    log.Println("server initialized")

    go func() {
        if err := server.Listen(); err != nil && !errors.Is(err, http.ErrServerClosed) {
            log.Fatalf("listen: %s\n", err)
        }
    }()

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := server.Shutdown(ctx); err != nil {
        log.Fatal("server shutdown:", err)
    }

    server.CloseConnections(ctx)
    <-ctx.Done()
    log.Println("server exiting")
}
