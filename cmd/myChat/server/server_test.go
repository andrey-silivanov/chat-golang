package server

import (
	"bytes"
	"encoding/json"
	"github.com/andrey-silivanov/chat-golang/cmd/myChat/models"
	"github.com/andrey-silivanov/chat-golang/cmd/myChat/store/pgstore"
	"github.com/andrey-silivanov/chat-golang/cmd/myChat/store/teststore"
	"github.com/andrey-silivanov/chat-golang/pkg/jwtToken"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	teststore.CreateTestDB(m)
}

func TestServer_Register(t *testing.T) {
	store := pgstore.New(teststore.DB)

	testCases := []struct {
		name         string
		payload      map[string]string
		expectedCode int
	}{
		{
			name: "register success",
			payload: map[string]string{
				"firstname": "John",
				"lastname":  "Smith",
				"email":     "john_smith@mail.com",
				"password":  "123456",
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "empty firstname",
			payload: map[string]string{
				"firstname": "",
				"lastname":  "Smith",
				"email":     "john_smith@mail.com",
				"password":  "123456",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "empty lastname",
			payload: map[string]string{
				"firstname": "John",
				"lastname":  "",
				"email":     "john_smith@mail.com",
				"password":  "123456",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "empty email",
			payload: map[string]string{
				"firstname": "John",
				"lastname":  "Smith",
				"email":     "",
				"password":  "123456",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid email",
			payload: map[string]string{
				"firstname": "John",
				"lastname":  "Smith",
				"email":     "invalid",
				"password":  "123456",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "empty password",
			payload: map[string]string{
				"firstname": "John",
				"lastname":  "Smith",
				"email":     "john_smith@mail.com",
				"password":  "",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
	}

	s := newServer(store)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/register", b)

			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}

}

func TestServer_Login(t *testing.T) {
	store := pgstore.New(teststore.DB)
	u := teststore.UsersFromTest[0]

	testCases := []struct {
		name                string
		payload             map[string]interface{}
		authorizationHeader bool
		expectedCode        int
	}{
		{
			name: "authenticated",
			payload: map[string]interface{}{
				"username": u.Email,
				"password": u.Password,
			},
			authorizationHeader: true,
			expectedCode:        http.StatusOK,
		},
		{
			name: "username missing",
			payload: map[string]interface{}{
				"username": "",
				"password": u.Password,
			},
			authorizationHeader: false,
			expectedCode:        http.StatusUnprocessableEntity,
		},
		{
			name: "password missing",
			payload: map[string]interface{}{
				"username": u.Email,
				"password": "",
			},
			authorizationHeader: false,
			expectedCode:        http.StatusUnprocessableEntity,
		},
		{
			name: "not authenticated",
			payload: map[string]interface{}{
				"username": "invalid@mail.com",
				"password": "654321",
			},
			authorizationHeader: false,
			expectedCode:        http.StatusUnauthorized,
		},
	}

	s := newServer(store)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/auth/login", b)

			s.ServeHTTP(rec, req)

			headers := rec.Header()
			response := &Response{}
			if tc.authorizationHeader == true {
				response.Data = &models.User{}
			}

			err := json.NewDecoder(rec.Body).Decode(response)
			if err != nil {
				log.Fatal(err)
			}

			authHeader, ok := headers["Authorization"]

			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.Equal(t, tc.authorizationHeader, ok)
			if ok {
				assert.NotNil(t, authHeader)
				userFromResponse := response.Data.(*models.User)
				assert.Equal(t, tc.payload["username"], userFromResponse.Email)
			}
		})
	}
}

func TestServer_User(t *testing.T) {
	store := pgstore.New(teststore.DB)
	u := teststore.UsersFromTest[0]

	testCases := []struct {
		name                string
		authorizationHeader bool
		validToken          bool
		expectedCode        int
	}{
		{
			name:                "authenticated user",
			authorizationHeader: true,
			validToken:          true,
			expectedCode:        http.StatusOK,
		},
		{
			name:                "not authenticated user",
			authorizationHeader: true,
			validToken:          false,
			expectedCode:        http.StatusUnauthorized,
		},
		{
			name:                "not set header",
			authorizationHeader: false,
			validToken:          false,
			expectedCode:        http.StatusBadRequest,
		},
	}

	s := newServer(store)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/auth/user", nil)
			if tc.authorizationHeader {
				var token string
				if tc.validToken {
					token, _ = jwtToken.GenerateJWTToken(&u)
				} else {
					expirationTime := time.Now().Add(-5 * time.Minute)
					token, _ = jwtToken.GenerateJWTTokenWithExpirationTime(&u, expirationTime)
				}

				req.Header.Set("Authorization", "Bearer "+token)
			}

			s.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)

			if tc.authorizationHeader && tc.validToken {
				response := &Response{
					Data: &models.User{},
				}

				err := json.NewDecoder(rec.Body).Decode(response)
				if err != nil {
					log.Fatal(err)
				}
				userFromResponse := response.Data.(*models.User)
				assert.Equal(t, u.Email, userFromResponse.Email)
			}
		})
	}
}

func TestServer_Refresh(t *testing.T) {
	store := pgstore.New(teststore.DB)
	u := teststore.UsersFromTest[0]

	testCases := []struct {
		name               string
		authorizationToken func() string
		authHeader         bool
		expectedCode       int
	}{
		{
			name: "don't refresh token",
			authorizationToken: func() string {
				expirationTime := time.Now().Add(5 * time.Minute)

				token, _ := jwtToken.GenerateJWTTokenWithExpirationTime(&u, expirationTime)
				return token
			},
			authHeader:   true,
			expectedCode: http.StatusOK,
		},
		{
			name: "refresh token",
			authorizationToken: func() string {
				expirationTime := time.Now().Add(1 * time.Minute)

				token, _ := jwtToken.GenerateJWTTokenWithExpirationTime(&u, expirationTime)
				return token
			},
			authHeader:   true,
			expectedCode: http.StatusOK,
		},
		{
			name: "expired token",
			authorizationToken: func() string {
				expirationTime := time.Now().Add(-5 * time.Minute)

				token, _ := jwtToken.GenerateJWTTokenWithExpirationTime(&u, expirationTime)
				return token
			},
			authHeader:   true,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "empty token",
			authorizationToken: func() string {
				return ""
			},
			authHeader:   true,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "header not set",
			authorizationToken: func() string {
				return ""
			},
			authHeader:   false,
			expectedCode: http.StatusBadRequest,
		},
	}

	s := newServer(store)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/auth/refresh", nil)
			if tc.authHeader {
				req.Header.Set("Authorization", "Bearer "+tc.authorizationToken())
			}

			s.ServeHTTP(rec, req)

			headers := rec.Header()
			authHeader, ok := headers["Authorization"]

			if tc.name == "refresh token" {
				assert.True(t, ok)
				assert.NotNil(t, authHeader)
			} else {
				assert.False(t, ok)
				assert.Nil(t, authHeader)
			}

			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
