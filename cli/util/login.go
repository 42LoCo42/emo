package util

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/42LoCo42/emo/api"
	"github.com/jamesruan/sodium"
)

func Login(
	client *api.Client,
	username,
	password []byte,
) (
	token []byte,
	err error,
) {
	key := MakeKey(username, password)

	raw, err := client.GetLoginUser(context.Background(), string(username))
	if err != nil {
		return nil, err
	}

	if raw.StatusCode != http.StatusOK {
		return nil, errors.New("login request failed")
	}

	resp, err := api.ParseGetLoginUserResponse(raw)
	if err != nil {
		return nil, err
	}

	encrypted, err := base64.StdEncoding.DecodeString(string(resp.Body))
	if err != nil {
		return nil, err
	}

	token, err = sodium.Bytes(encrypted).SealedBoxOpen(key)
	if err != nil {
		return nil, err
	}

	return token, nil
}
