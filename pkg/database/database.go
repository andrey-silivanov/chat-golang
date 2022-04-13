package database

import (
	"database/sql"
	"github.com/andrey-silivanov/chat-golang/cmd/myChat/config"
	_ "github.com/lib/pq"
	"log"
)

func GetConnection(config *config.Config) (*sql.DB, error) {

	log.Println("Connect db to ", config.DbUrl)

	db, err := sql.Open("postgres", config.DbUrl)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
