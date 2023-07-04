package songs

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/42LoCo42/emo/api"
	"github.com/42LoCo42/emo/shared"
	ht "github.com/ogen-go/ogen/http"
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
	res, err := shared.Client().SongsGet(context.Background())
	if err != nil {
		shared.Die(err, "get songs request failed")
	}

	return res
}

func getSongNames() []string {
	songs := getSongs()
	names := make([]string, len(songs))

	for i, song := range songs {
		names[i] = string(song.Name)
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
			song, err := shared.Client().SongsNameGet(context.Background(), api.SongsNameGetParams{
				Name: api.SongName(songname),
			})
			if err != nil {
				if shared.Is404(err) {
					song = &api.Song{
						Name: api.SongName(songname),
					}
				} else {
					shared.Die(err, "get song request failed")
				}
			}

			var songFile ht.MultipartFile

			cmd.Flags().Visit(func(f *pflag.Flag) {
				if !f.Changed {
					return
				}

				switch f.Name {
				case "file":
					realFile, err := os.Open(*file)
					if err != nil {
						shared.Die(err, "could not open song file")
					}

					songFile = ht.MultipartFile{
						Name: *file,
						File: realFile,
					}
				}
			})

			if err := shared.Client().SongsPost(context.Background(), api.NewOptSongsPostReq(api.SongsPostReq{
				Song: *song,
				File: songFile,
			})); err != nil {
				shared.Die(err, "post song request failed")
			}
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

			res, err := shared.Client().SongsNameGet(context.Background(), api.SongsNameGetParams{
				Name: api.SongName(songname),
			})
			if err != nil {
				shared.Die(err, "get song request failed")
			}

			fmt.Println("Song name: ", res.Name)
			fmt.Println("Song path: ", res.ID)
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

			if err := shared.Client().SongsNameDelete(context.Background(), api.SongsNameDeleteParams{
				Name: api.SongName(songname),
			}); err != nil {
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

			res, err := shared.Client().SongsNameFileGet(context.Background(), api.SongsNameFileGetParams{
				Name: api.SongName(songname),
			})
			if err != nil {
				shared.Die(err, "get song file request failed")
			}

			if _, err := io.Copy(os.Stdout, res.Data); err != nil {
				shared.Die(err, "could not copy file to stdout")
			}
		},
	}
}
