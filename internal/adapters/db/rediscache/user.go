package rediscache

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/boldyrevn/mod-example/internal/model"
)

func getUserKey(id string) string {
    return fmt.Sprintf("user:%s", id)
}

func (db *DataBase) setUser(ctx context.Context, user model.User) error {
    rawUser, err := json.Marshal(user)
    if err != nil {
        return err
    }
    return db.conn.Set(ctx, getUserKey(user.ID), string(rawUser), 0).Err()
}

func (db *DataBase) CreateUser(ctx context.Context, user model.User) error {
    if err := db.internalDB.CreateUser(ctx, user); err != nil {
        return err
    }
    return db.setUser(ctx, user)
}

func (db *DataBase) GetUser(ctx context.Context, id string) (model.User, error) {
    var res model.User

    userRaw, err := db.conn.Get(ctx, getUserKey(id)).Bytes()
    if err == nil {
        if err := json.Unmarshal(userRaw, &res); err != nil {
            return model.User{}, err
        }
        return res, nil
    }

    res, err = db.internalDB.GetUser(ctx, id)
    if err != nil {
        return model.User{}, err
    }

    return res, db.setUser(ctx, res)
}

func (db *DataBase) UpdateUser(ctx context.Context, user model.User) error {
    if err := db.internalDB.UpdateUser(ctx, user); err != nil {
        return err
    }
    return db.setUser(ctx, user)
}

func (db *DataBase) DeleteUser(ctx context.Context, id string) error {
    if err := db.internalDB.DeleteUser(ctx, id); err != nil {
        return err
    }
    return db.conn.Del(ctx, getUserKey(id)).Err()
}
