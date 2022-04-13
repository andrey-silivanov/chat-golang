package handler

import (
	"github.com/andrey-silivanov/chat-golang/cmd/myChat/models"
	"github.com/andrey-silivanov/chat-golang/cmd/myChat/store"
	"log"
	"net/http"
)

func Home(store store.Store) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		u := &models.User{
			Firstname: "dsad",
			Lastname:  "2",
			Email:     "dsad",
			Password:  "dsadsad",
		}

		rep := store.GetUserRepository()

		err := rep.Create(u)
		if err != nil {
			log.Fatal(err)
		}
	}
}
