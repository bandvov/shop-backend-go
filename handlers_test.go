package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestApp_getUsers(t *testing.T) {
	db, mock, err := sqlmock.New()
	user := User{
		UserId:   1,
		Email:    "test@test.aa",
		Phone:    "123-123-123",
		FullName: "test",
	}
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	app := &App{db: db}
	query := getAllUsersQuery

	rows := sqlmock.NewRows([]string{"UserId", "Email", "Phone", "FullName"}).
		AddRow(user.UserId, user.Email, user.Phone, user.FullName)

	mock.ExpectQuery(query).WillReturnRows(rows)
	ts := httptest.NewServer(http.HandlerFunc(app.getUsers))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(res.Body)

	defer res.Body.Close()

	var users []User

	users = append(users, user)

	b, _ := json.Marshal(users)

	if !bytes.Equal(body, b) {
		t.Errorf("Error. want: %v | have: %v", string(body), string(b))
	}
}
