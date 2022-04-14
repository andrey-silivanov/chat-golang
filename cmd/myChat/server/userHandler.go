package server

import (
	"github.com/andrey-silivanov/chat-golang/cmd/myChat/models"
	"github.com/thedevsaddam/govalidator"
	"net/http"
)

func (s *server) userHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(ctxKeyUser).(*models.User)
	response(w, user)
}

type searchUserRequestBody struct {
	Username string `json:"username"`
}

func (s *server) searchUserHandler(w http.ResponseWriter, r *http.Request) {
	authUser := r.Context().Value(ctxKeyUser).(*models.User)

	requestBody, validationErrors := validationUserSearchRequest(r)

	if len(validationErrors) != 0 {
		ValidationError(w, validationErrors)

		return
	}

	userRepository := s.store.GetUserRepository()
	users, err := userRepository.SearchUser(requestBody.Username, authUser)

	if err != nil {
		StatusInternalServerError(w, err.Error())
		return
	}

	response(w, users)
}

func validationUserSearchRequest(r *http.Request) (searchUserRequestBody, map[string][]string) {
	var requestBody searchUserRequestBody

	rules := govalidator.MapData{
		"username": []string{"required", "email"},
	}
	opts := govalidator.Options{
		Request: r,
		Data:    &requestBody,
		Rules:   rules,
	}

	v := govalidator.New(opts)
	validationErrors := v.ValidateJSON()

	return requestBody, validationErrors
}
