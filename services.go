package main

import (
	"context"

	"github.com/jackc/pgx/v4"
)

func checkUserExists(e string, conn *pgx.Conn) bool {
	var email string

	// Replace 3 with an ID from your database or another random
	// value to test the no rows use case.
	row := conn.QueryRow(context.Background(), getUserByEmailQuery, e)
	_ = row.Scan(&email)
	if email != "" {
		return true
	} else {
		return false
	}
}
