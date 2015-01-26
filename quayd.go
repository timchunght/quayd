package quayd

import (
	"strings"

	"code.google.com/p/goauth2/oauth"

	"github.com/ejholmes/go-github/github"
)

var (
	// Context is the string that will be displayed when showing the commit
	// status.
	Context = "Docker Image"

	// DefaultStatusesRepository is the default StatusesRepository to use.
	DefaultStatusesRepository = &statusesRepository{}

	// DefaultCommitResolver is the default CommitResolver to use.
	DefaultCommitResolver = &commitResolver{}

	// Default is the default Quayd to use.
	Default = &Quayd{}
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

// statusesRepository is a fake implementation of the StatusesRepository
// interface.
type statusesRepository struct {
	statuses []*Status
}

// Create implements StatusesRepository Create.
func (r *statusesRepository) Create(status *Status) error {
	r.statuses = append(r.statuses, status)

	return nil
}

// Reset resets the collection of Statuses.
func (r *statusesRepository) Reset() {
	r.statuses = nil
}

// GitHubStatusesRepository is an implementation of the StatusesRepository
// interface backed by a github.Client.
type GitHubStatusesRepository struct {
	RepositoriesService interface {
		CreateStatus(owner, repo, ref string, status *github.RepoStatus) (*github.RepoStatus, *github.Response, error)
	}
}

// Create implements StatusesRepository Create.
func (r *GitHubStatusesRepository) Create(status *Status) error {
	st := &github.RepoStatus{
		State:   &status.State,
		Context: &status.Context,
	}

	// Split `owner/repo` into ["owner", "repo"].
	c := strings.Split(status.Repo, "/")

	_, _, err := r.RepositoriesService.CreateStatus(
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

// commitResolver returns the short sha prefixed with the string "long".
type commitResolver struct{}

// Resolve implements CommitResolver Resolve.
func (cr *commitResolver) Resolve(repo, short string) (string, error) {
	return "long-" + short, nil
}

// GitHubCommitResolver is an implementation of CommitResolver backed by a
// github.Client.
type GitHubCommitResolver struct {
	RepositoriesService interface {
		GetCommit(owner, repo, sha string) (*github.RepositoryCommit, *github.Response, error)
	}
}

// Resolve implements CommitResolver Resolve.
func (cr *GitHubCommitResolver) Resolve(repo, short string) (string, error) {
	// Split `owner/repo` into ["owner", "repo"].
	c := strings.Split(repo, "/")

	cm, _, err := cr.RepositoriesService.GetCommit(
		c[0],
		c[1],
		short,
	)
	if err != nil {
		return "", err
	}

	return *cm.SHA, nil
}

// Quayd provides a Handle method for adding a GitHub Commit Status and tagging
// the docker image.
type Quayd struct {
	StatusesRepository
	CommitResolver
}

// New returns a new Quayd instance backed by GitHub implementations.
func New(token string) *Quayd {
	t := &oauth.Transport{
		Token: &oauth.Token{AccessToken: token},
	}

	gh := github.NewClient(t.Client())

	return &Quayd{
		StatusesRepository: &GitHubStatusesRepository{gh.Repositories},
		CommitResolver:     &GitHubCommitResolver{gh.Repositories},
	}
}

// Handle resolves the ref to a full 40 character sha, then creates a new GitHub
// Commit Status for that sha.
func (q *Quayd) Handle(repo, ref, state string) error {
	sha, err := q.commitResolver().Resolve(repo, ref)
	if err != nil {
		return err
	}

	return q.statusesRepository().Create(&Status{
		Repo:    repo,
		Ref:     sha,
		State:   state,
		Context: Context,
	})
}

func (q *Quayd) commitResolver() CommitResolver {
	if q.CommitResolver == nil {
		return DefaultCommitResolver
	}

	return q.CommitResolver
}

func (q *Quayd) statusesRepository() StatusesRepository {
	if q.StatusesRepository == nil {
		return DefaultStatusesRepository
	}

	return q.StatusesRepository
}
