package migrations

import (
    "context"
    "github.com/jackc/pgx/v4/pgxpool"
)

const migrationQueryUp = `
create table if not exists "Users"
(
    id         varchar(64) not null primary key,
    first_name varchar(64) not null,
    last_name  varchar(64) not null,
    age        integer     not null,
    profile    json
);

alter table "Users"
    owner to postgres;
`

const migrationQueryDown = `
drop table "Users"
`

func Up(ctx context.Context, conn *pgxpool.Pool) error {
    _, err := conn.Exec(ctx, migrationQueryUp)
    return err
}

func Down(ctx context.Context, conn *pgxpool.Conn) error {
    _, err := conn.Exec(ctx, migrationQueryDown)
    return err
}
