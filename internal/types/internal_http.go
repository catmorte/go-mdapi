package types

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/catmorte/go-mdapi/internal/vars"
)

type internalHTTP string

const (
	InternalHTTPMethodField   FieldVar = "method"
	InternalHTTPURLField      FieldVar = "url"
	InternalHTTPBodyField     FieldVar = "body"
	InternalHTTPBodyFileField FieldVar = "bodyFile"
	InternalHTTPHeadersField  FieldVar = "headers"
	InternalHTTPFormField     FieldVar = "form"
)

//go:embed internal_http_new_api.md
var internalHTTPTemplate internalHTTP

func (d internalHTTP) GetName() string {
	return "http"
}

func (d internalHTTP) NewAPI() string {
	return string(internalHTTPTemplate)
}

func (d internalHTTP) Run(vrs vars.Vars) error {
	requestURL, ok := InternalHTTPURLField.Get(vrs)
	if !ok {
		return errors.New("missing url field")
	}

	method, ok := InternalHTTPMethodField.Get(vrs)
	if !ok {
		method = "GET"
	}

	headers, err := d.headers(vrs)
	if err != nil {
		return err
	}

	requestBody, contentType, err := d.buildRequestBody(vrs)
	if err != nil {
		return err
	}
	if requestBody == nil {
		requestBody = http.NoBody
	}

	predefinedContentType := headers.Get("Content-Type")
	if len(predefinedContentType) == 0 && contentType != "" {
		headers.Set("Content-Type", contentType)
	}

	rq, err := http.NewRequest(method, requestURL, requestBody)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	rq.Header = headers

	resp, err := http.DefaultClient.Do(rq)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
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
		return fmt.Errorf("error reading body: %w", err)
	}
	err = os.WriteFile(bodyFile, body, 0x775)
	if err != nil {
		return fmt.Errorf("error writing body: %w", err)
	}

	return nil
}

func (d internalHTTP) headers(vrs vars.Vars) (http.Header, error) {
	var headers http.Header
	headersRaw, ok := InternalHTTPHeadersField.Get(vrs)
	if ok {
		headersLines := strings.Split(headersRaw, "\n")
		for _, line := range headersLines {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid header line: %s", line)
			}
			if headers == nil {
				headers = http.Header{}
			}
			headers.Add(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
		}
	}
	return headers, nil
}

func (d internalHTTP) buildRequestBody(vrs vars.Vars) (io.Reader, string, error) {
	if body, ok := InternalHTTPBodyField.Get(vrs); ok {
		return strings.NewReader(body), "", nil
	}

	if filePath, ok := InternalHTTPBodyFileField.Get(vrs); ok {
		file, err := os.Open(filePath)
		if err != nil {
			return nil, "", fmt.Errorf("failed to open request file: %w", err)
		}
		return file, "", nil
	}

	if form, ok := InternalHTTPFormField.Get(vrs); ok {
		bodyBuf := &bytes.Buffer{}
		writer := multipart.NewWriter(bodyBuf)

		lines := strings.Split(form, "\n")
		for _, line := range lines {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) != 2 {
				return nil, "", fmt.Errorf("invalid form line: %s", line)
			}

			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])

			if strings.HasPrefix(key, "@") {
				filePath := strings.TrimPrefix(key, "@")
				file, err := os.Open(filePath)
				if err != nil {
					return nil, "", fmt.Errorf("failed to open form file %q: %w", filePath, err)
				}
				defer file.Close()

				part, err := writer.CreateFormFile("file", filepath.Base(filePath))
				if err != nil {
					return nil, "", fmt.Errorf("failed to create form file part: %w", err)
				}
				if _, err := io.Copy(part, file); err != nil {
					return nil, "", fmt.Errorf("failed to copy form file: %w", err)
				}
			} else {
				if strings.HasPrefix(val, "\\@") {
					val = strings.TrimPrefix(val, "\\")
				}
				if err := writer.WriteField(key, val); err != nil {
					return nil, "", fmt.Errorf("failed to write form field: %w", err)
				}
			}
		}

		if err := writer.Close(); err != nil {
			return nil, "", fmt.Errorf("failed to close multipart writer: %w", err)
		}

		return bodyBuf, writer.FormDataContentType(), nil
	}

	return nil, "", nil
}

func (d internalHTTP) Compile(vrs vars.Vars) error {
	fmt.Println("not supported for internal commands")
	return nil
}
