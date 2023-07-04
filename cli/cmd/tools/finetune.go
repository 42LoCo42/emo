package tools

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"

	"github.com/42LoCo42/emo/cli/util"
	"github.com/42LoCo42/emo/shared"
	"github.com/spf13/cobra"
)

func finetune() *cobra.Command {
	var (
		editor *string
	)

	cmd := &cobra.Command{
		Use:   "finetune",
		Short: "Visually edit all stats of the current user",
		Run: func(cmd *cobra.Command, args []string) {
			stats, err := shared.Client().StatsUserGet(context.Background())
			if err != nil {
				shared.Die(err, "could not get current stats")
			}

			iR, iW := io.Pipe()
			oR, oW := io.Pipe()
			wg := &sync.WaitGroup{}

			// run pipeline
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer log.Print("Pipeline done!")

				if err := util.RunPipeline(
					iR, oW,
					[]string{"column", "-t"},
					[]string{*editor},
				); err != nil {
					shared.Die(err, "could not run pipeline")
				}
			}()

			// send input
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer log.Print("Receiver done!")

				fmt.Fprintln(iW, "#ID\tSong\tCount\tBoost\tSum")
				for _, stat := range stats {
					fmt.Fprintf(
						iW,
						"%v\t%v\t%v\t%v\t%v\n",
						stat.ID,
						stat.Song,
						stat.Count,
						stat.Boost,
						stat.Count+stat.Boost,
					)
				}
				iW.Close()
			}()

			// receive output
			scn := bufio.NewScanner(oR)
			for scn.Scan() {
				line := scn.Text()
				if strings.HasPrefix(line, "#") {
					continue
				}

				log.Print(line)
			}

			log.Print("Waiting for pipeline to stop...")
			wg.Wait()
			log.Print("All done!")
		},
	}

	editor = cmd.Flags().StringP(
		"editor",
		"e",
		"vim",
		"What editor should be used?",
	)

	return cmd
}
