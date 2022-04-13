package server

import (
	"fmt"
	"github.com/andrey-silivanov/chat-golang/pkg/jwtToken"
	"net/http"
	"strings"
)

func (s *server) userHandler(w http.ResponseWriter, r *http.Request) {
	reqToken, ok := r.Header["Authorization"]
	if !ok {
		UnauthorizedError(w, http.StatusText(http.StatusUnauthorized))

		return
	}
	fmt.Println(1111, reqToken)

	splitToken := strings.Split(reqToken[0], "Bearer ")
	tknStr := splitToken[1]

	claims, err := jwtToken.ParseJWTToken(tknStr)

	if err != nil {
		UnauthorizedError(w, err.Error())

		return
	}

	repository := s.store.GetUserRepository()

	user, err := repository.GetUserByEmail(claims.Email)
	if err != nil {
		UnauthorizedError(w, http.StatusText(http.StatusUnauthorized))

		return
	}

	response(w, user)
}
