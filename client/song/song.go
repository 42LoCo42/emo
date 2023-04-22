package song

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/42LoCo42/emo/client/util"
	"github.com/42LoCo42/emo/shared"
	"github.com/cristalhq/acmd"
)

var Subcommand = []acmd.Command{
	{
		Name:        "list",
		Description: "Get a list of all songs",
		ExecFunc: func(ctx context.Context, args []string) error {
			_, json, err := util.JsonRequest[[]shared.Song](
				util.Token(), http.MethodGet, shared.ENDPOINT_SONGS,
			)
			if err != nil {
				return err
			}

			fmt.Println(json)
			return nil
		},
	},
	{
		Name:        "get",
		Description: "Get information of a song",
		ExecFunc: func(ctx context.Context, args []string) error {
			if len(args) != 1 {
				return errors.New("Usage: get <song name>")
			}

			name := args[0]
			_, json, err := util.JsonRequest[shared.Song](
				util.Token(), http.MethodGet,
				shared.ENDPOINT_SONGS+"?name="+url.PathEscape(name),
			)
			if err != nil {
				return err
			}

			fmt.Println(json)
			return nil
		},
	},
	{
		Name:        "upload",
		Description: "Upload a song, overriding it if already present",
		ExecFunc: func(ctx context.Context, args []string) error {
			if len(args) < 1 {
				return errors.New("Usage: upload <song files...>")
			}

			for _, songName := range args {
				songFile, err := os.Open(songName)
				if err != nil {
					return err
				}

				if _, err := util.RequestWithBody(
					util.Token(), http.MethodPost,
					shared.ENDPOINT_SONGS+"?name="+url.PathEscape(songName),
					songFile,
				); err != nil {
					return err
				}
			}

			return nil
		},
	},
	{
		Name:        "download",
		Description: "Download a song, print it to stdout",
		ExecFunc: func(ctx context.Context, args []string) error {
			if len(args) != 1 {
				return errors.New("Usage: download <song name>")
			}

			name := args[0]
			song, err := util.Request(
				util.Token(), http.MethodGet,
				shared.ENDPOINT_SONGS+"/"+url.PathEscape(name),
			)
			if err != nil {
				return err
			}

			if _, err := io.Copy(os.Stdout, bytes.NewReader(song)); err != nil {
				return err
			}

			return nil
		},
	},
	{
		Name:        "del",
		Description: "Delete a song",
		ExecFunc: func(ctx context.Context, args []string) error {
			if len(args) < 1 {
				return errors.New("Usage: delete <song names...>")
			}

			for _, songName := range args {
				if _, err := util.Request(
					util.Token(), http.MethodDelete,
					shared.ENDPOINT_SONGS+"?name="+url.PathEscape(songName),
				); err != nil {
					return err
				}
			}

			return nil
		},
	},
}
