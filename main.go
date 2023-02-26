package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aerogo/aero"
	"github.com/pkg/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	jwtKey []byte
	db     *gorm.DB
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func die(err error, msg string) {
	if msg != "" {
		err = errors.Wrap(err, msg)
	}

	log.Print(err)

	if trace, ok := err.(stackTracer); ok {
		for _, frame := range trace.StackTrace() {
			fmt.Fprintf(os.Stderr, "%+v\n", frame)
		}
	}

	os.Exit(1)
}

func main() {
	var err error

	// get JWT key from environment
	jwtKeyS, ok := os.LookupEnv("EMO_KEY")
	if !ok {
		die(errors.New("Could not get JWT key from EMO_KEY environment variable"), "")
	}
	jwtKey = []byte(jwtKeyS)

	// connect to database
	db, err = gorm.Open(sqlite.Open("emo.db"))
	if err != nil {
		die(err, "Could not open database")
	}

	// setup database
	db.AutoMigrate(
		&User{},
		&Song{},
	)

	// create admin:admin
	admin := User{
		Name: "admin",
		PublicKey: []byte{
			0x47, 0xb0, 0x98, 0x99, 0x49, 0x54, 0xaf, 0xd8,
			0x81, 0xd2, 0xc5, 0x6f, 0x68, 0x61, 0x5a, 0xd7,
			0x4b, 0xfe, 0xe7, 0x06, 0xda, 0x3b, 0x33, 0xde,
			0xc7, 0x17, 0x94, 0x50, 0xaa, 0x67, 0x94, 0x99,
		},
		Admin: true,
	}
	db.FirstOrCreate(&admin, admin)

	// create aero application
	app := aero.New()
	app.Config.Ports.HTTP = 37812
	app.Use(logRequests, secureEndpointCheck)

	// serve static files
	app.Get("/*file", func(ctx aero.Context) error {
		return ctx.File("static/" + ctx.Get("file"))
	})

	app.Post("/auth", auth)

	app.Run()
}
