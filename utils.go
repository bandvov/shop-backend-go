package main

import (
	"errors"
	"fmt"
	"os"
)

func getEnvVariable(name string) (string, error) {
	var errMessage = "No %v in environment variables"
	connString, exists := os.LookupEnv(name)
	if !exists {
		return "", errors.New(fmt.Sprintf(errMessage, name))
	}
	return connString, nil
}
