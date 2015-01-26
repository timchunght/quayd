package quayd

import (
	"strings"

	"github.com/ejholmes/go-github/github"
)

var (
	DefaultStatusesRepository = &FakeStatusesRepository{}
	DefaultCommitResolver     = &FakeCommitResolver{}
	DefaultStatusesService    = &StatusesService{
		StatusesRepository: DefaultStatusesRepository,
		CommitResolver:     DefaultCommitResolver,
	}
)

// Status represents a GitHub Commit Status.
type Status struct {
	Repo    string
	Ref     string
	State   string
	Context string
}

// StatusesRepository is an interface that can be implemented for creating
// Commit Statuses.
type StatusesRepository interface {
	Create(*Status) error
}

// FakeStatusesRepository is a fake implementation of the StatusesRepository
// interface.
type FakeStatusesRepository struct {
	Statuses []*Status
}

// Create implements StatusesRepository Create.
func (r *FakeStatusesRepository) Create(status *Status) error {
	r.Statuses = append(r.Statuses, status)

	return nil
}

// Reset resets the collection of Statuses.
func (r *FakeStatusesRepository) Reset() {
	r.Statuses = nil
}

// GitHubStatusesRepository is an implementation of the StatusesRepository
// interface backed by a github.Client.
type GitHubStatusesRepository struct {
	*github.Client
}

// Create implements StatusesRepository Create.
func (r *GitHubStatusesRepository) Create(status *Status) error {
	st := &github.RepoStatus{
		State:   &status.State,
		Context: &status.Context,
	}

	// Split `owner/repo` into ["owner", "repo"].
	c := strings.Split(status.Repo, "/")

	_, _, err := r.Client.Repositories.CreateStatus(
		c[0],
		c[1],
		status.Ref,
		st,
	)
	return err
}

// CommitResolver is an interface for resolving a short sha to a full 40
// character sha.
type CommitResolver interface {
	Resolve(repo, short string) (string, error)
}

// FakeCommitResolver returns the short sha prefixed with the string "long".
type FakeCommitResolver struct{}

// Resolve implements CommitResolver Resolve.
func (cr *FakeCommitResolver) Resolve(repo, short string) (string, error) {
	return "long-" + short, nil
}

// GitHubCommitResolver is an implementation of CommitResolver backed by a
// github.Client.
type GitHubCommitResolver struct {
	*github.Client
}

// Resolve implements CommitResolver Resolve.
func (cr *GitHubCommitResolver) Resolve(repo, short string) (string, error) {
	// Split `owner/repo` into ["owner", "repo"].
	c := strings.Split(repo, "/")

	cm, _, err := cr.Client.Repositories.GetCommit(
		c[0],
		c[1],
		short,
	)
	if err != nil {
		return "", err
	}

	return *cm.SHA, nil
}

// StatusesService provides a convenient server for creating new commit
// statuses.
type StatusesService struct {
	StatusesRepository
	CommitResolver
}

func (s *StatusesService) Create(repo, ref, state string) error {
	sha, err := s.Resolve(repo, ref)
	if err != nil {
		return err
	}

	return s.StatusesRepository.Create(&Status{
		Repo:    repo,
		Ref:     sha,
		State:   state,
		Context: "docker",
	})
}
