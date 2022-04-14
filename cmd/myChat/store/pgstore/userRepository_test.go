package pgstore

import (
	"github.com/andrey-silivanov/chat-golang/cmd/myChat/models"
	"github.com/andrey-silivanov/chat-golang/cmd/myChat/store"
	"github.com/andrey-silivanov/chat-golang/cmd/myChat/store/teststore"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserRepository_Create(t *testing.T) {
	st := New(teststore.DB)
	repository := st.GetUserRepository()
	u := &models.User{
		Firstname: "John",
		Lastname:  "Doe",
		Email:     "test@mail.com",
		Password:  "123456",
	}

	err := repository.Create(u)

	assert.NoError(t, err)
	assert.NotNil(t, u.Id)
}

func TestUserRepository_GetUserByFirstname(t *testing.T) {
	st := New(teststore.DB)
	repository := st.GetUserRepository()

	for _, expectedUser := range teststore.UsersFromTest {
		user, err := repository.GetUserByFirstname(expectedUser.Firstname)
		assert.Equal(t, expectedUser.Firstname, user.Firstname)
		assert.Equal(t, expectedUser.Lastname, user.Lastname)
		assert.Equal(t, expectedUser.Email, user.Email)
		assert.Nil(t, err)
	}

	emptyResult, err := repository.GetUserByFirstname("RandomFirstName")

	assert.Error(t, store.ErrRecordNotFound, err)
	assert.Nil(t, emptyResult)
}

func TestUserRepository_GetUserById(t *testing.T) {
	st := New(teststore.DB)
	repository := st.GetUserRepository()
	expectedUser := teststore.UsersFromTest[0]

	result, err := repository.GetUserById(1)

	assert.EqualValues(t, 1, result.Id)
	assert.Equal(t, expectedUser.Firstname, result.Firstname)
	assert.Nil(t, err)

	emptyResult, err := repository.GetUserById(33)
	assert.Error(t, store.ErrRecordNotFound, err)
	assert.Nil(t, emptyResult)
}

func TestUserRepository_GetUserByEmail(t *testing.T) {
	st := New(teststore.DB)
	repository := st.GetUserRepository()
	expectedUser := teststore.UsersFromTest[0]

	result, err := repository.GetUserByEmail(expectedUser.Email)

	assert.EqualValues(t, expectedUser.Email, result.Email)
	assert.Equal(t, expectedUser.Firstname, result.Firstname)
	assert.NotNil(t, expectedUser.Password)
	assert.Nil(t, err)

	emptyResult, err := repository.GetUserById(33)
	assert.Error(t, store.ErrRecordNotFound, err)
	assert.Nil(t, emptyResult)
}

func TestUserRepository_SearchUser(t *testing.T) {
	st := New(teststore.DB)
	repository := st.GetUserRepository()
	authUser := teststore.UsersFromTest[0]
	expectedUser := teststore.UsersFromTest[1]

	result, err := repository.SearchUser(expectedUser.Email, &authUser)

	for _, item := range result {
		assert.EqualValues(t, expectedUser.Email, item.Email)
		assert.Equal(t, expectedUser.Firstname, item.Firstname)
		assert.Equal(t, expectedUser.Lastname, item.Lastname)
	}

	assert.Nil(t, err)

	emptyResult, err := repository.SearchUser("not_found@gmail.com", &authUser)
	assert.Error(t, store.ErrRecordNotFound, err)
	assert.Nil(t, emptyResult)
}
