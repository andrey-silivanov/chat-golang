package server

import "net/http"

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*") // @TODO изменить на конкретный url
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
