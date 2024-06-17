package postgres

import (
    "context"
    "github.com/boldyrevn/mod-example/internal/model"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "testing"
)

func TestDataBase_CreateUser(t *testing.T) {
    ctx := context.Background()
    db, err := New(ctx, "user=postgres password=povar host=localhost port=8080 dbname=go_app")
    require.NoError(t, err)
    defer db.Close(ctx)

    _, err = db.conn.Exec(ctx, `TRUNCATE "Users"`)
    assert.NoError(t, err)

    user := model.User{
        ID:        "124",
        FirstName: "Nikita",
        LastName:  "Boldyrev",
        Age:       20,
        Profile: model.ProfileDescription{
            Bio:       "my dick is very big",
            Interests: []string{"gay parties"},
        },
    }

    err = db.CreateUser(ctx, user)
    assert.NoError(t, err)

    actualUser, err := db.GetUser(ctx, user.ID)
    assert.NoError(t, err)
    assert.Equal(t, user, actualUser)

    user.LastName = "Amogusov"
    err = db.UpdateUser(ctx, user)
    assert.NoError(t, err)

    actualUser, err = db.GetUser(ctx, user.ID)
    assert.NoError(t, err)
    assert.Equal(t, user, actualUser)

    err = db.DeleteUser(ctx, user.ID)
    assert.NoError(t, err)

    _, err = db.GetUser(ctx, user.ID)
    assert.Error(t, err)
    assert.ErrorContains(t, err, "result has no rows")
}
