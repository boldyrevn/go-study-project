package postgres

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/boldyrevn/mod-example/internal/model"
)

func (db *DataBase) CreateUser(ctx context.Context, user model.User) error {
    profile, err := json.Marshal(user.Profile)
    if err != nil {
        return err
    }

    _, err = db.conn.Exec(ctx, createUserQuery, user.ID, user.FirstName, user.LastName, user.Age, string(profile))
    if err != nil {
        return fmt.Errorf("failed to create user: %w", err)
    }

    return nil
}

func (db *DataBase) GetUser(ctx context.Context, id string) (model.User, error) {
    result, err := db.conn.Query(ctx, getUserQuery, id)
    db.conn.Query(context.Background(), getUserQuery, nil, id)
    defer result.Close()
    if err != nil {
        return model.User{}, fmt.Errorf("failed to get user: %w", err)
    }

    hasRows := result.Next()
    if !hasRows {
        return model.User{}, fmt.Errorf("result has no rows")
    }

    user := model.User{}
    err = result.Scan(
        &user.ID,
        &user.FirstName,
        &user.LastName,
        &user.Age,
        &user.Profile,
    )
    if err != nil {
        return model.User{}, fmt.Errorf("failed to scan data: %w", err)
    }

    return user, err
}

func (db *DataBase) DeleteUser(ctx context.Context, id string) error {
    _, err := db.conn.Exec(ctx, deleteUserQuery, id)
    if err != nil {
        return fmt.Errorf("failed to delete user: %w", err)
    }
    return nil
}

func (db *DataBase) UpdateUser(ctx context.Context, user model.User) error {
    rawProfile, err := json.Marshal(user.Profile)
    if err != nil {
        return err
    }

    _, err = db.conn.Exec(ctx, updateUserQuery, user.ID, user.FirstName, user.LastName, user.Age, string(rawProfile))
    if err != nil {
        return fmt.Errorf("failed to update user: %w", err)
    }

    return nil
}
