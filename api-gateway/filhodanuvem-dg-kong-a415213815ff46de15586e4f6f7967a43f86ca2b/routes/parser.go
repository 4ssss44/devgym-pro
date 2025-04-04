package routes

import (
	"regexp"
)

func Parse(rawRoute string) (*regexp.Regexp, error) {
	r, _ := regexp.Compile("{[a-zA-Z0-9]+}")

	replaced := r.ReplaceAll([]byte(rawRoute), []byte("([a-zA-Z0-9]+)"))

	return regexp.Compile("^" + string(replaced) + "$")
}
