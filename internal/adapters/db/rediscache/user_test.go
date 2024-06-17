package rediscache

import (
    "context"
    "errors"
    "github.com/boldyrevn/mod-example/internal/adapters/db/mocks"
    "github.com/boldyrevn/mod-example/internal/model"
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestDataBase(t *testing.T) {
    var (
        host = "127.0.0.1"
        port = 6379
        ctx  = context.Background()
    )

    dbMock := mocks.NewDB(t)
    client, err := New(ctx, host, port, dbMock)
    assert.NoError(t, err)

    err = client.conn.FlushDB(ctx).Err()
    assert.NoError(t, err)

    user := model.User{
        ID:        "123",
        FirstName: "Test",
        LastName:  "Testov",
        Age:       23,
    }

    dbMock.On("GetUser", ctx, user.ID).Once().Return(user, nil)
    actualUser, err := client.GetUser(ctx, "123")
    dbMock.AssertNumberOfCalls(t, "GetUser", 1)
    assert.Equal(t, user, actualUser)

    actualUser, err = client.GetUser(ctx, "123")
    dbMock.AssertNumberOfCalls(t, "GetUser", 1)
    assert.Equal(t, user, actualUser)

    dbMock.On("DeleteUser", ctx, user.ID).Return(nil)
    err = client.DeleteUser(ctx, user.ID)
    assert.NoError(t, err)

    dbMock.On("GetUser", ctx, user.ID).Return(model.User{}, errors.New("some error"))
    _, err = client.GetUser(ctx, user.ID)
    assert.Error(t, err)
}
