package main

import (
	"flag"
	"log"
	"net/http"

	"code.google.com/p/goauth2/oauth"

	"github.com/ejholmes/go-github/github"
	"github.com/remind101/quayd"
)

func main() {
	var (
		port  = flag.String("port", "8080", "The port to run the server on.")
		token = flag.String("github-token", "", "The GitHub API Token to use when creating commit statuses.")
	)
	flag.Parse()

	t := &oauth.Transport{
		Token: &oauth.Token{AccessToken: *token},
	}

	gh := github.NewClient(t.Client())
	s := quayd.NewServer(quayd.NewStatusesService(gh))

	log.Fatal(http.ListenAndServe(":"+*port, s))
}
