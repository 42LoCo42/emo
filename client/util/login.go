package util

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/42LoCo42/emo/api"
	"github.com/jamesruan/sodium"
	"golang.org/x/crypto/argon2"
)

func Login(
	client *api.Client,
	username,
	password []byte,
) (
	token []byte,
	err error,
) {
	var seed sodium.BoxSeed
	seed.Bytes = argon2.IDKey(
		password,
		username,
		uint32(sodium.CryptoPWHashOpsLimitInteractive),
		uint32(sodium.CryptoPWHashMemLimitInteractive>>10),
		1,
		uint32(seed.Size()),
	)

	kp := sodium.SeedBoxKP(seed)

	raw, err := client.GetLoginUser(context.Background(), string(username))
	if err != nil {
		return nil, err
	}

	if raw.StatusCode != http.StatusOK {
		return nil, errors.New("login request not successful")
	}

	resp, err := api.ParseGetLoginUserResponse(raw)
	if err != nil {
		return nil, err
	}

	encrypted, err := base64.StdEncoding.DecodeString(string(resp.Body))
	if err != nil {
		return nil, err
	}

	token, err = sodium.Bytes(encrypted).SealedBoxOpen(kp)
	if err != nil {
		return nil, err
	}

	return token, nil
}
