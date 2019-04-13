package store

import (
	"github.com/jmoiron/sqlx"
	//_ "github.com/lib/pq"
)

type Store struct {
	*sqlx.DB
}

func NewStore() (*Store, error) {
	store := &Store{}
	db, err := sqlx.Connect("postgres", "user=postgres password=postgres dbname=apigen host=localhost port=5432 sslmode=disable")
	if err != nil {
		return nil, err
	}
	store.DB = db
	return store, nil
}
