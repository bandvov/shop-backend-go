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
		defer r.Body.Close()
		if r.Body == http.NoBody {
			log.Println("No body")
			http.Error(w, "No body", http.StatusBadRequest)
			return
		}
		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			log.Printf("Error reading body: %v\n", err)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}
		fmt.Println(user)
		hashedPassword, err := hashPassword(user.Password)
		if err != nil {
			log.Printf("Error while hashing: %v\n", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		user.Password = hashedPassword
		validate(user, "email", "phone", "password")

		_, err = conn.Query(context.Background(), `INSERT INTO users (email, phone, full_name, password) values ($1, $2, $3, $4) RETURNING *;`, &user.Email, &user.Phone, &user.FullName, &user.Password)
		if err != nil {
			log.Fatal(err)
		}
		w.Write([]byte("User created"))
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
			if err := rows.Scan(&user.UserId, &user.Email, &user.Password, &user.Phone, &user.FullName); err != nil {
				log.Fatalf("Scan: %v", err)
			}

			users = append(users, user)
		}
		marshaled, err := json.Marshal(users)
		w.Write(marshaled)
	}
}
