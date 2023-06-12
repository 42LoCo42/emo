package songs

import (
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

func getSongs() []api.Song {
	resp, err := shared.Client().GetSongs(context.Background())
	if err != nil || resp.StatusCode != http.StatusOK {
		shared.Die(err, "get songs request failed")
	}

	data, err := api.ParseGetSongsResponse(resp)
	if err != nil {
		shared.Die(err, "could not parse get songs response")
	}

	return *data.JSON200
}

func getSongNames() []string {
	songs := getSongs()
	names := make([]string, len(songs))

	for i, song := range songs {
		names[i] = song.Name
	}

	return names
}

func ArgsSongNames(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return getSongNames(), cobra.ShellCompDirectiveDefault
}

func list() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Get a list of all songs",
		Run: func(cmd *cobra.Command, args []string) {
			names := getSongNames()
			for _, name := range names {
				fmt.Println(name)
			}
		},
	}
}

func set() *cobra.Command {
	var (
		file *string
	)

	cmd := &cobra.Command{
		Use:               "set songname",
		Short:             "Create or change a song",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: ArgsSongNames,
		Run: func(cmd *cobra.Command, args []string) {
			songname := args[0]

			// get song info
			song := api.Song{Name: songname}

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

			done := make(chan int) // channel for waiting
			pR, pW := io.Pipe()    // pipe for data exchange
			body := multipart.NewWriter(pW)

			// sender goroutine
			go func() {
				// upload song
				resp, err = shared.Client().PostSongsWithBody(
					context.Background(),
					"multipart/form-data; boundary = "+body.Boundary(),
					pR,
				)
				if err != nil || resp.StatusCode != http.StatusOK {
					shared.Die(err, "post song request failed")
				}

				done <- 0 // signal done
			}()

			// write info to multipart body
			infoWriter, err := body.CreateFormField("Info")
			if err != nil {
				shared.Die(err, "could not create info field")
			}
			if err := json.NewEncoder(infoWriter).Encode(song); err != nil {
				shared.Die(err, "could not write to info field")
			}

			// handle flags
			cmd.Flags().Visit(func(f *pflag.Flag) {
				if !f.Changed {
					return
				}

				switch f.Name {
				case "file":
					fileWriter, err := body.CreateFormFile("File", *file)
					if err != nil {
						shared.Die(err, "could not create file field")
					}

					file, err := os.Open(*file)
					if err != nil {
						shared.Die(err, "could not open song file")
					}

					// write song to multipart body
					if _, err := io.Copy(fileWriter, file); err != nil {
						shared.Die(err, "could not write to file field")
					}
				}
			})

			body.Close() // no more data in body
			pW.Close()   // request finished
			<-done       // wait for sender to quit
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
		Use:               "get songname",
		Short:             "Get a song's information",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: ArgsSongNames,
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
			fmt.Println("Song path:", data.JSON200.ID)
		},
	}
}

func del() *cobra.Command {
	return &cobra.Command{
		Use:               "delete songname",
		Short:             "Delete a song",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: ArgsSongNames,
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
		Use:               "file songname",
		Short:             "Get a song file",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: ArgsSongNames,
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
