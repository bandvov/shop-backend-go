package main

import (
	"github.com/jackc/pgx/v4"
)

type App struct {
	Conn *pgx.Conn
}

func NewApp(conn *pgx.Conn) *App {

	return &App{
		Conn: conn,
	}
}
