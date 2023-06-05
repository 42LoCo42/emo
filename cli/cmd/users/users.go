package users

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"github.com/42LoCo42/emo/api"
	"github.com/42LoCo42/emo/cli/util"
	"github.com/42LoCo42/emo/shared"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "users",
		Short: "User management",
	}

	cmd.AddCommand(
		list(),
		set(),
		get(),
		del(),
	)

	return cmd
}

func prettyPrintUser(user *api.User) {
	fmt.Println("Username:         ", user.Name)
	fmt.Println("Is admin:         ", user.IsAdmin)
	fmt.Println("Can upload songs: ", user.CanUploadSongs)
	fmt.Println("Public key:       ", base64.StdEncoding.EncodeToString(user.PublicKey))
}

func list() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Get a list of all users",
		Run: func(cmd *cobra.Command, args []string) {
			resp, err := shared.Client().GetUsers(context.Background())
			if err != nil || resp.StatusCode != http.StatusOK {
				shared.Die(err, "get users request failed")
			}

			data, err := api.ParseGetUsersResponse(resp)
			if err != nil {
				shared.Die(err, "could not parse get users response")
			}

			for _, user := range *data.JSON200 {
				fmt.Println(user.Name)
			}
		},
	}
}

func set() *cobra.Command {
	var (
		isAdmin        *bool
		canUploadSongs *bool
		password       *bool
	)

	cmd := &cobra.Command{
		Use:   "set username",
		Short: "Create or set a user",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			username := args[0]

			// get user

			resp, err := shared.Client().GetUsersName(context.Background(), username)
			if err != nil || resp.StatusCode != http.StatusOK {
				shared.Die(err, "get user by name request failed")
			}

			data, err := api.ParseGetUsersNameResponse(resp)
			if err != nil {
				shared.Die(err, "could not parse get user by name response")
			}

			// for each flag: if set, apply value to user
			cmd.Flags().Visit(func(f *pflag.Flag) {
				if !f.Changed {
					return
				}

				switch f.Name {
				case "isAdmin":
					data.JSON200.IsAdmin = *isAdmin
				case "canUploadSongs":
					data.JSON200.CanUploadSongs = *canUploadSongs
				case "password":
					if !*password {
						break
					}

					password, err := util.AskPassword()
					if err != nil {
						shared.Die(err, "could not read password")
					}

					data.JSON200.PublicKey = util.
						MakeKey([]byte(username), password).
						PublicKey.Bytes
				}
			})

			// post changed user

			resp, err = shared.Client().PostUsers(context.Background(), *data.JSON200)
			if err != nil || resp.StatusCode != http.StatusOK {
				shared.Die(err, "post user request failed")
			}

			prettyPrintUser(data.JSON200)
		},
	}

	isAdmin = cmd.Flags().Bool(
		"isAdmin",
		false,
		"Is this user an administrator?",
	)

	canUploadSongs = cmd.Flags().Bool(
		"canUploadSongs",
		false,
		"Can this user upload songs?",
	)

	password = cmd.Flags().Bool(
		"password",
		false,
		"The password of this user",
	)

	return cmd
}

func get() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get username",
		Short: "Get a user",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			username := args[0]

			resp, err := shared.Client().GetUsersName(context.Background(), username)
			if err != nil || resp.StatusCode != http.StatusOK {
				shared.Die(err, "get user by name request failed")
			}

			data, err := api.ParseGetUsersNameResponse(resp)
			if err != nil {
				shared.Die(err, "could not parse get user by name response")
			}

			prettyPrintUser(data.JSON200)
		},
	}

	return cmd
}

func del() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete username",
		Short: "Delete a user",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			username := args[0]

			resp, err := shared.Client().DeleteUsersName(context.Background(), username)
			if err != nil || resp.StatusCode != http.StatusOK {
				shared.Die(err, "delete user request failed")
			}

			log.Print("Done!")
		},
	}

	return cmd
}
