package mongodb

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "go.mongodb.org/mongo-driver/mongo/options"
    "os"
    "testing"
    "time"
)

func TestDataBase(t *testing.T) {
    var (
        mongoDBName = os.Getenv("MONGODB_NAME")
        mongoDBHost = os.Getenv("MONGODB_HOST")
        mongoDBPort = os.Getenv("MONGODB_PORT")
    )

    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    client, err := New(
        ctx,
        fmt.Sprintf("%s_test", mongoDBName),
        options.Client().SetHosts([]string{
            fmt.Sprintf("%s:%s", mongoDBHost, mongoDBPort),
        }))
    require.NoError(t, err)

    someLog := struct {
        Name string `json:"name,omitempty"`
        Time int    `json:"time,omitempty"`
    }{
        Name: "amogus",
        Time: 123,
    }

    rawLog, err := json.Marshal(someLog)
    require.NoError(t, err)

    n, err := client.Write(rawLog)
    assert.NoError(t, err)
    assert.Equal(t, len(rawLog), n)

    err = client.conn.Database("test_db").Drop(ctx)
    require.NoError(t, err)

    err = client.Close(ctx)
}
