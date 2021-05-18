package env

import (
	"fmt"
	"os"
)

func UsernamePassword() (string, string, error) {
	username, err := getenv("ENPHOTO_USERNAME")
	if err != nil {
		return "", "", err
	}
	password, err := getenv("ENPHOTO_PASSWORD")
	if err != nil {
		return "", "", err
	}
	return username, password, nil
}

func getenv(key string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		return "", fmt.Errorf("environment variable %s is not defined", key)
	}
	return val, nil
}
