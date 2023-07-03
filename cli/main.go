package main

import (
	"fmt"
	"log"
	"os"

	"github.com/42LoCo42/emo/cli/cmd/daemon"
	"github.com/42LoCo42/emo/cli/cmd/login"
	"github.com/42LoCo42/emo/cli/cmd/script"
	"github.com/42LoCo42/emo/cli/cmd/songs"
	"github.com/42LoCo42/emo/cli/cmd/stats"
	"github.com/42LoCo42/emo/cli/cmd/users"
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
	if err := shared.LoadConfig(); err != nil {
		log.Fatal(shared.Wrap(err, "could not load config"))
	}

	if err := shared.InitClient(); err != nil {
		log.Fatal(shared.Wrap(err, "could not init API client"))
	}

	shared.TengoSetup()

	rootCmd.AddCommand(
		daemon.Cmd(),
		login.Login,
		script.Cmd(),
		songs.Cmd(),
		stats.Cmd(),
		users.Cmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
