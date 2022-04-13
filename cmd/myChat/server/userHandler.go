package server

import (
	"github.com/andrey-silivanov/chat-golang/cmd/myChat/models"
	"net/http"
)

func (s *server) userHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(ctxKeyUser).(*models.User)
	response(w, user)
}
