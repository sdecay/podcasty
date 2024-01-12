package auth

import (
	"errors"
	"net/http"
	"strings"
)

// change to use respondWithError
func GetAPIKey(headers http.Header) (string, error) {
	h := headers.Get("Authorization")
	if h == "" {
		return "", errors.New("no auth info found")
	}

	// FIXME: this is dumb as hell
	vals := strings.Split(h, " ")
	if len(vals) != 2 || vals[0] != "ApiKey" {
		return "", errors.New("malformed auth header")
	}

	return vals[1], nil

	// if len(vals[1]) == 64 {
	// 	return vals[1], nil
	// }
	//
	// return "", errors.New("Malformed auth header")
}
