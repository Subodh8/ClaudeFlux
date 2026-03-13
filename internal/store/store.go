package store

type Store struct{}

func Open(dir string) (*Store, error) {
	return &Store{}, nil
}

func (s *Store) Close() error {
	return nil
}

func (s *Store) CreateRun(runID, name string) error { return nil }
func (s *Store) FailRun(runID, err string) error    { return nil }
func (s *Store) CompleteRun(runID string) error     { return nil }

func (s *Store) CreateAgentRun(runID, name string) error               { return nil }
func (s *Store) FailAgentRun(runID, name, err string) error            { return nil }
func (s *Store) CompleteAgentRun(runID, name string, result any) error { return nil }
