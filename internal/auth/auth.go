package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetApiKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")

	if val == "" {
		return "", errors.New("Authorization header not found")
	}
	vals := strings.Split(val, " ")

	if len(vals) != 2 {
		return "", errors.New("Malformed Authorization header")
	}

	if vals[0] != "ApiKey" {
		return "", errors.New("Malformed Authorization header")
	}

	return vals[1], nil
}
