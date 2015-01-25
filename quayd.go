package quayd

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

// NewStatusesRepository returns a new StatusesRepository implementation.
func NewStatusesRepository(kind string) StatusesRepository {
	switch kind {
	case "fake":
		return &FakeStatusesRepository{}
	default:
		panic("Not a valid statuses repository.")
	}
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

// CommitResolver is an interface for resolving a short sha to a full 40
// character sha.
type CommitResolver interface {
	Resolve(short string) (string, error)
}

// FakeCommitResolver just returns the short sha.
type FakeCommitResolver struct{}

// Resolve implements CommitResolver Resolve.
func (cr *FakeCommitResolver) Resolve(short string) (string, error) {
	return short, nil
}

// StatusesService provides a convenient server for creating new commit
// statuses.
type StatusesService struct {
	StatusesRepository
	CommitResolver
}

func (s *StatusesService) Create(repo, ref, state string) error {
	sha, err := s.Resolve(ref)
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
