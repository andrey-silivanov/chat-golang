package server

import (
	"github.com/andrey-silivanov/chat-golang/cmd/myChat/models"
	"github.com/andrey-silivanov/chat-golang/pkg/jwtToken"
	"github.com/dgrijalva/jwt-go"
	"github.com/thedevsaddam/govalidator"
	"net/http"
	"strings"
	"time"
)

type loginRequestBody struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

func (s *server) loginHandler(w http.ResponseWriter, r *http.Request) {

	requestBody, validationErrors := validationLoginRequest(r)

	if len(validationErrors) != 0 {
		ValidationError(w, validationErrors)

		return
	}

	userRepository := s.store.GetUserRepository()

	user, err := userRepository.GetUserByEmail(requestBody.Username)

	if err != nil {
		UnauthorizedError(w, "Password or username invalid")

		return
	}

	if !checkPassword(requestBody, user) {
		UnauthorizedError(w, "Password or username invalid")

		return
	}

	tokenString, err := jwtToken.GenerateJWTToken(user)
	if err != nil {
		StatusInternalServerError(w, err.Error())
		return
	}

	w.Header().Set("Authorization", "Bearer "+tokenString)

	response(w, user)
}

func validationLoginRequest(r *http.Request) (loginRequestBody, map[string][]string) {
	var requestBody loginRequestBody

	rules := govalidator.MapData{
		"username": []string{"required", "between:3,50"},
		"password": []string{"required", "min:4", "max:20"},
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

/*@ TODO переделать логику обновления токена */
func (s *server) refreshHandler(w http.ResponseWriter, r *http.Request) {

	header := r.Header.Get("Authorization")
	bearerToken := strings.Split(header, " ")
	if len(bearerToken) != 2 {
		BadRequestError(w, "Cannot read token")

		return
	}

	token := bearerToken[1]

	if bearerToken[0] != "Bearer" || len(token) == 0 {
		BadRequestError(w, "Error in authorization token. it needs to be in form of 'Bearer <token>'")

		return
	}

	claims, err := jwtToken.ParseJWTToken(token)
	if err != nil {

		if err == jwt.ErrSignatureInvalid {
			BadRequestError(w, err.Error())
			return
		}
		UnauthorizedError(w, err.Error())
		return
	}

	//// (END) The code up-till this point is the same as the first part of the `Welcome` route
	//
	//// We ensure that a new token is not issued until enough time has elapsed
	//// In this case, a new token will only be issued if the old token is within
	//// 30 seconds of expiry. Otherwise, return a bad request status
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 60*time.Second {

		return
	}

	userRepository := s.store.GetUserRepository()

	user, err := userRepository.GetUserByEmail(claims.Email)
	if err != nil {
		UnauthorizedError(w, err.Error())

		return
	}
	tokenString, err := jwtToken.GenerateJWTToken(user)
	if err != nil {
		StatusInternalServerError(w, err.Error())
		return
	}

	w.Header().Set("Authorization", "Bearer "+tokenString)

	response(w, user)

}

func checkPassword(requestBody loginRequestBody, user *models.User) bool {

	return requestBody.Password == user.Password
}
