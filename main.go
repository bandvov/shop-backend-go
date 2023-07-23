package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v4"
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
	config, err := pgx.ParseConfig(connString)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatal(err)
	}

	app := NewApp(conn)
	handlers := Handlers{}

	fmt.Print(app)

	http.HandleFunc("/", handlers.getRoot)
	http.HandleFunc("/hello", handlers.getHello(conn))

	err = http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
