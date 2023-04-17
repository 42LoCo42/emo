package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

func Request(
	method, endpoint string,
) (string, error) {
	req, err := http.NewRequest(
		method,
		"http://localhost:37812"+endpoint,
		nil,
	)
	if err != nil {
		return "", err
	}

	req.Header.Add("authorization", token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		return "", errors.New(res.Status)
	}

	buf := new(strings.Builder)
	if _, err := io.Copy(buf, res.Body); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func JsonRequest[T any](method, endpoint string) (*T, error) {
	var t T

	body, err := Request(method, endpoint)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(body), &t); err != nil {
		return nil, err
	}

	return &t, nil
}
