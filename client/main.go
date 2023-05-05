package main

import (
	"fmt"
	"os"

	"github.com/42LoCo42/emo/client/cmd/login"
	"github.com/42LoCo42/emo/client/cmd/users"
	"github.com/42LoCo42/emo/client/util"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "emo",
		Short: "easy music organizer",
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(
		&util.Address,
		"address",
		"",
		"Address of the emo server",
	)
	rootCmd.MarkPersistentFlagRequired("address")

	rootCmd.AddCommand(
		login.Login,
		users.Init(),
	)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
