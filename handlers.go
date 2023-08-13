package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

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

		errors := validate(user)

		if len(errors) > 0 {
			w.WriteHeader(500)
			marshalledErrors, _ := json.Marshal(errors)
			w.Write(marshalledErrors)
			return
		}

		userExists := checkUserExists(user.Email, conn)
		if &userExists != nil {
			http.Error(w, "User already exists", http.StatusBadRequest)
			return
		}

		hashedPassword, err := hashPassword(user.Password)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		user.Password = hashedPassword

		ctx := context.Background()
		ctxWithCancel, cancel := context.WithCancel(ctx)

		defer cancel()

		rows, err := conn.Query(ctxWithCancel, CreateUserQuery, &user.Email, &user.Phone, &user.FullName, &user.Password)

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
		ctx := context.Background()
		ctxWithCancel, cancel := context.WithCancel(ctx)

		rows, err := conn.Query(ctxWithCancel, `SELECT * FROM users`)
		if err != nil {
			log.Fatal(err)
		}
		defer cancel()
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

func (h *Handlers) login(conn *pgx.Conn) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		if r.Body == http.NoBody {
			http.Error(w, "No body", http.StatusBadRequest)
			return
		}

		var body User

		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}

		errors := validate(body)

		if len(errors) > 0 {
			w.WriteHeader(http.StatusInternalServerError)
			marshalledErrors, _ := json.Marshal(errors)
			w.Write(marshalledErrors)
			return
		}

		userExists := checkUserExists(body.Email, conn)
		if &userExists == nil {
			http.Error(w, "User does not exist", http.StatusBadRequest)
			return
		}
		token, err := generateJWT(userExists)
		if err != nil {
			fmt.Println(fmt.Errorf("JWT generating error: %+v", err))
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "access-token",
			Value:    token,
			Path:     "/",
			Secure:   true,
			SameSite: http.SameSiteDefaultMode,
			Expires:  time.Now().Add(24 * time.Hour),
		})
		fmt.Println(userExists)
		res, _ := json.Marshal(userExists)
		w.Write(res)
	}
}
