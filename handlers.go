package main

import (
	"fmt"
	"io"
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
