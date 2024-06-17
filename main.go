package main

import (
    "log/slog"
    "os"
)

type User struct {
    Id   uint64 `json:"id,omitempty"`
    Name string `json:"name,omitempty"`
    Age  uint   `json:"age,omitempty"`
}

func main() {
    logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

    user := &User{
        Id:   123,
        Name: "Alyosha",
        Age:  23,
    }
    logger.Info("user has been extracted successfully", slog.Any("user", user))
}
