package pgstore

import (
	"database/sql"
	"github.com/andrey-silivanov/chat-golang/cmd/myChat/store"
)

type Store struct {
	db             *sql.DB
	userRepository *UserRepository
}

func (s *Store) GetUserRepository() store.UserRepository {
	if s.userRepository == nil {
		s.userRepository = &UserRepository{db: s.db}
	}

	return s.userRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}
