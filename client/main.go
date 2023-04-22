package main

import (
	"context"
	"log"
	"os"

	"github.com/42LoCo42/emo/client/song"
	"github.com/42LoCo42/emo/client/stat"
	"github.com/42LoCo42/emo/client/util"
	"github.com/cristalhq/acmd"
)

func main() {
	cmds := []acmd.Command{
		{
			Name:        "song",
			Subcommands: song.Subcommand,
		},
		{
			Name:        "stat",
			Subcommands: stat.Subcommand,
		},
		{
			Name:        "login",
			Description: "Log in to an emo server",
			ExecFunc: func(ctx context.Context, args []string) error {
				tokenPath, err := util.TokenFilePath()
				if err != nil {
					return err
				}

				token := util.Login(util.InputCreds())
				return os.WriteFile(tokenPath, []byte(token), 0600)
			},
		},
	}

	r := acmd.RunnerOf(cmds, acmd.Config{
		AppName:        "emo",
		AppDescription: "easy music organizer",
		Version:        "0.0.1",
	})

	if err := r.Run(); err != nil {
		log.Fatal("Error: ", err)
	}
}
