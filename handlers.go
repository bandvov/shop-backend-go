package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/jackc/pgx/v4"
)

type Handlers struct{}

func (h *Handlers) getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	io.WriteString(w, "This is my website!\n")
}
func (h *Handlers) getHello(conn *pgx.Conn) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("got /hello request\n")
		io.WriteString(w, "Hello, HTTP!\n")
	}
}
func (h *Handlers) addUser(conn *pgx.Conn) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var users = make([]*User, 0)

		rows, err := conn.Query(context.Background(), `INSERT INTO users (email) values ($1) RETURNING *;`, "new")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			user := &User{}
			if err := rows.Scan(&user.UserId, &user.Email); err != nil {
				log.Fatalf("Scan: %v", err)
			}

			users = append(users, user)
		}
		marshaled, err := json.Marshal(users)

		w.Write(marshaled)
	}

}

func (h *Handlers) getUsers(conn *pgx.Conn) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var users = make([]*User, 0)
		rows, err := conn.Query(context.Background(), `SELECT * FROM users`)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			user := &User{}
			if err := rows.Scan(&user.UserId, &user.Email); err != nil {
				log.Fatalf("Scan: %v", err)
			}

			users = append(users, user)
		}
		marshaled, err := json.Marshal(users)
		w.Write(marshaled)
	}
}
