package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jackc/pgx/v4"
)

type Handlers struct{}

func (h *Handlers) createUser(conn *pgx.Conn) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		if r.Body == http.NoBody {
			http.Error(w, "No body", http.StatusBadRequest)
			return
		}

		var user User

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}

		isNotValid, errors := validate(user)

		if isNotValid {
			w.WriteHeader(500)
			marshalledErrors, _ := json.Marshal(errors)
			w.Write(marshalledErrors)
			return
		}

		userExists := checkUserExists(user.Email, conn)
		if userExists {
			http.Error(w, "User already exists", http.StatusBadRequest)
			return
		}

		hashedPassword, err := hashPassword(user.Password)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		user.Password = hashedPassword

		rows, err := conn.Query(context.Background(), CreateUserQuery, &user.Email, &user.Phone, &user.FullName, &user.Password)
		defer rows.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
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
