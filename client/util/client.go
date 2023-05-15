package util

import (
	"context"
	"log"
	"net/http"

	"github.com/42LoCo42/emo/api"
	"github.com/42LoCo42/emo/shared"
)

var (
	Address string // set by main from flag
	client  *api.Client
)

func InitClient() error {
	token, err := LoadToken()
	if err != nil {
		log.Print(shared.Wrap(err, "warning: no token loaded"))
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
		return shared.Wrap(err, "could not create client")
	}

	return nil
}

func Client() *api.Client {
	client.Server = Address
	return client
}
