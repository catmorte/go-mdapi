package types

import (
	_ "embed"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/catmorte/go-mdapi/internal/file"
	"github.com/catmorte/go-mdapi/internal/vars"
)

type internalHTTP string

type FieldVar string

const (
	MethodField  FieldVar = "method"
	URLField     FieldVar = "url"
	BodyField    FieldVar = "body"
	HeadersField FieldVar = "headers"
)

func (f FieldVar) Get(vrs vars.Vars) (string, bool) {
	v, ok := vrs[string(f)]
	return v, ok
}

//go:embed internal_http_new_api.md
var internalHTTPTemplate internalHTTP

func (d internalHTTP) GetName() string {
	return "http"
}

func (d internalHTTP) NewAPI() string {
	return string(d)
}

func (d internalHTTP) Run(vrs vars.Vars) error {
	requestURL, ok := URLField.Get(vrs)
	if !ok {
		return errors.New("missing url field")
	}

	method, ok := MethodField.Get(vrs)
	if !ok {
		method = "GET"
	}

	requestBodyRaw, ok := BodyField.Get(vrs)
	var requestBody io.Reader
	if ok {
		requestBody = strings.NewReader(requestBodyRaw)
	}

	rq, err := http.NewRequest(method, requestURL, requestBody)

	var headers http.Header
	headersRaw, ok := HeadersField.Get(vrs)
	if ok {
		headersLines := strings.Split(headersRaw, "\n")
		for _, line := range headersLines {
			parts := strings.Split(line, ":")
			if len(parts) != 2 {
				return fmt.Errorf("invalid header line: %s", line)
			}
			if headers == nil {
				headers = http.Header{}
			}
			headers.Add(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
		}
	}

	if headers != nil {
		rq.Header = headers
	}

	if err != nil {
		return fmt.Errorf("error creating request: %s", err)
	}

	resp, err := http.DefaultClient.Do(rq)
	if err != nil {
		return fmt.Errorf("error making request: %s", err)
	}

	defer resp.Body.Close()

	resultDir := vrs.GetResultDir()

	statusFile := filepath.Join(resultDir, "status")
	headersFile := filepath.Join(resultDir, "headers")
	bodyFile := filepath.Join(resultDir, "body")

	err = os.WriteFile(statusFile, []byte(resp.Status), 0x775)
	if err != nil {
		return fmt.Errorf("error writing status: %w", err)
	}

	sb := strings.Builder{}
	for key, values := range resp.Header {
		for _, value := range values {
			sb.WriteString(fmt.Sprintf("%s: %s\n", key, value))
		}
	}
	err = os.WriteFile(headersFile, []byte(sb.String()), 0x775)
	if err != nil {
		return fmt.Errorf("error writing headers: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading body: %s", err)
	}
	err = os.WriteFile(bodyFile, body, 0x775)
	if err != nil {
		return fmt.Errorf("error writing body: %w", err)
	}

	return nil
}

func compute(key string, typeFields map[string]file.TypedComponent, vrs map[string]string, orDefault ...string) (string, error) {
	v, ok := typeFields[key]
	if !ok {
		if len(orDefault) > 0 {
			return strings.Join(orDefault, ""), nil
		}
		return "", fmt.Errorf("missing %s field", key)
	}

	val, err := v.Compute(vrs)
	if err != nil {
		return "", fmt.Errorf("error computing url: %s", err)
	}

	return val, nil
}
