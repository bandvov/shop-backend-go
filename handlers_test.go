package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func TestApp_login(t *testing.T) {
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	user := User{
		UserId:   1,
		Email:    "test@test.aa",
		Password: "$2a$14$Gp3Kfg79XEZR5rV6ahOAjeIK225slynTJgZMhiHzlDLSHFrQKYolq",
	}

	app := &App{db: db}
	query := getUserByEmailQuery

	rows := sqlmock.NewRows([]string{"UserId", "Email", "Password"}).
		AddRow(user.UserId, user.Email, user.Password)

	mock.ExpectQuery(query).WillReturnRows(rows)

	ts := httptest.NewServer(http.HandlerFunc(app.login))
	defer ts.Close()

	marshaled, _ := json.Marshal(User{Email: user.Email, Password: "1"})

	requestBody := bytes.NewBuffer(marshaled)

	fmt.Println(requestBody)
	res, err := http.Post(ts.URL, "application/json", requestBody)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(fmt.Errorf("error in body: %+v", err))
	}
	defer res.Body.Close()

	b, _ := json.Marshal(User{Email: user.Email, UserId: user.UserId})

	if !bytes.Equal(body, b) {
		t.Errorf("Error. want: %+v | have: %+v", string(body), string(b))
	}

}
