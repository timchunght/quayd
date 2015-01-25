package quayd

var (
	DefaultStatusesRepository = &statusesRepository{}
	DefaultCommitResolver     = &commitResolver{}
	DefaultStatusesService    = &StatusesService{
		StatusesRepository: DefaultStatusesRepository,
		CommitResolver:     DefaultCommitResolver,
	}
)

// StatusesRepository is an interface that can be implemented for creating
// Commit Statuses.
type StatusesRepository interface {
	Create(*Status) error
}

type Status struct {
	Repo    string
	Ref     string
	State   string
	Context string
}

// statusesRepository is a fake implementation of the StatusesRepository
// interface.
type statusesRepository struct {
	s []*Status
}

func (r *statusesRepository) Create(status *Status) error {
	r.s = append(r.s, status)

	return nil
}

// CommitResolver is an interface for resolving a short sha to a full 40
// character sha.
type CommitResolver interface {
	Resolve(short string) (string, error)
}

// commitResolver just returns the short sha.
type commitResolver struct{}

func (c *commitResolver) Resolve(short string) (string, error) {
	return short, nil
}

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
