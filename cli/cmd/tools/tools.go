package tools

import "github.com/spf13/cobra"

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tools",
		Short: "Useful tools to interact with Emo",
	}

	cmd.AddCommand(
		dashboard(),
	)

	return cmd
}
