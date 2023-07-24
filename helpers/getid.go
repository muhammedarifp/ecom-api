package helpers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func GetUserIDFromJwt(ctx *gin.Context) string {
	token := ctx.GetHeader("token")

	res, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return USER_AUTH_SECRET, nil
	})

	if err != nil {
		return ""
	} else if !res.Valid {
		return ""
	} else {
		userID := res.Claims.(jwt.MapClaims)["UserID"]
		userIDstr := fmt.Sprintf("%v", userID)
		return userIDstr
	}
}
