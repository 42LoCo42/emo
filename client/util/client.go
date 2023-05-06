package util

import (
	"context"
	"net/http"

	"github.com/42LoCo42/emo/api"
	"github.com/42LoCo42/emo/shared"
	"github.com/pkg/errors"
)

var (
	Address string // set by main from flag
	client  *api.Client
)

func InitClient() error {
	token, err := LoadToken()
	if err != nil {
		return errors.Wrap(err, "could not load token")
	}

	client, err = api.NewClient(
		"undefined:",
		api.WithRequestEditorFn(
			func(ctx context.Context, req *http.Request) error {
				req.Header.Add(shared.AUTH_HEADER, string(token))
				return nil
			},
		),
	)
	if err != nil {
		return errors.Wrap(err, "could not create client")
	}

	return nil
}

func Client() *api.Client {
	client.Server = Address
	return client
}
