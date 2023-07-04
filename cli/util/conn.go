package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/42LoCo42/emo/shared"
)

type DaemonConn struct {
	net.Conn
	jsonModeInit bool
	jsonDecoder  *json.Decoder
}

func NewDaemonConn() *DaemonConn {
	socketPath := shared.GetConfig().Socket
	log.Print("Connecting to ", socketPath)

	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		shared.Die(err, "could not connect")
	}

	return &DaemonConn{
		Conn: conn,
	}
}

func (c *DaemonConn) RunShell(
	i io.Reader,
	o io.Writer,
) {
	// init shell mode
	fmt.Fprint(c, "s")

	// copy input to daemon, close on EOF
	go func() {
		io.Copy(c, i)
		c.Close()
	}()

	// copy daemon to output
	io.Copy(o, c)
}

func (c *DaemonConn) CMD(o io.Writer, args []string) (out, err string) {
	if !c.jsonModeInit {
		fmt.Fprint(c, "j")
		c.jsonModeInit = true
		c.jsonDecoder = json.NewDecoder(c)
	}

	buf := bytes.NewBuffer(nil)
	if err := json.NewEncoder(buf).Encode(args); err != nil {
		shared.Die(err, "could not encode args")
	}

	if _, err := io.Copy(c, buf); err != nil {
		shared.Die(err, "could not copy args to connection")
	}

	var res []string
	if err := c.jsonDecoder.Decode(&res); err != nil {
		shared.Die(err, "could not decode result of JSON command")
	}

	return res[0], res[1]
}
