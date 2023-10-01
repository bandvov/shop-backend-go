package main

import "time"

type User struct {
	UserId    int       `db:"user_id" json:"user_id,omitempty"`
	Email     string    `db:"email" json:"email,omitempty" validate:"required"`
	Password  string    `db:"password" json:"password,omitempty" validate:"required"`
	Phone     string    `db:"phone,omitempty" json:"phone,omitempty"`
	FullName  string    `db:"full_name,omitempty" json:"full_name,omitempty"`
	CreatedAt time.Time `db:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt time.Time `db:"updated_at,omitempty" json:"updated_at,omitempty"`
}

type ValidationErrors map[string][]string

type LoginResponse struct {
	UserId    int       `json:"user_id,omitempty"`
	Email     string    `json:"email,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at,omitempty"`
}
