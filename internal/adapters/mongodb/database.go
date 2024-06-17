package mongodb

import (
    "context"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type DataBase struct {
    conn *mongo.Client
    name string
}

func New(ctx context.Context, name string, opts ...*options.ClientOptions) (*DataBase, error) {
    conn, err := mongo.Connect(ctx, opts...)
    if err != nil {
        return nil, err
    }

    if err := conn.Ping(ctx, nil); err != nil {
        return nil, err
    }

    return &DataBase{
        conn: conn,
        name: name,
    }, nil
}
