package main

var CreateUserQuery = `INSERT INTO users (email, phone, full_name, password) values ($1, $2, $3, $4) RETURNING *;`

var getUserByEmailQuery = `SELECT user_id, email, phone from users where email=$1;`

var getAllUsersQuery = `SELECT user_id,email,phone,full_name FROM users`
