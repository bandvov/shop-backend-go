package main

import "time"

type User struct {
	UserId    int       `db:"user_id" json:"user_id,omitempty"`
	Email     string    `db:"email" json:"email,omitempty" validate:"required"`
	Password  string    `db:"password" json:"password,omitempty" validate:"required"`
	Phone     string    `db:"phone" json:"phone,omitempty"`
	FullName  string    `db:"full_name" json:"full_name,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at,omitempty"`
}
