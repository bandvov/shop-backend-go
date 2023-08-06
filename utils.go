package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

func getEnvVariable(name string) (string, error) {
	var errMessage = "No %v in environment variables"
	connString, exists := os.LookupEnv(name)
	if !exists {
		return "", errors.New(fmt.Sprintf(errMessage, name))
	}
	return connString, nil
}
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func capitalize(str string) string {
	runes := []rune(str)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func validate(body interface{}, fields ...string) (bool, []error) {

	var bodyMap map[string]interface{}
	inrec, _ := json.Marshal(body)
	json.Unmarshal(inrec, &bodyMap)

	// var errors []error
	for _, v := range fields {

		fmt.Println(bodyMap[v])
	}
	return true, nil
}
