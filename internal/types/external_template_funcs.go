package types

import (
	"encoding/base64"
	"net/url"
	"os"
	"strings"
	"text/template"
)

var templateFuncs = template.FuncMap{
	"split": strings.Split,
	"trim": func(s string) (string, error) {
		return strings.TrimSpace(s), nil
	},
	"urlencode": func(s string) (string, error) {
		return url.QueryEscape(s), nil
	},
	"base64": func(s string) (string, error) {
		return base64.StdEncoding.EncodeToString([]byte(s)), nil
	},
	"base64decode": func(s string) (string, error) {
		res, err := base64.StdEncoding.DecodeString(s)
		if err != nil {
			return "", err
		}
		return string(res), nil
	},
	"readFile": func(s string) (string, error) {
		raw, err := os.ReadFile(s)
		if err != nil {
			return "", err
		}
		return string(raw), nil
	},
}
