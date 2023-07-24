package helpers

import (
	"github.com/muhammedarif/Ecomapi/config"
	"github.com/muhammedarif/Ecomapi/models"
)

// If user banned return false
// Otherwise return true
func IsUserBlocked(userID string) bool {
	db := *config.GetDb()
	var userData models.Users
	db.First(&userData, userData)
	if userData.Status {
		return true
	} else {
		return false
	}
}
