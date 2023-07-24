package helpers

import (
	"golang.org/x/crypto/bcrypt"
)

func PassToHash(pass string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		panic("Hash err")
	}

	return string(hash)
}

func CompareHashPass(orgpass string, hashpass string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashpass), []byte(orgpass)); err != nil {
		return false
	}

	return true
}
