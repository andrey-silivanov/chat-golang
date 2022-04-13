package pgstore

import (
	"database/sql"
	"github.com/andrey-silivanov/chat-golang/cmd/myChat/models"
	"github.com/andrey-silivanov/chat-golang/cmd/myChat/store"
)

type UserRepository struct {
	db *sql.DB
}

func (r *UserRepository) Create(u *models.User) error {

	sqlStatement := "INSERT INTO users (firstname, lastname, email, password) VALUES ($1, $2, $3, $4) RETURNING id"

	err := r.db.QueryRow(sqlStatement,
		u.Firstname,
		u.Lastname,
		u.Email,
		u.Password,
	).Scan(&u.Id)

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetUserByFirstname(firstname string) (*models.User, error) {
	result := &models.User{}

	query := "SELECT id, firstname, lastname, email FROM users where firstname = $1"
	row := r.db.QueryRow(query, firstname)
	err := row.Scan(&result.Id, &result.Firstname, &result.Lastname, &result.Email)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return result, nil
}

func (r *UserRepository) GetUserById(id int) (*models.User, error) {
	result := &models.User{}

	query := "SELECT id, firstname, lastname, email FROM users where id = $1"
	row := r.db.QueryRow(query, id)

	err := row.Scan(&result.Id, &result.Firstname, &result.Lastname, &result.Email)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return result, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	result := &models.User{}

	query := "SELECT id, firstname, lastname, email, password FROM users where email = $1"
	row := r.db.QueryRow(query, email)

	err := row.Scan(&result.Id, &result.Firstname, &result.Lastname, &result.Email, &result.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return result, nil
}
