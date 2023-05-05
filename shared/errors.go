package shared

import (
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"
)

type StackTracer interface {
	StackTrace() errors.StackTrace
}

func Die(err error, msg string) {
	if err == nil {
		err = errors.New(msg)
	} else if msg != "" {
		err = errors.Wrap(err, msg)
	}

	log.Print(err)
	Trace(err)
	os.Exit(1)
}

func Trace(err error) {
	if trace, ok := err.(StackTracer); ok {
		for _, frame := range trace.StackTrace() {
			fmt.Fprintf(os.Stderr, "%+v\n", frame)
		}
	}
}
