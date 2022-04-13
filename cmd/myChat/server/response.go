package server

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	HTTPCode int         `json:"-"`
	Code     int         `json:"code"`
	Data     interface{} `json:"data"`
}

func response(w http.ResponseWriter, data interface{}) {
	r := &Response{
		HTTPCode: http.StatusOK,
		Code:     http.StatusOK,
		Data:     data,
	}

	toJson(w, r)
}

func toJson(w http.ResponseWriter, r *Response) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(r.HTTPCode)
	json.NewEncoder(w).Encode(r)
}
