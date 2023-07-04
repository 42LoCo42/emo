package daemon

import (
	"fmt"
	"os"
	"reflect"

	"github.com/42LoCo42/emo/cli/util"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "daemon <shell | command>",
		Short: "Starts or interacts with the emo player daemon",
		Args:  cobra.MinimumNArgs(1),

		Run: func(cmd *cobra.Command, args []string) {
			if reflect.DeepEqual(args, []string{"shell"}) {
				// start interactive shell
				util.NewDaemonConn().RunShell(os.Stdin, os.Stdout)
			} else {
				// send args as command
				out, err := util.NewDaemonConn().CMD(cmd.OutOrStdout(), args)
				if err != "" {
					fmt.Fprint(cmd.ErrOrStderr(), err)
				} else {
					fmt.Fprint(cmd.OutOrStderr(), out)
				}
			}
		},
	}
}
