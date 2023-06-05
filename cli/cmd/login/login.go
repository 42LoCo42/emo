package login

import (
	"log"

	"github.com/42LoCo42/emo/cli/util"
	"github.com/42LoCo42/emo/shared"
	"github.com/spf13/cobra"
)

var Login = &cobra.Command{
	Use:   "login username",
	Short: "Log in to an emo server",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]

		password, err := util.AskPassword()
		if err != nil {
			shared.Die(err, "could not read password")
		}

		token, err := util.Login(shared.Client(), []byte(username), password)
		if err != nil {
			shared.Die(err, "login failed")
		}

		shared.GetConfig().Token = token
		if err := shared.SaveConfig(); err != nil {
			shared.Die(err, "could not save token")
		}

		log.Print("Login successful!")
	},
}
