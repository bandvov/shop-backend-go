package main

import (
	"github.com/jackc/pgx/v4"
)

type App struct {
	Conn     *pgx.Conn
	Handlers Handlers
}

func NewApp(conn *pgx.Conn, h Handlers) App {

	return App{
		Conn:     conn,
		Handlers: Handlers{},
	}
}
