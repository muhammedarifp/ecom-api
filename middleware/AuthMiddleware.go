package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/muhammedarif/Ecomapi/helpers"
)

func AuthUserMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authToken := ctx.GetHeader("token")
		if authToken == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "permision denaid",
				"User":  true,
			})
			return
		}

		_, err := jwt.Parse(authToken, func(t *jwt.Token) (interface{}, error) {
			return helpers.USER_AUTH_SECRET, nil
		})

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":  "permision denaid",
				"socond": err.Error(),
			})
			return
		}

		ctx.Next()
	}
}

// Admin

func AuthAdminMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authToken := ctx.GetHeader("token")
		if authToken == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "permision denaid",
			})
			return
		}

		_, err := jwt.Parse(authToken, func(t *jwt.Token) (interface{}, error) {
			return helpers.ADMIN_AUTH_SECRET, nil
		})

		if err != nil {
			fmt.Println(err.Error())
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": false,
				"error":  "permision denaid, you are not a admin",
			})
			return
		}

		ctx.Next()
	}
}
