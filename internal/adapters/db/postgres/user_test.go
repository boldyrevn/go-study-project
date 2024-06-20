package postgres

import (
    "context"
    "fmt"
    "github.com/boldyrevn/mod-example/internal/model"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "os"
    "testing"
    "time"
)

func TestDataBase(t *testing.T) {
    var (
        userDB   = os.Getenv("POSTGRES_USER")
        password = os.Getenv("POSTGRES_PASSWORD")
        host     = os.Getenv("POSTGRES_TEST_HOST")
        port     = os.Getenv("POSTGRES_PORT")
        dataBase = os.Getenv("POSTGRES_DB")
    )

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    connString := fmt.Sprintf(
        "user=%s password=%s host=%s port=%s dbname=%s",
        userDB, password, host, port, dataBase,
    )
    db, err := New(ctx, connString)

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
            Bio:       "my brain is very big",
            Interests: []string{"after parties"},
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
