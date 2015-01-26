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

	// DefaultTagger is the default Tagger to use.
	DefaultTagger = &tagger{}

	// DefaultTagResovler is the default TagResolver to use.
	DefaultTagResolver = &tagResolver{}

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
	// Create creates a GitHub Commit Status.
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
	// Resolve resolves the short sha to a full 40 character sha.
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

// Tagger is an interface for tagging a docker image with a tag.
type Tagger interface {
	// Tag tags the imageID with the given tag.
	Tag(imageID, tag string) error
}

// tagger is a fake implementation of the Tagger interface.
type tagger struct {
	tags map[string]string
}

// Tag implements Tagger Tag.
func (t *tagger) Tag(imageID, tag string) error {
	if t.tags == nil {
		t.tags = make(map[string]string)
	}

	t.tags[imageID] = tag

	return nil
}

// TagResolver resolves a docker tag to an image id.
type TagResolver interface {
	Resolve(tag string) (string, error)
}

// tagResolver is a fake implementation of the TagResolver interface.
type tagResolver struct{}

func (r *tagResolver) Resolve(tag string) (string, error) {
	return "1234", nil
}

// DockerTagResolver is an implementation of the TagResolver that resolves an
// image tag to a docker image id, using the docker api.
type DockerTagResolver struct {
	// TODO
}

func (r *DockerTagResolver) Resolve(tag string) (string, error) {
	// TODO Something with the docker api
	return "", nil
}

// Quayd provides a Handle method for adding a GitHub Commit Status and tagging
// the docker image.
type Quayd struct {
	StatusesRepository
	CommitResolver
	Tagger
	TagResolver
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

// TagImage takes a docker tag, resolves it to a image id, then tags it with the
// given tag.
func (q *Quayd) TagImage(tag, repo, ref string) error {
	sha, err := q.commitResolver().Resolve(repo, ref)
	if err != nil {
		return err
	}

	// Something that resolves the `tag` into an image id.
	imageID, err := q.tagResolver().Resolve(tag)
	if err != nil {
		return err
	}

	return q.tagger().Tag(imageID, sha)
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

func (q *Quayd) tagger() Tagger {
	if q.Tagger == nil {
		q.Tagger = DefaultTagger
	}

	return q.Tagger
}

func (q *Quayd) tagResolver() TagResolver {
	if q.TagResolver == nil {
		q.TagResolver = DefaultTagResolver
	}

	return q.TagResolver
}
