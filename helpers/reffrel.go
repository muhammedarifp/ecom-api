package helpers

import (
	"encoding/base64"
	"math/rand"
)

func CreateReffrerelCode() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)

	reffrelcode := base64.URLEncoding.EncodeToString(bytes)[:8]
	return reffrelcode
}
