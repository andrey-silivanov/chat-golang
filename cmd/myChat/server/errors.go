package server

import "net/http"

func ValidationError(w http.ResponseWriter, msg map[string][]string) {

	r := &Response{
		HTTPCode: http.StatusUnprocessableEntity,
		Code:     http.StatusUnprocessableEntity,
		Data:     msg,
	}

	toJson(w, r)
}

func UnauthorizedError(w http.ResponseWriter, message string) {
	r := &Response{
		HTTPCode: http.StatusUnauthorized,
		Code:     http.StatusUnauthorized,
		Data:     message,
	}

	toJson(w, r)
}

func StatusInternalServerError(w http.ResponseWriter, msg string) {
	r := &Response{
		HTTPCode: http.StatusInternalServerError,
		Code:     http.StatusInternalServerError,
		Data:     msg,
	}

	toJson(w, r)
}

func BadRequestError(w http.ResponseWriter, msg string) {
	r := &Response{
		HTTPCode: http.StatusBadRequest,
		Code:     http.StatusBadRequest,
		Data:     msg,
	}

	toJson(w, r)
}
