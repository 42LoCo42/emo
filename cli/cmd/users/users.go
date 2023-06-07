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
			user := api.User{Name: username}

			resp, err := shared.Client().GetUsersName(context.Background(), username)
			if err != nil || (resp.StatusCode != http.StatusOK &&
				resp.StatusCode != http.StatusNotFound) {
				shared.Die(err, "get user by name request failed")
			}

			data, err := api.ParseGetUsersNameResponse(resp)
			if err != nil {
				shared.Die(err, "could not parse get user by name response")
			}
			if data.JSON200 != nil {
				user = *data.JSON200
			}

			// for each flag: if set, apply value to user
			cmd.Flags().Visit(func(f *pflag.Flag) {
				if !f.Changed {
					return
				}

				switch f.Name {
				case "isAdmin":
					user.IsAdmin = *isAdmin
				case "canUploadSongs":
					user.CanUploadSongs = *canUploadSongs
				case "password":
					if !*password {
						break
					}

					password, err := util.AskPassword()
					if err != nil {
						shared.Die(err, "could not read password")
					}

					user.PublicKey = util.
						MakeKey([]byte(username), password).
						PublicKey.Bytes
				}
			})

			// post changed user

			resp, err = shared.Client().PostUsers(context.Background(), user)
			if err != nil || resp.StatusCode != http.StatusOK {
				shared.Die(err, "post user request failed")
			}

			prettyPrintUser(&user)
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
