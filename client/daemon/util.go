package daemon

import (
	"log"
	"os"
	"os/user"
	"path"

	"github.com/pkg/errors"
)

var SocketPath string

func GetSocketPath() (string, error) {
	if SocketPath != "" {
		return SocketPath, nil
	}

	if xdgRuntimeDir := os.Getenv("XDG_RUNTIME_DIR"); xdgRuntimeDir != "" {
		return path.Join(xdgRuntimeDir, "emo.socket"), nil
	}

	log.Print("Caution: XDG_RUNTIME_DIR not set, using fallback socket path")

	user, err := user.Current()
	if err != nil {
		return "", errors.Wrap(err, "could not get current user")
	}

	return path.Join("/tmp", user.Username+"-emo.socket"), nil
}
