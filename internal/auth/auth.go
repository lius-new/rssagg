package auth

import (
	"errors"
	"net/http"
	"strings"
)

// GetAPIKey extacts an API Key from
// the headers of an http request
// example:
// Authorization: ApiKey {insert apikey here}
// desc: 实际上就是判断header中是否包含Authorization字段,且字段格式是否是: ("ApiKey hsjfhs274sjf-sdfhjsj-234sdfhjsfh")
func GetAPIKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("no Authorization info used")
	}

	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("malformed auth header")
	}

	if vals[0] != "ApiKey" {
		return "", errors.New("malformed first part of auth header")
	}
	return vals[1], nil
}
