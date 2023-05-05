package users

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"github.com/42LoCo42/emo/api"
	"github.com/42LoCo42/emo/client/util"
	"github.com/42LoCo42/emo/shared"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func Init() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "users",
		Short: "User management",
	}

	cmd.AddCommand(
		list(),
		get(),
		set(),
		del(),
	)

	return cmd
}

func list() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Get a list of all users",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := util.NewClient()
			if err != nil {
				shared.Die(err, "could not create client")
			}

			resp, err := client.GetUsers(context.Background())
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

	return cmd
}

func get() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get username",
		Short: "Get a user",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			username := args[0]

			client, err := util.NewClient()
			if err != nil {
				shared.Die(err, "could not create client")
			}

			resp, err := client.GetUsersName(context.Background(), username)
			if err != nil || resp.StatusCode != http.StatusOK {
				shared.Die(err, "get user by name request failed")
			}

			data, err := api.ParseGetUsersNameResponse(resp)
			if err != nil {
				shared.Die(err, "could not parse get user by name response")
			}

			fmt.Println("Username:         ", data.JSON200.Name)
			fmt.Println("Is admin:         ", data.JSON200.IsAdmin)
			fmt.Println("Can upload songs: ", data.JSON200.CanUploadSongs)
			fmt.Println("Public key:       ",
				base64.StdEncoding.EncodeToString(data.JSON200.PublicKey))
		},
	}

	return cmd
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

			client, err := util.NewClient()
			if err != nil {
				shared.Die(err, "could not create client")
			}

			// get user

			resp, err := client.GetUsersName(context.Background(), username)
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

			resp, err = client.PostUsers(context.Background(), *data.JSON200)
			if err != nil || resp.StatusCode != http.StatusOK {
				shared.Die(err, "post user request failed")
			}

			log.Print("Done!")
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

func del() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete username",
		Short: "Delete a user",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			username := args[0]

			client, err := util.NewClient()
			if err != nil {
				shared.Die(err, "could not create client")
			}

			resp, err := client.DeleteUsersName(context.Background(), username)
			if err != nil || resp.StatusCode != http.StatusOK {
				shared.Die(err, "delete user request failed")
			}

			log.Print("Done!")
		},
	}

	return cmd
}
