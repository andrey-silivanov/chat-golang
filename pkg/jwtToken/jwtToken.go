package jwtToken

import (
	"errors"
	"github.com/andrey-silivanov/chat-golang/cmd/myChat/models"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// Create the JWT key used to create the signature
var jwtKey = []byte("my_secret_key")

// Claims Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	FirstName string `json:"username"`
	Email     string `json:"email"`
	jwt.StandardClaims
}

func GenerateJWTToken(user *models.User) (string, error) {
	// Declare the expiration time of the token
	expirationTime := time.Now().Add(60 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	return GenerateJWTTokenWithExpirationTime(user, expirationTime)
}

func ParseJWTToken(token string) (claims *Claims, err error) {
	claims = &Claims{}

	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return
	}

	if !tkn.Valid {
		err = errors.New("token is not valid")
		return
	}

	return
}

func GenerateJWTTokenWithExpirationTime(user *models.User, expirationTime time.Time) (string, error) {
	claims := getClaims(user, expirationTime)
	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the JWT string
	return token.SignedString(jwtKey)
}

func getClaims(u *models.User, expirationTime time.Time) *Claims {
	return &Claims{
		FirstName: u.Firstname,
		Email:     u.Email,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}
}
