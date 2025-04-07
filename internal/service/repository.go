package service

import (
	"fmt"
	"math/rand"
)

type Round struct {
	SpyID   int
	Place   Place
	Members int
	Roles   []int
	seed    *rand.Rand
}

func (r Round) Equal(other Round) bool {
	if r.SpyID != other.SpyID || r.Members != other.Members {
		return false
	}
	if r.Place != other.Place {
		return false
	}
	if len(r.Roles) != len(other.Roles) {
		return false
	}
	for i := range r.Roles {
		if r.Roles[i] != other.Roles[i] {
			return false
		}
	}
	return true
}

func (r Round) String() string {
	return fmt.Sprintf("SpyID: %d, Place: %s, Members: %d, Roles: %v", r.SpyID, r.Place.name, r.Members, r.Roles)
}

type Store struct {
	chats map[int64]Round
}

func NewStore() *Store {
	return &Store{make(map[int64]Round)}
}

func (s *Store) Get(chaID int64) (Round, bool) {
	round, ok := s.chats[chaID]
	return round, ok
}

func (s *Store) Set(chaID int64, round Round) {
	s.chats[chaID] = round
}
