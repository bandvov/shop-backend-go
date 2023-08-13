package main

import "time"

type User struct {
	UserId    int       `db:"user_id" json:"user_id"`
	Email     string    `db:"email" json:"email" validate:"required"`
	Password  string    `db:"password" json:"password" validate:"required"`
	Phone     string    `db:"phone" json:"phone,omitempty"`
	FullName  string    `db:"full_name" json:"full_name,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at,omitempty"`
}
