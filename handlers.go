package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func (a *App) createUser(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	conn, err := a.db.Conn(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if r.Body == http.NoBody {
		http.Error(w, "No body", http.StatusBadRequest)
		return
	}

	var user User

	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	errors := validate(user)

	if len(errors) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		marshalledErrors, _ := json.Marshal(errors)
		w.Write(marshalledErrors)
		return
	}

	userExists := checkUserExists(user.Email, conn)

	if userExists.Email != "" {
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	user.Password = hashedPassword

	stmt, err := conn.PrepareContext(ctx, CreateUserQuery)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, &user.Email, &user.Phone, &user.FullName, &user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Write([]byte("User created"))
}

func (a *App) getUsers(w http.ResponseWriter, r *http.Request) {
	var users = make([]*User, 0)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	conn, err := a.db.Conn(ctx)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := conn.QueryContext(ctx, `SELECT user_id,email,phone,full_name FROM users`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		user := &User{}
		if err := rows.Scan(&user.UserId, &user.Email, &user.Phone, &user.FullName); err != nil {
			log.Fatalf("Scan: %v", err)
		}

		users = append(users, user)
	}
	marshaled, err := json.Marshal(users)
	w.Write(marshaled)
}

func (a *App) login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	conn, err := a.db.Conn(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if r.Body == http.NoBody {
		http.Error(w, "No body", http.StatusBadRequest)
		return
	}

	var body User

	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	errors := validate(body, "password")

	if len(errors) > 0 {
		w.WriteHeader(http.StatusInternalServerError)
		marshalledErrors, _ := json.Marshal(errors)
		w.Write(marshalledErrors)
		return
	}

	userExists := checkUserExists(body.Email, conn)
	if userExists.Email == "" {
		http.Error(w, "User does not exist", http.StatusBadRequest)
		return
	}
	hashedPassword, err := hashPassword(body.Password)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if hashedPassword != userExists.Password{
		http.Error(w, "Credentials not correct", http.StatusBadRequest)
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
	res, _ := json.Marshal(userExists)
	w.Write(res)
}
