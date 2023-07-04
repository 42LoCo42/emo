package users

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"

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
	fmt.Println("Username:         ", user.ID)
	fmt.Println("Is admin:         ", user.IsAdmin)
	fmt.Println("Can upload songs: ", user.CanUploadSongs)
	fmt.Println("Public key:       ", base64.StdEncoding.EncodeToString(user.PublicKey))
}

func getUsers() []api.User {
	res, err := shared.Client().UsersGet(context.Background())
	if err != nil {
		shared.Die(err, "get users request failed")
	}

	return res
}

func getUserNames() []string {
	users := getUsers()
	names := make([]string, len(users))

	for i, user := range users {
		names[i] = string(user.ID)
	}

	return names
}

func ArgsUserNames(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return getUserNames(), cobra.ShellCompDirectiveNoFileComp
}

func list() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Get a list of all users",
		Run: func(cmd *cobra.Command, args []string) {
			names := getUserNames()
			for _, name := range names {
				fmt.Println(name)
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
		Use:               "set username",
		Short:             "Create or set a user",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: ArgsUserNames,
		Run: func(cmd *cobra.Command, args []string) {
			username := args[0]

			// get user
			user, err := shared.Client().UsersNameGet(context.Background(), api.UsersNameGetParams{
				Name: api.UserName(username),
			})
			if err != nil {
				if shared.Is404(err) {
					user = &api.User{
						ID: api.UserName(username),
					}
				} else {
					shared.Die(err, "get user by name request failed")
				}
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
			if err := shared.Client().UsersPost(context.Background(), api.NewOptUser(*user)); err != nil {
				shared.Die(err, "post user request failed")
			}

			prettyPrintUser(user)
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
		Use:               "get username",
		Short:             "Get a user",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: ArgsUserNames,
		Run: func(cmd *cobra.Command, args []string) {
			username := args[0]

			user, err := shared.Client().UsersNameGet(context.Background(), api.UsersNameGetParams{
				Name: api.UserName(username),
			})
			if err != nil {
				shared.Die(err, "get user by name request failed")
			}

			prettyPrintUser(user)
		},
	}

	return cmd
}

func del() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "delete username",
		Short:             "Delete a user",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: ArgsUserNames,
		Run: func(cmd *cobra.Command, args []string) {
			username := args[0]

			if err := shared.Client().UsersNameDelete(context.Background(), api.UsersNameDeleteParams{
				Name: api.UserName(username),
			}); err != nil {
				shared.Die(err, "delete user request failed")
			}

			log.Print("Done!")
		},
	}

	return cmd
}
