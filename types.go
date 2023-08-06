package main

import "time"

type User struct {
	UserId    int       `db:"user_id" json:"user_id"`
	FullName  string    `db:"full_name" json:"full_name"`
	Email     string    `db:"email" json:"email"`
	Phone     string    `db:"phone" json:"phone,omitempty"`
	Password  string    `db:"password" json:"password"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
