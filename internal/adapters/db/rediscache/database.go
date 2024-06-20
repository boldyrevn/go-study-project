package rediscache

import (
    "context"
    "fmt"
    "github.com/boldyrevn/mod-example/internal/adapters/db"
    "github.com/redis/go-redis/v9"
)

type DataBase struct {
    internalDB db.DB
    conn       *redis.Client
}

func New(ctx context.Context, host, port string, db int, wrappedDB db.DB) (*DataBase, error) {
    conn := redis.NewClient(&redis.Options{
        Addr: fmt.Sprintf("%s:%s", host, port),
        DB:   db,
    })
    if err := conn.Ping(ctx).Err(); err != nil {
        return nil, err
    }
    return &DataBase{
        internalDB: wrappedDB,
        conn:       conn,
    }, nil
}

func (db *DataBase) Close() error {
    return db.conn.Close()
}
