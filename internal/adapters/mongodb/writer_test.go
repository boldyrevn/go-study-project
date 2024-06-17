package mongodb

import (
    "context"
    "encoding/json"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "go.mongodb.org/mongo-driver/mongo/options"
    "testing"
)

func TestDataBase_Write(t *testing.T) {
    ctx := context.Background()

    client, err := New(ctx, "test_db", options.Client().SetHosts([]string{"localhost:27017"}))
    assert.NoError(t, err)

    err = client.conn.Database("test_db").Drop(ctx)
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
}
