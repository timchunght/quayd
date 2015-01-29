package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/remind101/quayd"
)

func main() {
	var (
		port  = flag.String("port", "8080", "The port to run the server on.")
		token = flag.String("github-token", "", "The GitHub API Token to use when creating commit statuses.")
		auth  = flat.String("registry-auth", "", "The Quay username and password user:pass ")
	)
	flag.Parse()

	q := quayd.New(*token, *auth)
	s := quayd.NewServer(q)

	log.Fatal(http.ListenAndServe(":"+*port, s))
}
