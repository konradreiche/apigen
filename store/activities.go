package store

import (
	"time"

	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"
	"github.com/konradreiche/apigen/model"
)

func (store *Store) SaveActivity(m *model.Activity) error {
	return saveActivity(store, m)
}

func saveActivity(e sqlx.Ext, m *model.Activity) error {
	if m == nil {
		return errors.New("provided model can not be nil")
	}
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now().UTC()
	}
	_, err := sqlx.NamedExec(e,
		`INSERT INTO activities
	(
		id,
		method,
		data,
		error,
		created_at
	)
	VALUES
	(
		:id,
		:method,
		:data,
		:error,
		:created_at
	)`, m)
	return errors.Wrap(err, "SaveActivity failed")
}
