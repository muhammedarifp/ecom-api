package admincontroller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/muhammedarif/Ecomapi/config"
	"github.com/muhammedarif/Ecomapi/helpers"
	"github.com/muhammedarif/Ecomapi/models"
)

type AdminLoginInp struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=5"`
}

func AdminLoginController() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var loginInp AdminLoginInp
		if err := ctx.ShouldBindJSON(&loginInp); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Error": err.Error(),
			})
			return
		}

		validator := validator.New()
		if err := validator.Struct(&loginInp); err != nil {
			ctx.AbortWithStatusJSON(400, gin.H{
				"Error": err.Error(),
			})
			return
		}

		db := *config.GetDb()
		var adminDeta models.Users
		if res := db.First(&adminDeta, `email = ?`, loginInp.Email); res.Error != nil {
			ctx.JSON(http.StatusInternalServerError, models.Response{
				Status:  false,
				Message: "Admin login failed",
				Error:   res.Error.Error(),
			})
			return
		}

		if !adminDeta.Isadmin {
			ctx.JSON(http.StatusBadRequest, models.Response{
				Status:  false,
				Message: "Admin login failed",
				Error:   "You are not a admin",
			})
			return
		}

		// Check user password
		if !helpers.CompareHashPass(loginInp.Password, adminDeta.Password) {
			ctx.JSON(http.StatusBadRequest, models.Response{
				Status:  false,
				Message: "Admin login failed",
				Error:   "your password is incorrect",
			})
			return
		}

		authToken := helpers.CreateAdminJwtToken(adminDeta.ID)

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Status":  true,
			"Message": "Admin login success",
			"Token":   authToken,
			"Error":   nil,
		})
	}
}
