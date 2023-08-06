package main

var AddUserQuery = `INSERT INTO users (email, phone, full_name, password) values ($1, $2, $3, $4) RETURNING *;`
