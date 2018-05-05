package kvstore

import "sync"

type Storage struct {
	Store map[string]HashMap
	mu    sync.Mutex
}

// WARNING: It is only for debugging purposes!
func NewInMemory() *Storage {
	return &Storage{
		Store: map[string]HashMap{},
	}
}

func (s *Storage) NewHash(ident string) HashMap {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Allocate new hash-map
	s.Store[ident] = HashMap{Data: map[string]string{}}

	return s.Store[ident]
}

func (s *Storage) DeleteHash(ident string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.Store, ident)
}
