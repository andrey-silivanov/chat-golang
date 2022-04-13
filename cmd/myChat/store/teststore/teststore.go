package teststore

import (
	"database/sql"
	"fmt"
	"github.com/andrey-silivanov/chat-golang/cmd/myChat/models"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const MIGRATION_PATH = "/var/www/go/chat/migrations/"

var DB *sql.DB

var UsersFromTest = []models.User{
	{
		Firstname: "John",
		Lastname:  "Doe",
		Email:     "test@mail.com",
		Password:  "123456",
	},
	{
		Firstname: "Mary",
		Lastname:  "Jones",
		Email:     "mary@mail.com",
		Password:  "123456",
	},
}

func CreateTestDB(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "13.3",
		Env: []string{
			"POSTGRES_PASSWORD=secret",
			"POSTGRES_USER=user_name",
			"POSTGRES_DB=dbname",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	/* @TODO переделать через env */
	databaseUrl := fmt.Sprintf("postgres://user_name:secret@%s/dbname?sslmode=disable", hostAndPort)

	log.Println("Connecting to database on url: ", databaseUrl)

	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		DB, err = sql.Open("postgres", databaseUrl)
		if err != nil {
			return err
		}
		return DB.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	err = migrate(DB)
	seedUsers(UsersFromTest)

	if err != nil {
		log.Fatalf("Migration error: %s", err)
	}

	//Run tests
	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func migrate(db *sql.DB) error {

	files, err := filepath.Glob(MIGRATION_PATH + "*_up.sql")
	if err != nil {
		return err
	}
	fmt.Println(files)
	for _, file := range files {
		fContent, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}
		err = executeSql(db, string(fContent))
		if err != nil {
			return err
		}
	}

	return nil
}

func executeSql(db *sql.DB, sqlQuery string) error {
	_, err := db.Exec(sqlQuery)
	if err != nil {
		return err
	}

	return nil
}

func seedUsers(users []models.User) {

	valueStrings := make([]string, 0, len(users))
	valueArgs := make([]interface{}, 0, len(users)*4)
	i := 0
	for _, u := range users {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d)", i*4+1, i*4+2, i*4+3, i*4+4))
		valueArgs = append(valueArgs, u.Firstname)
		valueArgs = append(valueArgs, u.Lastname)
		valueArgs = append(valueArgs, u.Email)
		valueArgs = append(valueArgs, u.Password)
		i++
	}
	stmt := fmt.Sprintf("INSERT INTO users(firstname, lastname, email, password) VALUES %s", strings.Join(valueStrings, ","))
	_, err := DB.Exec(stmt, valueArgs...)
	if err != nil {
		log.Fatal(err)
	}
}
