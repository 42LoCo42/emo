package shared

import (
	"github.com/42LoCo42/emo/api"
)

var client *api.Client

func InitClient() error {
	var err error
	client, err = api.NewClient(
		GetConfig().Server,
		// TODO add token
	)
	if err != nil {
		return Wrap(err, "could not create client")
	}

	return nil
}

func Client() *api.Client {
	return client
}
