package tools

import (
	"fmt"
	"time"

	"github.com/42LoCo42/emo/cli/util"
	"github.com/spf13/cobra"
)

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
				fmt.Print("[1J[H") // reset screen

				out, err := conn.CMD(cmd.OutOrStdout(), []string{
					"playback",
					"show",
				})
				if err != "" {
					fmt.Print(err)
				}

				fmt.Print(out)

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
