package converters

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
)

type Conv func(string) (string, error)

var convs = map[string]Conv{
	"upper": func(s string) (string, error) {
		return strings.ToUpper(s), nil
	},
	"lower": func(s string) (string, error) {
		return strings.ToLower(s), nil
	},
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
	"q1_2tick": func(s string) (string, error) {
		return strings.ReplaceAll(s, "'", "\\`"), nil
	},
	"q2_2tick": func(s string) (string, error) {
		return strings.ReplaceAll(s, "\"", "\\`"), nil
	},
	"q2_escape": func(s string) (string, error) {
		return escapeUnescaped(s, '"'), nil
	},
	"q1_escape": func(s string) (string, error) {
		return escapeUnescaped(s, '\''), nil
	},
}

func escapeUnescaped(s string, ch byte) string {
	var b strings.Builder
	escaped := false

	for i := 0; i < len(s); i++ {
		c := s[i]

		if c == '\\' && !escaped {
			escaped = true
			b.WriteByte(c)
			continue
		}

		if c == '"' {
			if !escaped {
				b.WriteByte('\\')
			}
			b.WriteByte(ch)
		} else {
			b.WriteByte(c)
		}

		escaped = false
	}

	return b.String()
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

func SupportedConvs() []string {
	var res []string
	for k := range convs {
		res = append(res, k)
	}
	return res
}
