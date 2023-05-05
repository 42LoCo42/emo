package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/42LoCo42/emo/api"
	"github.com/42LoCo42/emo/client/util"
	"github.com/42LoCo42/emo/shared"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var LoginCMD = &cobra.Command{
	Use:   "login username",
	Short: "Log in to an emo server",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]

		client, err := api.NewClient("http://localhost:37812")
		if err != nil {
			shared.Die(err, "could not create client")
		}

		fmt.Fprint(os.Stderr, "Password: ")
		password, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			shared.Die(err, "could not read credentials")
		}
		fmt.Fprintln(os.Stderr)

		token, err := util.Login(client, []byte(username), password)
		if err != nil {
			shared.Die(err, "login failed")
		}

		if err := util.SaveToken(token); err != nil {
			shared.Die(err, "could not save token")
		}

		log.Print("Login successful!")
	},
}
