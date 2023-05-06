package daemon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"reflect"

	"github.com/42LoCo42/emo/client/daemon"
	"github.com/42LoCo42/emo/shared"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "daemon [shell | command]",
		Short: "Starts or interacts with the emo player daemon",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				// start daemon
				if err := daemon.Run(); err != nil {
					shared.Die(err, "error in daemon")
				}
			} else {
				// get socket path
				socketPath, err := daemon.GetSocketPath()
				if err != nil {
					shared.Die(err, "could not get socket path")
				}

				log.Print("Connecting to ", socketPath)

				// connect
				conn, err := net.Dial("unix", socketPath)
				if err != nil {
					shared.Die(err, "could not connect")
				}

				if reflect.DeepEqual(args, []string{"shell"}) {
					// start interactive shell
					fmt.Fprintln(conn) // a byte other than 'j' starts shell

					// copy input to daemon, close on EOF
					go func() {
						io.Copy(conn, os.Stdin)
						conn.Close()
					}()

					// copy daemon to output
					io.Copy(os.Stdout, conn)
				} else {
					// send args as JSON

					// encode command
					buf := bytes.NewBuffer(nil)
					if err := json.NewEncoder(buf).Encode(args); err != nil {
						shared.Die(err, "could not encode args")
					}

					// send
					fmt.Fprintln(conn, "j", buf.String())

					// get output
					if _, err := io.Copy(os.Stdout, conn); err != nil {
						shared.Die(err, "could not copy output")
					}
				}
			}
		},
	}
}