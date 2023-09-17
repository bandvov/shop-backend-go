package main

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"
	"unicode"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var sampleSecretKey = []byte(os.Getenv("JWT_SECRET"))

func generateHmacKey() []byte {
	h := sha256.New()
	h.Write(sampleSecretKey)
	return h.Sum(nil)
}

func generateJWT(data interface{}) (string, error) {

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":    data,
		"Expires": time.Now().Add(24 * time.Hour),
	})
	hmacKey := generateHmacKey()

	fmt.Println(hmacKey)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(hmacKey)
	if err != nil {
		fmt.Println(fmt.Errorf("jwt error%+v", err))
	}

	return tokenString, nil
}

func verifyJWT(next func(writer http.ResponseWriter, request *http.Request)) http.HandlerFunc {

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		accessTokenCookie, err := request.Cookie("access-token")
		if err != nil {
			writer.WriteHeader(http.StatusUnauthorized)
			writer.Write([]byte("You're Unauthorized!"))
			return
		}

		token, err := jwt.Parse(accessTokenCookie.Value, func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodECDSA)
			if !ok {
				writer.WriteHeader(http.StatusUnauthorized)
				writer.Write([]byte("You're Unauthorized!"))
			}
			return "", nil
		})
		if err != nil {
			writer.WriteHeader(http.StatusUnauthorized)
			writer.Write([]byte("You're Unauthorized due to error parsing the JWT"))
		}

		if token.Valid {
			next(writer, request)
		} else {
			writer.WriteHeader(http.StatusUnauthorized)
			_, err := writer.Write([]byte("You're Unauthorized due to invalid token"))
			if err != nil {
				return
			}
		}
	})
}

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
func validate(body interface{}, fields ...string) map[string][]string {

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

	return validationErrors
}
