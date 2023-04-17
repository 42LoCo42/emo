package main

import (
	"log"
	"net/http"

	"github.com/42LoCo42/emo/shared"
)

var token string

func main() {
	Login(InputCreds())
	// Login([]byte("admin"), []byte("admin"))

	log.Print()

	StatQuery()
}

func StatQuery() {
	result, err := JsonRequest[[]shared.StatQuery](
		http.MethodGet, shared.ENDPOINT_STATS)
	if err != nil {
		log.Fatal(err)
	}

	for _, stat := range *result {
		log.Printf(
			"Song %s (ID: %s) - Count %d, Boost %d",
			stat.Name, stat.ID,
			stat.Count, stat.Boost,
		)
	}
}
