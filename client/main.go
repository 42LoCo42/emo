package main

import (
	"fmt"
	"os"

	"github.com/42LoCo42/emo/client/cmd"
	"github.com/spf13/cobra"
)

var (
	address string

	rootCmd = &cobra.Command{
		Use:   "emo",
		Short: "easy music organizer",
	}
)

func init() {
	rootCmd.AddCommand(
		cmd.LoginCMD,
	)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
