package main

import (
	"context"
	"database/sql"
	"log"
	"time"
)

func checkUserExists(e string, conn *sql.Conn) User {
	var user User

	// Replace 3 with an ID from your database or another random
	// value to test the no rows use case.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()

	row := conn.QueryRowContext(ctx, getUserByEmailQuery, e)
	err := row.Scan(&user.UserId, &user.Email, &user.Password)
	if err != nil {
		log.Fatal(err)
	}
	return user
}
