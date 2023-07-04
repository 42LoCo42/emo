package util

import (
	"context"
	"encoding/base64"
	"io/ioutil"

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

	res, err := client.LoginUserGet(context.Background(), api.LoginUserGetParams{
		User: api.UserName(username),
	})
	if err != nil {
		return nil, err
	}

	encrypted, err := ioutil.ReadAll(
		base64.NewDecoder(base64.StdEncoding, res.Data),
	)
	if err != nil {
		return nil, err
	}

	token, err = sodium.Bytes(encrypted).SealedBoxOpen(key)
	if err != nil {
		return nil, err
	}

	return token, nil
}
