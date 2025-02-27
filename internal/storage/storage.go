package storage

import (
	"github.com/philipslstwoyears/calculator-go/internal/dto"
	"sort"
	"sync"
)

type Storage struct {
	data   map[int]dto.Expression
	mu     *sync.RWMutex
	lastID int
}

func New() *Storage {
	return &Storage{
		data: make(map[int]dto.Expression),
		mu:   &sync.RWMutex{},
	}
}
func (s *Storage) Add(e dto.Expression) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	e.Id = s.lastID
	s.data[e.Id] = e
	s.lastID++
	return e.Id
}
func (s *Storage) Update(e dto.Expression) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[e.Id] = e
}

func (s *Storage) Get(id int) (dto.Expression, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.data[id]
	return e, ok
}
func (s *Storage) GetAll() []dto.Expression {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var slice []dto.Expression
	for _, e := range s.data {
		slice = append(slice, e)
	}
	sort.Slice(slice, func(i, j int) bool {
		return slice[i].Id < slice[j].Id
	})
	return slice
}
