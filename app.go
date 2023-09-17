package main

import (
	"database/sql"
)

type App struct {
	db *sql.DB
}

func NewApp(db *sql.DB) App {

	return App{
		db: db,
	}
}
