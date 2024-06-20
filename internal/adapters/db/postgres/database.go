package postgres

import (
    "context"
    "github.com/boldyrevn/mod-example/internal/migrations"
    "github.com/jackc/pgx/v4/pgxpool"
)

type DataBase struct {
    conn *pgxpool.Pool
}

func New(ctx context.Context, connString string) (*DataBase, error) {
    conf, err := pgxpool.ParseConfig(connString)
    if err != nil {
        return nil, err
    }

    conn, err := pgxpool.ConnectConfig(ctx, conf)
    if err != nil {
        return nil, err
    }

    if err := migrations.Up(context.Background(), conn); err != nil {
        return nil, err
    }

    return &DataBase{conn: conn}, nil
}

func (db *DataBase) Close(ctx context.Context) error {
    db.conn.Close()
    return nil
}
