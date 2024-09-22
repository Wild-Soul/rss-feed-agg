package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// ExtractApiKey extaracts an API Key from the headers of an HTTP Request
// Example:
// Authorization: ApiKey {USER_API_KEY}
func ExtractApiKey(headers http.Header) (string, error) {
	authorizationHeader := headers.Get("Authorization")
	if authorizationHeader == "" {
		return "", errors.New("no authentication info found")
	}

	vals := strings.Split(authorizationHeader, " ")
	if len(vals) != 2 {
		return "", errors.New("malformed auth header (api key not provided)")
	}

	if vals[0] != "ApiKey" {
		return "", errors.New("malformed auth header")
	}

	fmt.Println("Request authenticated")
	return vals[1], nil
}
