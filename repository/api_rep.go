package repository

import "github.com/jmoiron/sqlx"

type MyapiPostgres struct {
	db *sqlx.DB
}

func NewMyapiPostgres(db *sqlx.DB) *MyapiPostgres {
	return &MyapiPostgres{db: db}
}
