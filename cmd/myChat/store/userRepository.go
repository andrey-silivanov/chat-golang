package store

import "github.com/andrey-silivanov/chat-golang/cmd/myChat/models"

type UserRepository interface {
	Create(u *models.User) error
	GetUserByFirstname(firstname string) (*models.User, error)
	GetUserById(id int) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	SearchUser(email string, excludedUser *models.User) ([]models.User, error)
}
