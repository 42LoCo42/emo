package songs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/42LoCo42/emo/api"
	"github.com/42LoCo42/emo/shared"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "songs",
		Short: "Song management",
	}

	cmd.AddCommand(
		list(),
		set(),
		get(),
		del(),
		file(),
	)

	return cmd
}

func list() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Get a list of all songs",
		Run: func(cmd *cobra.Command, args []string) {
			resp, err := shared.Client().GetSongs(context.Background())
			if err != nil || resp.StatusCode != http.StatusOK {
				shared.Die(err, "get songs request failed")
			}

			data, err := api.ParseGetSongsResponse(resp)
			if err != nil {
				shared.Die(err, "could not parse get songs response")
			}

			for _, song := range *data.JSON200 {
				fmt.Println(song.Name)
			}
		},
	}
}

func set() *cobra.Command {
	var (
		file *string
	)

	cmd := &cobra.Command{
		Use:   "set songname",
		Short: "Create or change a song",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			songname := args[0]

			// get song info
			song := api.SongInfo{
				Name: songname,
			}

			resp, err := shared.Client().GetSongsName(context.Background(), songname)
			if err != nil || (resp.StatusCode != http.StatusOK &&
				resp.StatusCode != http.StatusNotFound) {
				shared.Die(err, "get song requesst failed")
			}

			data, err := api.ParseGetSongsNameResponse(resp)
			if err != nil {
				shared.Die(err, "parse get song request failed")
			}

			if data.JSON200 != nil {
				song = *data.JSON200
			}

			// build post body
			buf := bytes.NewBuffer([]byte{})
			writer := multipart.NewWriter(buf)

			// handle flags
			cmd.Flags().Visit(func(f *pflag.Flag) {
				if !f.Changed {
					return
				}

				switch f.Name {
				case "file":
					fileWriter, err := writer.CreateFormFile("File", *file)
					if err != nil {
						shared.Die(err, "could not create file field")
					}

					file, err := os.Open(*file)
					if err != nil {
						shared.Die(err, "could not open song file")
					}

					if _, err := io.Copy(fileWriter, file); err != nil {
						shared.Die(err, "could not write to file field")
					}
				}
			})

			// write info to multipart body
			infoWriter, err := writer.CreateFormField("Info")
			if err != nil {
				shared.Die(err, "could not create info field")
			}

			if err := json.NewEncoder(infoWriter).Encode(song); err != nil {
				shared.Die(err, "could not write to info field")
			}

			writer.Close()

			// upload song
			resp, err = shared.Client().PostSongsWithBody(
				context.Background(),
				"multipart/form-data; boundary = "+writer.Boundary(),
				buf,
			)
			if err != nil || resp.StatusCode != http.StatusOK {
				shared.Die(err, "post song request failed")
			}

			log.Print("Done!")
		},
	}

	file = cmd.Flags().StringP(
		"file",
		"f",
		"",
		"The song file to upload",
	)

	return cmd
}

func get() *cobra.Command {
	return &cobra.Command{
		Use:   "get songname",
		Short: "Get a song's information",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			songname := args[0]

			resp, err := shared.Client().GetSongsName(context.Background(), songname)
			if err != nil || resp.StatusCode != http.StatusOK {
				shared.Die(err, "get song request failed")
			}

			data, err := api.ParseGetSongsNameResponse(resp)
			if err != nil {
				shared.Die(err, "could not parse get song response")
			}

			fmt.Println("Song name: ", data.JSON200.Name)
		},
	}
}

func del() *cobra.Command {
	return &cobra.Command{
		Use:   "delete songname",
		Short: "Delete a song",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			songname := args[0]

			resp, err := shared.Client().DeleteSongsName(context.Background(), songname)
			if err != nil || resp.StatusCode != http.StatusOK {
				shared.Die(err, "delete song request failed")
			}

			log.Print("Done!")
		},
	}
}

func file() *cobra.Command {
	return &cobra.Command{
		Use:   "file songname",
		Short: "Get a song file",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			songname := args[0]

			resp, err := shared.Client().GetSongsNameFile(context.Background(), songname)
			if err != nil || resp.StatusCode != http.StatusOK {
				shared.Die(err, "get song file request failed")
			}

			if _, err := io.Copy(os.Stdout, resp.Body); err != nil {
				shared.Die(err, "could not copy file to stdout")
			}
		},
	}
}
