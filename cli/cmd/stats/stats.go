package stats

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/42LoCo42/emo/api"
	"github.com/42LoCo42/emo/shared"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Statistics management",
	}

	cmd.AddCommand(
		list(),
		set(),
		get(),
		del(),
		ofMyself(),
		ofUser(),
		ofSong(),
	)

	return cmd
}

func prettyPrintStat(stat *api.Stat) {
	fmt.Println("ID:    ", stat.ID)
	fmt.Println("User:  ", stat.User)
	fmt.Println("Song:  ", stat.Song)
	fmt.Println("Count: ", stat.Count)
	fmt.Println("Boost: ", stat.Boost)
	fmt.Println("Time:  ", stat.Time)
}

func list() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Get a list of statistics",
		Run: func(cmd *cobra.Command, args []string) {
			resp, err := shared.Client().GetStats(context.Background())
			if err != nil || resp.StatusCode != http.StatusOK {
				shared.Die(err, "get stats request failed")
			}

			data, err := api.ParseGetStatsResponse(resp)
			if err != nil {
				shared.Die(err, "could not parse get stats response")
			}

			for _, stat := range *data.JSON200 {
				prettyPrintStat(&stat)
				fmt.Println()
			}
		},
	}
}

func set() *cobra.Command {
	var (
		user  *string
		song  *string
		count *int64
		boost *int64
		time  *float64
	)

	cmd := &cobra.Command{
		Use:   "set ID",
		Short: "Create or set a statistic",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				shared.Die(err, "could not parse ID")
			}

			// get stat
			stat := api.Stat{ID: id}

			resp, err := shared.Client().GetStatsId(context.Background(), id)
			if err != nil || (resp.StatusCode != http.StatusOK &&
				resp.StatusCode != http.StatusNotFound) {
				shared.Die(err, "get stat request failed")
			}

			data, err := api.ParseGetStatsIdResponse(resp)
			if err != nil {
				shared.Die(err, "could not parse get stat response")
			}

			if data.JSON200 != nil {
				stat = *data.JSON200
			}

			// handle flags
			cmd.Flags().Visit(func(f *pflag.Flag) {
				if !f.Changed {
					return
				}

				switch f.Name {
				case "user":
					stat.User = *user
				case "song":
					stat.Song = *song
				case "count":
					stat.Count = *count
				case "boost":
					stat.Boost = *boost
				case "time":
					stat.Time = *time
				}
			})

			// upload new stat
			resp, err = shared.Client().PostStats(context.Background(), stat)
			if err != nil || resp.StatusCode != http.StatusOK {
				shared.Die(err, "post stat request failed")
			}

			prettyPrintStat(&stat)
			log.Print("Done!")
		},
	}

	user = cmd.Flags().StringP(
		"user",
		"u",
		"",
		"The user of this statistic",
	)

	song = cmd.Flags().StringP(
		"song",
		"s",
		"",
		"The song of this statistic",
	)

	count = cmd.Flags().Int64P(
		"count",
		"c",
		0,
		"The count of this statistic",
	)

	boost = cmd.Flags().Int64P(
		"boost",
		"b",
		0,
		"The boost of this statistic",
	)

	time = cmd.Flags().Float64P(
		"time",
		"t",
		0,
		"The total listening time of the statistic",
	)

	return cmd
}

func get() *cobra.Command {
	return &cobra.Command{
		Use:   "get ID",
		Short: "Get a statistic",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				shared.Die(err, "could not parse ID")
			}

			resp, err := shared.Client().GetStatsId(context.Background(), id)
			if err != nil || resp.StatusCode != http.StatusOK {
				shared.Die(err, "get stat request failed")
			}

			data, err := api.ParseGetStatsIdResponse(resp)
			if err != nil {
				shared.Die(err, "could not parse get stat response")
			}

			prettyPrintStat(data.JSON200)
		},
	}
}

func del() *cobra.Command {
	return &cobra.Command{
		Use:   "delete ID",
		Short: "Delete a statistic",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				shared.Die(err, "could not parse ID")
			}

			resp, err := shared.Client().DeleteStatsId(context.Background(), id)
			if err != nil || resp.StatusCode != http.StatusOK {
				shared.Die(err, "get stat request failed")
			}

			log.Print("Done!")
		},
	}
}

func ofMyself() *cobra.Command {
	return &cobra.Command{
		Use:   "ofMyself",
		Short: "Get the statistics of the currently logged in user",
		Run: func(cmd *cobra.Command, args []string) {
			resp, err := shared.Client().GetStatsUser(context.Background())
			if err != nil || resp.StatusCode != http.StatusOK {
				shared.Die(err, "get stats of user request failed")
			}

			data, err := api.ParseGetStatsUserUserResponse(resp)
			if err != nil {
				shared.Die(err, "could not parse stats of user response")
			}

			for _, stat := range *data.JSON200 {
				prettyPrintStat(&stat)
				fmt.Println()
			}
		},
	}
}

func ofUser() *cobra.Command {
	return &cobra.Command{
		Use:   "ofUser user",
		Short: "Get the statistics of a user",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			username := args[0]

			resp, err := shared.Client().GetStatsUserUser(context.Background(), username)
			if err != nil || resp.StatusCode != http.StatusOK {
				shared.Die(err, "get stats of user request failed")
			}

			data, err := api.ParseGetStatsUserUserResponse(resp)
			if err != nil {
				shared.Die(err, "could not parse stats of user response")
			}

			for _, stat := range *data.JSON200 {
				prettyPrintStat(&stat)
				fmt.Println()
			}
		},
	}
}

func ofSong() *cobra.Command {
	return &cobra.Command{
		Use:   "ofSong song",
		Short: "Get the statistics of a song",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			songname := args[0]

			resp, err := shared.Client().GetStatsSongSong(context.Background(), songname)
			if err != nil || resp.StatusCode != http.StatusOK {
				shared.Die(err, "get stats of song request failed")
			}

			data, err := api.ParseGetStatsSongSongResponse(resp)
			if err != nil {
				shared.Die(err, "could not parse stats of song response")
			}

			for _, stat := range *data.JSON200 {
				prettyPrintStat(&stat)
				fmt.Println()
			}
		},
	}
}
