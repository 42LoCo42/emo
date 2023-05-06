package daemon

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/42LoCo42/emo/client/daemon/playback"
	"github.com/42LoCo42/emo/client/daemon/queue"
	"github.com/42LoCo42/emo/client/daemon/state"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func Run() error {
	socketPath, err := GetSocketPath()
	if err != nil {
		return errors.Wrap(err, "could not get socket path")
	}

	os.Remove(socketPath)

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		return errors.Wrap(err, "could not open listener")
	}

	state := state.State{}
	log.Print("Daemon has started, listening on ", socketPath)

	for {
		client, err := listener.Accept()
		if err != nil {
			return errors.Wrap(err, "could not accept client")
		}

		go handleClient(client, &state)
	}
}

func handleClient(client net.Conn, state *state.State) {
	defer client.Close()

	cmd := &cobra.Command{}

	cmd.AddCommand(
		queue.Cmd(state),
		playback.Cmd(state),
	)

	cmd.SetOut(client)
	cmd.SetErr(client)

	var mode [1]byte
	if _, err := client.Read(mode[:]); err != nil {
		log.Print("Could not read client mode: ", err)
		return
	}

	if mode[0] == 'j' {
		// json & single command only
		var args []string
		if err := json.NewDecoder(client).Decode(&args); err != nil {
			log.Print("Could not decode JSON command: ", err)
			return
		}

		cmd.SetArgs(args)
		if err := cmd.Execute(); err != nil {
			fmt.Fprintln(client, "Command error: ", err)
		}

		client.Close()
	} else {
		// shell mode
		fmt.Fprintln(
			client,
			"Connected to emo daemon! Run help for a command list",
		)

		scn := bufio.NewReader(client)

		for {
			fmt.Fprint(client, "🎵 > ")
			line, err := scn.ReadString('\n')
			if err != nil && err != io.EOF {
				log.Print("Could not read client command: ", err)
				return
			}

			cmd.SetArgs(strings.Split(strings.TrimSpace(line), " "))
			if err := cmd.Execute(); err != nil {
				fmt.Fprintln(client, "Command error: ", err)
			}
		}
	}
}
