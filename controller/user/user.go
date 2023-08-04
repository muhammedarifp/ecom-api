package usercontroller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/muhammedarif/Ecomapi/config"
	"github.com/muhammedarif/Ecomapi/helpers"
	"github.com/muhammedarif/Ecomapi/models"
	"gorm.io/gorm"
)

type userLoginForm struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func UserLoginController() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Create a variable using user model
		var userLoginData userLoginForm

		// Bind user enterd val use struct
		if err := ctx.ShouldBindJSON(&userLoginData); err != nil {
			ctx.JSON(http.StatusBadRequest, models.Response{
				Status:  false,
				Message: "Login failed",
				Error:   err.Error(),
			})
			return
		}

		// Validate user enter details
		validate := validator.New()
		if err := validate.Struct(userLoginData); err != nil {
			ctx.JSON(http.StatusBadRequest, models.Response{
				Status:  false,
				Message: "Login failed",
				Error:   err.Error(),
			})
			return
		}

		// Create database instance and serach user data using user email
		db := *config.GetDb()
		var userData models.Users
		if res := db.First(&userData, `email = ?`, userLoginData.Email); res.Error != nil {
			ctx.JSON(http.StatusBadRequest, models.Response{
				Status:  false,
				Message: "Login failed",
				Error:   "User not found",
			})
			return
		}

		// This case work is success. status ok [200]
		if userData.Email == userLoginData.Email && helpers.CompareHashPass(userLoginData.Password, userData.Password) {
			if userData.Status {

				token := helpers.CreateUserJwtToken(userData.ID)

				ctx.JSON(http.StatusOK, gin.H{
					"Status":  true,
					"Message": "Login Success",
					"Error":   nil,
					"Token":   token,
				})
				return
			} else {
				ctx.AbortWithStatusJSON(403, models.Response{
					Status:  false,
					Message: "Login failed",
					Error:   "This user is banned || blocked",
				})
				return
			}
		}

		// This case work email or password incorrect
		ctx.JSON(http.StatusNetworkAuthenticationRequired, models.Response{
			Status:  false,
			Message: "Login failed",
			Error:   "Email or password incorrect",
		})
	}
}

// End the User login controller
// ----------------------------------------------------------------------------------------------
// User Signup helper controller

type UserSignupFormData struct {
	FirstName    string `json:"first_name" validate:"required"`
	LastName     string `json:"last_name" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
	Mobile       string `json:"mobile" validate:"required"`
	Password     string `json:"password" validate:"required,min=6"`
	RefferelCode string `json:"refferel_code"`
}

func UserSignupController() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Bind form data using struct
		var userSignupData UserSignupFormData
		if err := ctx.ShouldBindJSON(&userSignupData); err != nil {
			ctx.JSON(http.StatusBadRequest, models.Response{
				Status:  false,
				Message: "Signup failed",
				Error:   err.Error(),
			})
			return
		}

		// Validate user entered data
		validate := validator.New()
		if err := validate.Struct(userSignupData); err != nil {
			ctx.JSON(http.StatusBadRequest, models.Response{
				Status:  false,
				Message: "Signup failed",
				Error:   err.Error(),
			})
			return
		}

		// Database init
		db := *config.GetDb()

		// Check email exist or not
		var preUser models.Users
		res := db.First(&preUser, `email = ?`, userSignupData.Email)

		// User entered email is aldredy in the database
		if res.Error == nil {
			ctx.JSON(http.StatusBadRequest, models.Response{
				Status:  false,
				Message: "Signup failed",
				Error:   "Email is aldredy exist",
			})

			// Create new user
		} else if errors.Is(res.Error, gorm.ErrRecordNotFound) {

			// Create new user struct
			pass := helpers.PassToHash(userSignupData.Password)
			reffrelCode := helpers.CreateReffrerelCode()
			newUser := models.Users{
				FirstName:    userSignupData.FirstName,
				LastName:     userSignupData.LastName,
				Email:        userSignupData.Email,
				Mobile:       userSignupData.Mobile,
				Password:     pass,
				Isadmin:      false,
				Isverified:   false,
				Status:       true,
				ReferralCode: reffrelCode,
			}

			// Upload this struct into db
			if response := db.Create(&newUser); response.Error != nil {
				ctx.JSON(http.StatusInternalServerError, models.Response{
					Status:  false,
					Message: "Signup failed",
					Error:   "Server issue",
				})
				return
			}

			response := gin.H{
				"status":  true,
				"message": "Signup Success",
				"error":   "Your signup is success !",
				"user": map[string]any{
					"user_id":       newUser.ID,
					"first_name":    newUser.FirstName,
					"last_name":     newUser.LastName,
					"is_verified":   newUser.Isverified,
					"refferel_code": newUser.ReferralCode,
				},
			}

			if userSignupData.RefferelCode != "" {
				var referedUserData models.Users
				if err := db.First(&referedUserData, `referral_code = ?`, userSignupData.RefferelCode); err == nil {
					var referdUserWallet models.Wallets
					var newUserWallet models.Wallets
					db.First(&referdUserWallet, `user_id = ?`, referedUserData.ID)
					db.First(&newUserWallet, `user_id = ?`, newUser.ID)
					referdUserWallet.UserID = referedUserData.ID
					referdUserWallet.Balance = referdUserWallet.Balance + 100
					newUserWallet.Balance = newUserWallet.Balance + 30
					newUserWallet.UserID = newUser.ID
					db.Save(&referdUserWallet)
					db.Save(&newUserWallet)
					response["Bonus"] = "Refferel code applied"
					ctx.JSON(200, response)
					return
				}
			}

			response["Bonus"] = "Refferel code empty or invalid"
			ctx.JSON(200, response)

			// More qury errors
		} else {
			ctx.JSON(http.StatusInternalServerError, models.Response{
				Status:  false,
				Message: "Signup failed",
				Error:   res.Error.Error(),
			})
		}
	}
}
