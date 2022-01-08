package app

type Store struct {
	apps []Meta
}

func (s *Store) Add(app Meta) {
	s.apps = append(s.apps, app)
}

func (s *Store) List() []Meta {
	return s.apps
}
