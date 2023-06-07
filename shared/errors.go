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

type Causer interface {
	Cause() error
}

func Die(err error, msg string) {
	err = Wrap(err, msg)
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

func Wrap(err error, message string) error {
	if err == nil {
		return errors.New(message)
	} else {
		return errors.Wrap(err, message)
	}
}

func RCause(err error) error {
	for {
		causer, ok := err.(Causer)
		if !ok {
			return err
		}

		err = causer.Cause()
	}
}
