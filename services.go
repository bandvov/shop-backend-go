package main

import (
	"context"
	"database/sql"
	"time"
)

func checkUserExists(e string, conn *sql.Conn) User {
	var user User

	// Replace 3 with an ID from your database or another random
	// value to test the no rows use case.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()

	row := conn.QueryRowContext(ctx, getUserByEmailQuery, e)
	_ = row.Scan(&user.Email)

	return user
}
