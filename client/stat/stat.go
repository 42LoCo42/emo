package stat

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/42LoCo42/emo/client/util"
	"github.com/42LoCo42/emo/shared"
	"github.com/cristalhq/acmd"
)

var Subcommand = []acmd.Command{
	{
		Name:        "get",
		Description: "Get the statistics for a song",
		ExecFunc: func(ctx context.Context, args []string) error {
			if len(args) != 1 {
				return errors.New("Usage: get <song name>")
			}

			songName := args[0]

			_, json, err := util.JsonRequest[shared.StatQuery](
				util.Token(), http.MethodGet,
				shared.ENDPOINT_STATS+"?name="+url.PathEscape(songName),
			)
			if err != nil {
				return err
			}

			fmt.Println(json)
			return nil
		},
	},
	{
		Name:        "count",
		Description: "Set the count value in a statistic of a song",
		ExecFunc: func(ctx context.Context, args []string) error {
			if len(args) != 2 {
				return errors.New("Usage: count <song name> <count value>")
			}

			songName := args[0]
			countVal := args[1]

			if _, err := util.Request(
				util.Token(), http.MethodPost,
				shared.ENDPOINT_STATS+fmt.Sprintf(
					"?name=%s&count=%s",
					url.PathEscape(songName),
					url.QueryEscape(countVal),
				),
			); err != nil {
				return err
			}

			return nil
		},
	},
	{
		Name:        "boost",
		Description: "Set the count value in a statistic of a song",
		ExecFunc: func(ctx context.Context, args []string) error {
			if len(args) != 2 {
				return errors.New("Usage: boost <song name> <boost value>")
			}

			songName := args[0]
			boostVal := args[1]

			if _, err := util.Request(
				util.Token(), http.MethodPost,
				shared.ENDPOINT_STATS+fmt.Sprintf(
					"?name=%s&boost=%s",
					url.PathEscape(songName),
					url.QueryEscape(boostVal),
				),
			); err != nil {
				return err
			}

			return nil
		},
	},
	{
		Name:        "del",
		Description: "Delete a song statistic",
		ExecFunc: func(ctx context.Context, args []string) error {
			if len(args) != 1 {
				return errors.New("Usage: del <song name>")
			}

			songName := args[0]

			if _, err := util.Request(
				util.Token(), http.MethodDelete,
				shared.ENDPOINT_STATS+"?name="+url.PathEscape(songName),
			); err != nil {
				return err
			}

			return nil
		},
	},
}
