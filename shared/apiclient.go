package shared

import (
	"github.com/42LoCo42/emo/api"
	"net/http"
)

var client *api.Client

type CustomClient struct {
	http.Client
	token string
}

// Do implements http.Client
func (c *CustomClient) Do(r *http.Request) (*http.Response, error) {
	r.Header.Add(AUTH_HEADER, c.token)
	return c.Client.Do(r)
}

func InitClient() error {
	var err error
	client, err = api.NewClient(
		GetConfig().Server,
		api.WithClient(&CustomClient{
			token: string(GetConfig().Token),
		}),
	)
	if err != nil {
		return Wrap(err, "could not create client")
	}

	return nil
}

func Client() *api.Client {
	return client
}
