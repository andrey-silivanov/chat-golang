package store

type Store interface {
	GetUserRepository() UserRepository
}
