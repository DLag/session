package session

type Session struct {
	name    string
	manager *Manager
}

// Set sets variable in Session
func (s *Session) Set(name string, value interface{}) error {
	if s.manager == nil {
		return ErrorEmptyManager
	}
	return s.manager.Set(s.name, name, value)
}

// Get gets variable from Session
func (s *Session) Get(name string) (interface{}, error) {
	if s.manager == nil {
		return nil, ErrorEmptyManager
	}
	return s.manager.Get(s.name, name)
}

// Name gets name of the Session
func (s *Session) Name() string {
	return s.name
}
