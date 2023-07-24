package usercontroller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/muhammedarif/Ecomapi/config"
	"github.com/muhammedarif/Ecomapi/helpers"
	"github.com/muhammedarif/Ecomapi/models"
)

func UserVerifyOtpController() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userid := ctx.PostForm("id")
		otp := helpers.GenarateOtp()
		if res, err := helpers.SendOtpEmail(userid, otp); !res {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
				Status:  false,
				Message: "User valifation error",
				Error:   err.Error(),
			})

			return
		}

		ctx.AbortWithStatusJSON(http.StatusOK, models.Response{
			Status:  true,
			Message: "Successfully sended",
			Error:   nil,
		})
	}
}

func UserVarifyValidateController() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userid := ctx.PostForm("id")
		otp := ctx.PostForm("otp")

		db := *config.GetDb()
		var userData models.Users
		if res := db.First(&userData, userid); res.Error != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"Error": "User details not valid",
			})
			return
		}
		result := helpers.ValidateOtp(otp, userid)

		if result {
			userData.Isverified = true
			if res := db.Save(userData); res.Error != nil {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"Error": "Internal server error",
				})
				return
			}

			ctx.JSON(200, gin.H{
				"Message": "User verified success",
				"Error":   nil,
				"User":    userData,
			})
		} else {
			ctx.JSON(200, gin.H{
				"Status": "Otp is not valid or correct",
			})
		}
	}
}
