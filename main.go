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

	postgresUser, err := getEnvVariable("POSTGRES_DB_USER")
	if err != nil {
		log.Fatal(err)
	}

	postgresUserPassword, err := getEnvVariable("POSTGRES_DB_USER_PASSWORD")
	if err != nil {
		log.Fatal(err)
	}

	postgresHost, err := getEnvVariable("POSTGRES_DB_HOST")
	if err != nil {
		log.Fatal(err)
	}

	postgresPort, err := getEnvVariable("POSTGRES_DB_PORT")
	if err != nil {
		log.Fatal(err)
	}

	connString := fmt.Sprintf("postgresql://%v:%v@%v:%v/postgres?sslmode=disable", postgresUser, postgresUserPassword, postgresHost, postgresPort)

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
