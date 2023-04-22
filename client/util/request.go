package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func RequestWithBody(
	token, method, endpoint string,
	body io.Reader,
) ([]byte, error) {
	req, err := http.NewRequest(
		method,
		"http://localhost:37812"+endpoint,
		body,
	)
	if err != nil {
		return nil, err
	}

	if token != "" {
		req.Header.Add("authorization", token)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(res.Status)
	}

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, res.Body); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func Request(token, method, endpoint string) ([]byte, error) {
	return RequestWithBody(token, method, endpoint, nil)
}

func JsonRequest[T any](token, method, endpoint string) (*T, string, error) {
	var t T

	body, err := Request(token, method, endpoint)
	if err != nil {
		return nil, "", err
	}

	if err := json.Unmarshal(body, &t); err != nil {
		return nil, "", err
	}

	return &t, string(body), nil
}
