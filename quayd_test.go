package quayd

import (
	"testing"

	"github.com/ejholmes/go-github/github"
)

// TODO Make this work.
func testGitHubStatusesRepository(t *testing.T) {
	repo := "ejholmes/docker-statsd"

	g := github.NewClient(nil)
	r := &GitHubStatusesRepository{RepositoriesService: g.Repositories}

	s := &Status{Repo: repo, Ref: "6607c19", State: "pending", Context: "test"}

	if err := r.Create(s); err != nil {
		t.Fatal(err)
	}
}

// TODO Don't make network requests.
func TestGitHubCommitResolver(t *testing.T) {
	repo := "ejholmes/docker-statsd"

	tests := []struct {
		in  string
		out string
	}{
		{"6607c19", "6607c19d3fd492ec53439f4104b39e4c62ece179"},
	}

	for _, tt := range tests {
		g := github.NewClient(nil)
		r := &GitHubCommitResolver{RepositoriesService: g.Repositories}

		sha, err := r.Resolve(repo, tt.in)
		if err != nil {
			t.Fatal(err)
		}

		if got, want := sha, tt.out; got != want {
			t.Fatalf("Sha => %s; want %s", got, want)
		}
	}
}
