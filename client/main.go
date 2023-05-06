package main

import (
	"fmt"
	"os"

	daemonCmd "github.com/42LoCo42/emo/client/cmd/daemon"
	"github.com/42LoCo42/emo/client/cmd/login"
	"github.com/42LoCo42/emo/client/cmd/songs"
	"github.com/42LoCo42/emo/client/cmd/stats"
	"github.com/42LoCo42/emo/client/cmd/users"
	"github.com/42LoCo42/emo/client/daemon"
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

	rootCmd.PersistentFlags().StringVarP(
		&daemon.SocketPath,
		"socket",
		"s",
		"",
		"Path of the daemon's control socket",
	)

	if err := util.InitClient(); err != nil {
		shared.Die(err, "could not create client")
	}

	rootCmd.AddCommand(
		login.Login,
		users.Cmd(),
		songs.Cmd(),
		stats.Cmd(),
		daemonCmd.Cmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
