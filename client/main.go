package main

import (
	"fmt"
	"os"

	"github.com/42LoCo42/emo/client/cmd/login"
	"github.com/42LoCo42/emo/client/cmd/songs"
	"github.com/42LoCo42/emo/client/cmd/stats"
	"github.com/42LoCo42/emo/client/cmd/users"
	"github.com/42LoCo42/emo/client/util"
	"github.com/42LoCo42/emo/shared"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "emo",
		Short: "easy music organizer",
	}
)

func main() {
	rootCmd.PersistentFlags().StringVarP(
		&util.Address,
		"address",
		"a",
		"",
		"Address of the emo server",
	)

	if err := util.InitClient(); err != nil {
		shared.Die(err, "could not create client")
	}

	rootCmd.AddCommand(
		login.Login,
		users.Cmd(),
		songs.Cmd(),
		stats.Cmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
