package shared

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"reflect"

	"github.com/d5/tengo/v2"
)

var tengoApiFunctions map[string]tengo.Object

type TengoAPIResponse struct {
	tengo.ObjectImpl
	*http.Response
}

func TengoUnpackArgs[T any](args []tengo.Object, f func(any) T) []T {
	values := make([]T, len(args))
	for i, arg := range args {
		values[i] = f(tengo.ToInterface(arg))
	}

	return values
}

func TengoSetup() {
	tengoApiFunctions = make(map[string]tengo.Object)

	initialArgs := []reflect.Value{
		reflect.ValueOf(Client()),
		reflect.ValueOf(context.Background()),
	}
	typ := initialArgs[0].Type()

	for i := 0; i < typ.NumMethod(); i++ {
		met := typ.Method(i)
		// log.Print(met.Name)

		tengoApiFunctions[met.Name] = &tengo.UserFunction{
			Value: func(args ...tengo.Object) (ret tengo.Object, err error) {
				results := met.Func.Call(append(
					initialArgs,
					TengoUnpackArgs(args, reflect.ValueOf)...,
				))

				if len(results) != 2 {
					Die(nil, "IMPOSSIBLE != 2 results of API function")
				}

				// handle error first, it's easier
				if raw := results[1].Interface(); raw != nil {
					err, ok := raw.(error)
					if !ok {
						Die(nil, "IMPOSSIBLE: 2nd result of API function != error")
					}

					if err != nil {
						return nil, err
					}
				}

				res, ok := results[0].Interface().(*http.Response)
				if !ok {
					Die(nil, "IMPOSSIBLE: 1st result of API function != *http.Response")
				}

				return tengo.FromInterface(&TengoAPIResponse{
					Response: res,
				})

			},
		}
	}
}

func RunTengo(
	src []byte,
	out io.Writer,
	functions map[string]func(...tengo.Object) (tengo.Object, error),
) (*tengo.Compiled, error) {
	var (
		mkTengoPrinter = func(f func(io.Writer, ...any) (int, error)) *tengo.UserFunction {
			return &tengo.UserFunction{
				Value: func(args ...tengo.Object) (tengo.Object, error) {
					values := TengoUnpackArgs(args, func(a any) any {
						return a
					})

					n, err := f(out, values...)
					if err != nil {
						return nil, err
					}

					return tengo.FromInterface(n)
				},
			}
		}

		tengoPrint   = mkTengoPrinter(fmt.Fprint)
		tengoPrintln = mkTengoPrinter(fmt.Fprintln)
		tengoPrintf  = mkTengoPrinter(func(w io.Writer, a ...any) (int, error) {
			switch len(a) {
			case 0:
				return -1, fmt.Errorf("printf needs a format string")
			case 1:
				return fmt.Fprint(w, a[0].(string))
			default:
				return fmt.Fprintf(w, a[0].(string), a[1:]...)
			}
		})
	)

	script := tengo.NewScript(src)

	// add retargetable printers
	script.Add("print", tengoPrint)
	script.Add("println", tengoPrintln)
	script.Add("printf", tengoPrintf)

	// add API client & associated functions
	for key, val := range tengoApiFunctions {
		script.Add(key, val)
	}

	// add provided functions
	for key, val := range functions {
		script.Add(key, &tengo.UserFunction{
			Value: val,
		})
	}

	return script.Run()
}
