package helpers

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var USER_AUTH_SECRET = []byte("user123")
var ADMIN_AUTH_SECRET = []byte("admin123")

type CustomClaims struct {
	UserID uint
	jwt.MapClaims
}

func CreateUserJwtToken(id uint) string {
	claims := &CustomClaims{
		UserID: id,
		MapClaims: jwt.MapClaims{
			"sub": id,
			"exp": time.Now().Add(24 * time.Hour).Unix(),
			"iss": "shoue-shop",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	finaltoken, err := token.SignedString(USER_AUTH_SECRET)

	if err != nil {
		fmt.Println(err)
	}

	return finaltoken
}

// admin

func CreateAdminJwtToken(id uint) string {
	claims := &CustomClaims{
		UserID: id,
		MapClaims: jwt.MapClaims{
			"sub": id,
			"exp": time.Now().Add(24 * time.Hour).Unix(),
			"iss": "shoue-shop",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	finaltoken, err := token.SignedString(ADMIN_AUTH_SECRET)

	if err != nil {
		fmt.Println(err)
	}

	return finaltoken
}
