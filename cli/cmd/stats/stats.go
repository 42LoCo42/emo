package stats

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/42LoCo42/emo/api"
	"github.com/42LoCo42/emo/cli/cmd/songs"
	"github.com/42LoCo42/emo/cli/cmd/users"
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

func getStats() []api.Stat {
	res, err := shared.Client().StatsGet(context.Background())
	if err != nil {
		shared.Die(err, "get stats request failed")
	}

	return res
}

func getStatIDStrings() []string {
	stats := getStats()
	ids := make([]string, len(stats))

	for i, stat := range stats {
		ids[i] = fmt.Sprint(stat.ID)
	}

	return ids
}

func ArgsStatIDs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return getStatIDStrings(), cobra.ShellCompDirectiveDefault
}

func list() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Get a list of statistics",
		Run: func(cmd *cobra.Command, args []string) {
			stats := getStats()
			for _, stat := range stats {
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
		Use:               "set ID",
		Short:             "Create or set a statistic",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: ArgsStatIDs,
		Run: func(cmd *cobra.Command, args []string) {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				shared.Die(err, "could not parse ID")
			}

			// get stat
			stat, err := shared.Client().StatsIDGet(context.Background(), api.StatsIDGetParams{
				ID: api.StatID(id),
			})
			if err != nil {
				if shared.Is404(err) {
					stat = &api.Stat{
						ID: api.StatID(id),
					}
				} else {
					shared.Die(err, "get stat request failed")
				}
			}

			// handle flags
			cmd.Flags().Visit(func(f *pflag.Flag) {
				if !f.Changed {
					return
				}

				switch f.Name {
				case "user":
					stat.User = api.UserName(*user)
				case "song":
					stat.Song = api.SongName(*song)
				case "count":
					stat.Count = *count
				case "boost":
					stat.Boost = *boost
				case "time":
					stat.Time = *time
				}
			})

			// upload new stat
			stat, err = shared.Client().StatsPost(context.Background(), api.NewOptStat(*stat))
			if err != nil {
				shared.Die(err, "post stat request failed")
			}

			prettyPrintStat(stat)
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
		Use:               "get ID",
		Short:             "Get a statistic",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: ArgsStatIDs,
		Run: func(cmd *cobra.Command, args []string) {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				shared.Die(err, "could not parse ID")
			}

			res, err := shared.Client().StatsIDGet(context.Background(), api.StatsIDGetParams{
				ID: api.StatID(id),
			})
			if err != nil {
				shared.Die(err, "get stat request failed")
			}

			prettyPrintStat(res)
		},
	}
}

func del() *cobra.Command {
	return &cobra.Command{
		Use:               "delete ID",
		Short:             "Delete a statistic",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: ArgsStatIDs,
		Run: func(cmd *cobra.Command, args []string) {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				shared.Die(err, "could not parse ID")
			}

			if err := shared.Client().StatsIDDelete(context.Background(), api.StatsIDDeleteParams{
				ID: api.StatID(id),
			}); err != nil {
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
			res, err := shared.Client().StatsUserGet(context.Background())
			if err != nil {
				shared.Die(err, "get stats of user request failed")
			}

			for _, stat := range res {
				prettyPrintStat(&stat)
				fmt.Println()
			}
		},
	}
}

func ofUser() *cobra.Command {
	return &cobra.Command{
		Use:               "ofUser user",
		Short:             "Get the statistics of a user",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: users.ArgsUserNames,
		Run: func(cmd *cobra.Command, args []string) {
			username := args[0]

			res, err := shared.Client().StatsUserUserGet(context.Background(), api.StatsUserUserGetParams{
				User: api.UserName(username),
			})
			if err != nil {
				shared.Die(err, "get stats of user request failed")
			}

			for _, stat := range res {
				prettyPrintStat(&stat)
				fmt.Println()
			}
		},
	}
}

func ofSong() *cobra.Command {
	return &cobra.Command{
		Use:               "ofSong song",
		Short:             "Get the statistics of a song",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: songs.ArgsSongNames,
		Run: func(cmd *cobra.Command, args []string) {
			songname := args[0]

			res, err := shared.Client().StatsSongSongGet(context.Background(), api.StatsSongSongGetParams{
				Song: api.SongName(songname),
			})
			if err != nil {
				shared.Die(err, "get stats of song request failed")
			}

			for _, stat := range res {
				prettyPrintStat(&stat)
				fmt.Println()
			}
		},
	}
}
