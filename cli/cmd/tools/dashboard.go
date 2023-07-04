package tools

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/42LoCo42/emo/cli/util"
	"github.com/42LoCo42/emo/shared"
	"github.com/spf13/cobra"
	"github.com/wangjia184/sortedset"
)

func MinInt(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func dashboard() *cobra.Command {
	var (
		delay *float64
	)

	cmd := &cobra.Command{
		Use:   "dashboard",
		Short: "Display a continuously updating dashboard",
		Run: func(cmd *cobra.Command, args []string) {
			conn := util.NewDaemonConn()

			for {
				// get current stats
				stats, err := shared.Client().StatsUserGet(context.Background())
				if err != nil {
					shared.Die(err, "could not get current stats")
				}

				// define top count
				top := MinInt(10, len(stats))

				// reset screen
				fmt.Print("[H[J")

				// print currently playing song
				{
					fmt.Println("[32;1mCurrently playing:[m")
					out, err := conn.CMD(cmd.OutOrStdout(), []string{
						"playback",
						"show",
					})
					if err != "" {
						fmt.Print(err)
					}

					fmt.Println(out)
				}

				{
					fmt.Printf("[32;1mTop %v highest count:[m\n", top)

					sort.Slice(stats, func(i, j int) bool {
						return stats[i].Count > stats[j].Count
					})
					for _, stat := range stats[:top] {
						fmt.Printf("%v: %v\n", stat.Song, stat.Count)
					}
					fmt.Println()
				}

				{
					fmt.Printf("[32;1mTop %v highest boost:[m\n", top)

					sort.Slice(stats, func(i, j int) bool {
						return stats[i].Boost > stats[j].Boost
					})
					for _, stat := range stats[:top] {
						fmt.Printf("%v: %v\n", stat.Song, stat.Boost)
					}
					fmt.Println()
				}

				{
					categories := sortedset.New()
					for _, stat := range stats {
						category := strings.Split(string(stat.Song), "/")[0]

						old := 0
						if tmp := categories.GetByKey(category); tmp != nil {
							old = int(tmp.Score())
						}

						categories.AddOrUpdate(
							category,
							sortedset.SCORE(old+int(stat.Count)),
							nil,
						)
					}

					top := MinInt(5, categories.GetCount())

					fmt.Printf("[32;1mTop %v categories:[m\n", top)
					for i := 0; i < top; i++ {
						fmt.Println(categories.PopMax().Key())
					}
					fmt.Println()
				}

				{
					totalCount := 0.0
					for _, stat := range stats {
						totalCount += stat.Time
					}

					fmt.Println("[32;1mTotal song count:[m", len(stats))
					fmt.Println("[32;1mTotal watchtime: [m", time.Duration(float64(time.Second)*totalCount))
				}

				time.Sleep(time.Duration(float64(time.Second) * *delay))
			}
		},
	}

	delay = cmd.Flags().Float64P(
		"delay",
		"d",
		1,
		"Delay between updates in seconds",
	)

	return cmd
}
