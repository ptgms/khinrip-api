package helper

import (
	"io/ioutil"
	"os"
	"strings"
)

func Validate_key(apikey string) bool {

	var allowed_keys []string

	content, err := ioutil.ReadFile("allowed_keys")
	if err != nil {
		_, _ = os.Create("allowed_keys")
		print("There are no allowed API-Keys yet! Add some to allowed_keys.")
	}

	text := string(content)

	allowed_keys = strings.Split(text, "\n")

	if contains(allowed_keys, apikey) {
		return true
	}

	return false
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if strings.Split(a, "/")[0] == e {
			return true
		}
	}
	return false
}
