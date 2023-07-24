package usercontroller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/muhammedarif/Ecomapi/config"
	"github.com/muhammedarif/Ecomapi/helpers"
	"github.com/muhammedarif/Ecomapi/models"
	"gorm.io/gorm"
)

// This function uses get user profile. This Function return user details.
func UserGetProfileController() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userid := ctx.Query("id")
		var usrdata models.Users
		db := *config.GetDb()

		if res := db.First(&usrdata, userid); res.Error != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
				Status:  false,
				Message: "User not found",
				Error:   res.Error.Error(),
			})

			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"created":       usrdata.CreatedAt,
			"updated":       usrdata.UpdatedAt,
			"first_name":    usrdata.FirstName,
			"last_name":     usrdata.LastName,
			"email":         usrdata.Email,
			"user_verified": usrdata.Isverified,
			"mobile":        usrdata.Mobile,
		})
	}
}

// Edit User Profile

type UpdatedDetails struct {
	UserID    uint     `form:"id" validate:"required"`
	FirstName string   `form:"first_name" validate:"required"`
	LastName  string   `form:"last_name" validate:"required"`
	Email     string   `form:"email" validate:"required,email"`
	Mobile    string   `form:"mobile" validate:"required,min=10,max=10"`
	Images    []string `form:""`
}

// This is user edited details
func UserEditProfileController() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var newData UpdatedDetails
		ctx.ShouldBind(&newData)
		validate := validator.New()
		if err := validate.Struct(&newData); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
				Status:  false,
				Message: "UserData not updated",
				Error:   err.Error(),
			})
			return
		}

		// Instance of database
		db := *config.GetDb()
		var UserData models.Users

		// database user serching context errors
		if res := db.Where(newData.UserID).First(&UserData); res.Error != nil {
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
					Status:  false,
					Message: "UserData not updated",
					Error:   "Entred user id not valid !",
				})
			} else {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
					Status:  false,
					Message: "UserData not updated",
					Error:   res.Error.Error(),
				})
			}

			return
		}

		UserData.FirstName = newData.FirstName
		UserData.LastName = newData.LastName
		UserData.Email = newData.Email
		UserData.Mobile = newData.Mobile

		if res := db.Save(&UserData); res.Error != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
				Status:  false,
				Message: "UserData not updated",
				Error:   res.Error.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status":   true,
			"message":  "UserData updated successfully",
			"Error":    nil,
			"UserData": UserData,
		})
	}
}

// Forgot password

func UserForgotPasswordAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		db := *config.GetDb()
		email := ctx.PostForm("email")
		var userData models.Users

		if res := db.First(&userData, `email = ?`, email); res.Error != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
				Status:  false,
				Message: "Forgot password attemp failed",
				Error:   res.Error.Error(),
			})
			return
		}
		otp := helpers.GenarateOtp()
		useridToStr := fmt.Sprint(userData.ID)
		if ok, err := helpers.SendOtpEmail(useridToStr, otp); !ok {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
				Status:  false,
				Message: "Forgot password attemp failed",
				Error:   err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, models.Response{
			Status:  true,
			Message: "Otp sended on your mail",
			Error:   nil,
		})
	}
}

func OtpValidateAndResetPassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		db := *config.GetDb()
		var userData models.Users
		userid := ctx.PostForm("userid")
		otp := ctx.PostForm("otp")
		newpass := ctx.PostForm("newpass")
		confirmpass := ctx.PostForm("confirmpass")

		if res := helpers.ValidateOtp(otp, userid); !res {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
				Status:  false,
				Message: "Forgot password attemp failed",
				Error:   "Entered Incorrect or not valid otp",
			})
			return
		}

		if newpass == confirmpass && len(confirmpass) >= 6 {
			hash := helpers.PassToHash(newpass)
			if res := db.First(&userData, userid); res.Error != nil {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
					Status:  false,
					Message: "Forgot password attemp failed",
					Error:   res.Error.Error(),
				})
				return
			}

			userData.Password = hash
			if res := db.Save(userData); res != nil {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
					Status:  false,
					Message: "Forgot password attemp failed",
					Error:   res.Error.Error(),
				})
				return
			}

			ctx.JSON(http.StatusOK, gin.H{
				"Status":  true,
				"Message": "Your password changed success",
				"Error":   nil,
				"User":    userData,
			})
		}
	}
}
