package shared

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/42LoCo42/emo/api"
	"github.com/pkg/errors"
)

type StackTracer interface {
	StackTrace() errors.StackTrace
}

type Causer interface {
	Cause() error
}

type Unwrap interface {
	Unwrap() error
}

func Die(err error, msg string) {
	err = Wrap(err, msg)
	Trace(err)
	os.Exit(1)
}

func Trace(err error) {
	log.Print(err)
	if trace, ok := err.(StackTracer); ok {
		for _, frame := range trace.StackTrace() {
			fmt.Fprintf(os.Stderr, "%+v\n", frame)
		}
	}
}

func WrapP(err error, message string, args ...any) error {
	return errors.Wrap(err, fmt.Sprintf(message, args...))
}

func Wrap(err error, message string, args ...any) error {
	if err == nil {
		return fmt.Errorf(message, args...)
	} else {
		return WrapP(err, message, args...)
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

func RUnwrap(err error) error {
	for {
		unwrapped, ok := err.(Unwrap)
		if !ok {
			return err
		}

		err = unwrapped.Unwrap()
	}
}

func ErrorStatus(err error) int {
	err = RUnwrap(err)
	real, ok := err.(*api.ErrorStatusCode)
	if !ok {
		return 0
	}

	return real.StatusCode
}

func Is404(err error) bool {
	return ErrorStatus(err) == http.StatusNotFound
}
