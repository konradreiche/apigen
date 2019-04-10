package model

import (
	"encoding/json"
	"time"
)

type Activity struct {
	ID        int             `db:"id"          json:"-"`
	Method    string          `db:"method"      json:"method"`
	Data      json.RawMessage `db:"data"        json:"data,omitempty"`
	Error     string          `db:"error"       json:"error,omitempty"`
	CreatedAt time.Time       `db:"created_at"  json:"-"`
}
