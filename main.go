package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	var err error
	connString, err := getEnvVariable("CONN_STRING")
	if err != nil {
		log.Fatal(err)
	}
	port, err := getEnvVariable("PORT")
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatal(err)
	}

	app := NewApp(db)

	http.HandleFunc("/register", app.createUser)
	http.HandleFunc("/login", app.login)
	http.HandleFunc("/users", app.getUsers)
	http.HandleFunc("/", http.NotFound)

	err = http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
