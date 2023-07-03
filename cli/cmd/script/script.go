package script

import (
	"os"

	"github.com/42LoCo42/emo/shared"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var (
		daemon *bool
	)

	cmd := &cobra.Command{
		Use:   "script file.tengo",
		Short: "Execute a Tengo script",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			path := args[0]

			file, err := os.ReadFile(path)
			if err != nil {
				shared.Die(err, "could not read script")
			}

			if *daemon {
				panic("TODO: run tengo on daemon")
			} else {
				if _, err := shared.RunTengo(
					file,
					os.Stdout,
					nil,
				); err != nil {
					shared.Die(err, "script execution failed")
				}
			}
		},
	}

	daemon = cmd.Flags().BoolP(
		"daemon",
		"d",
		false,
		"Run on daemon instead of CLI",
	)

	return cmd
}
