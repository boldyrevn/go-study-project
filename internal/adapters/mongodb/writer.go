package mongodb

import (
    "context"
    "encoding/json"
)

func (db *DataBase) Write(p []byte) (int, error) {
    var val map[string]any
    if err := json.Unmarshal(p, &val); err != nil {
        return 0, err
    }

    _, err := db.conn.Database(db.name).Collection("logs").InsertOne(context.Background(), val)
    if err != nil {
        return 0, err
    }

    return len(p), nil
}
