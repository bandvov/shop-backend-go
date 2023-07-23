package main

type User struct {
	UserId int    `db:"user_id"`
	Email  string `db:"email"`
}
