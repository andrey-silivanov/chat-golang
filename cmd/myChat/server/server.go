package server

import (
	"context"
	"github.com/andrey-silivanov/chat-golang/cmd/myChat/config"
	"github.com/andrey-silivanov/chat-golang/cmd/myChat/store"
	"github.com/andrey-silivanov/chat-golang/cmd/myChat/store/pgstore"
	"github.com/andrey-silivanov/chat-golang/pkg/database"
	"github.com/andrey-silivanov/chat-golang/pkg/jwtToken"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"strings"
)

const (
	ctxKeyUser ctxKey = iota
)

type ctxKey uint8

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
	//	s.router.HandleFunc("/auth/user", s.userHandler).Methods(http.MethodGet)
	s.router.HandleFunc("/auth/refresh", s.refreshHandler).Methods(http.MethodGet)

	private := s.router.PathPrefix("/").Subrouter()
	private.Use(s.authenticateUserMiddleware)
	private.HandleFunc("/auth/user", s.userHandler).Methods(http.MethodGet)
	//private.HandleFunc("/auth/refresh", s.refreshHandler).Methods(http.MethodGet)
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

func (s *server) authenticateUserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		repository := s.store.GetUserRepository()

		user, err := repository.GetUserByEmail(claims.Email)
		if err != nil {
			UnauthorizedError(w, http.StatusText(http.StatusUnauthorized))

			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, user)))
	})
}
