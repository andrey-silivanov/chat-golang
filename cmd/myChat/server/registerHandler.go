package server

import (
	"github.com/andrey-silivanov/chat-golang/cmd/myChat/models"
	"github.com/thedevsaddam/govalidator"
	"net/http"
)

type registerRequestBody struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (s *server) registerHandler(w http.ResponseWriter, r *http.Request) {
	requestBody, validationErrors := validationRegisterRequest(r)

	if len(validationErrors) != 0 {
		ValidationError(w, validationErrors)

		return
	}

	userRepository := s.store.GetUserRepository()

	u := &models.User{
		Firstname: requestBody.Firstname,
		Lastname:  requestBody.Lastname,
		Email:     requestBody.Email,
		Password:  requestBody.Password,
	}

	err := userRepository.Create(u)
	if err != nil {
		StatusInternalServerError(w)
	}

	response(w, u)
}

func validationRegisterRequest(r *http.Request) (registerRequestBody, map[string][]string) {
	var requestBody registerRequestBody

	rules := govalidator.MapData{
		"firstname": []string{"required", "between:3,50"},
		"lastname":  []string{"required", "between:3,50"},
		"email":     []string{"required", "min:4", "max:100", "email"},
		"password":  []string{"required", "min:4", "max:20"},
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
