package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
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

// capitalize is used to capitalize string
func capitalize(str string) string {
	runes := []rune(str)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// validate is used to validate fields in struct
func validate(body interface{}, fields ...string) (bool, map[string][]string) {

	var bodyMap map[string]interface{}
	inrec, _ := json.Marshal(body)
	json.Unmarshal(inrec, &bodyMap)

	t := reflect.TypeOf(body)

	var validationErrors = make(map[string][]string)

	for fieldName := range bodyMap {
		field, found := t.FieldByName(capitalize(fieldName))
		if !found {
			continue
		}

		validateTagValue := strings.Split(field.Tag.Get("validate"), ",")

		for _, v := range validateTagValue {
			if v == "required" && len(bodyMap[fieldName].(string)) == 0 {
				validationErrors[fieldName] = append(validationErrors[fieldName], fmt.Sprintf("%v is required", fieldName))
			}
		}
	}
	var isNotValid = len(validationErrors) > 0
	return isNotValid, validationErrors
}
