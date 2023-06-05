package shared

import (
	"context"
	"net/http"

	"github.com/42LoCo42/emo/api"
)

var client *api.Client

func InitClient() error {
	var err error
	client, err = api.NewClient(
		GetConfig().Server,
		api.WithRequestEditorFn(
			func(ctx context.Context, req *http.Request) error {
				req.Header.Add(AUTH_HEADER, string(GetConfig().Token))
				return nil
			},
		),
	)
	if err != nil {
		return Wrap(err, "could not create client")
	}

	return nil
}

func Client() *api.Client {
	return client
}
