package server

import (
	"errors"
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
		StatusInternalServerError(w)
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

//func welcome(w http.ResponseWriter, r *http.Request) {
//
//	tokenStr, err := getTokenFromHeader(r)
//	if err != nil {
//		httpErrors.BadRequest(w, err.Error())
//	}
//	fmt.Println("reqToken", tokenStr)
//
//	// Initialize a new instance of `Claims`
//
//	// Parse the JWT string and store the result in `claims`.
//	// Note that we are passing the key in this method as well. This method will return an error
//	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
//	// or if the signature does not match
//	claims, err := jwtToken.ParseJWTToken(tokenStr)
//
//	if err != nil {
//		if err == jwt.ErrSignatureInvalid {
//			httpErrors.UnauthorizedError(w)
//			return
//		}
//		httpErrors.BadRequest(w, err.Error())
//		return
//	}
//
//	user, ok := getUser(claims.FirstName)
//	if !ok {
//		httpErrors.UnauthorizedError(w) // @TODO поменять ошибку
//
//		return
//	}
//
//	httpResponse.ToJson(w, user)
//}

func getTokenFromHeader(r *http.Request) (string, error) {
	reqToken, ok := r.Header["Authorization"]
	if !ok {
		return "", errors.New("header not set")
	}

	splitToken := strings.Split(reqToken[0], "Bearer ")

	return splitToken[1], nil
}

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
		StatusInternalServerError(w)
		return
	}

	w.Header().Set("Authorization", "Bearer "+tokenString)

	response(w, user)

}

//func getUser(username string) (models.User, bool) {
//	dbConn := database.GetConnection()
//	defer dbConn.Close()
//
//	userRepository := repository.CreateUserRepository(dbConn)
//
//	return userRepository.GetUserByFirstname(username)
//}

func checkPassword(requestBody loginRequestBody, user *models.User) bool {

	return requestBody.Password == user.Password
}
