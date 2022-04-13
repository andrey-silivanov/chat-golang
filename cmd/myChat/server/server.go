package server

import (
	"github.com/andrey-silivanov/chat-golang/cmd/myChat/config"
	"github.com/andrey-silivanov/chat-golang/cmd/myChat/store"
	"github.com/andrey-silivanov/chat-golang/cmd/myChat/store/pgstore"
	"github.com/andrey-silivanov/chat-golang/pkg/database"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
)

type server struct {
	router *mux.Router
	logger *logrus.Logger
	store  store.Store
}

func newServer(store store.Store) *server {
	s := &server{
		router: mux.NewRouter(),
		logger: logrus.New(),
		store:  store,
	}
	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.HandleFunc("/register", s.registerHandler).Methods(http.MethodPost)
	s.router.HandleFunc("/auth/login", s.loginHandler).Methods(http.MethodPost)
	s.router.HandleFunc("/auth/user", s.userHandler).Methods(http.MethodGet)
	s.router.HandleFunc("/auth/refresh", s.refreshHandler).Methods(http.MethodGet)
}

func Start() error {
	appConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.GetConnection(appConfig)
	if err != nil {
		return err
	}

	defer db.Close()

	storePg := pgstore.New(db)

	srv := newServer(storePg)
	srv.logger.Info("start server")

	return http.ListenAndServe(":"+appConfig.ServerPort, srv.corsHandler(appConfig))
}

func (s *server) corsHandler(config *config.Config) http.Handler {
	credentials := handlers.AllowCredentials()
	methods := handlers.AllowedMethods([]string{http.MethodGet, http.MethodPost})
	ttl := handlers.MaxAge(3600)
	origins := handlers.AllowedOrigins([]string{config.FrontUrl})
	headers := handlers.AllowedHeaders([]string{"Authorization", "Content-Type"})
	ex := handlers.ExposedHeaders([]string{"Authorization"})

	return handlers.CORS(credentials, methods, origins, ttl, headers, ex)(s.router)
}
