package converters

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
)

type Conv func(string) (string, error)

var convs = map[string]Conv{
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
}

func Convert(text string, cs []string) (string, error) {
	for _, cKey := range cs {
		c, ok := convs[cKey]
		if !ok {
			return "", fmt.Errorf("unknown converter %s", cKey)
		}

		var err error
		text, err = c(text)
		if err != nil {
			return "", err
		}
	}
	return text, nil
}
