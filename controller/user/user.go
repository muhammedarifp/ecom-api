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
	FirstName string `form:"first_name" validate:"required"`
	LastName  string `form:"last_name" validate:"required"`
	Email     string `form:"email" validate:"required,email"`
	Mobile    string `form:"mobile" validate:"required"`
	Password  string `form:"password" validate:"required,min=6"`
}

func UserSignupController() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// Bind form data using struct
		var userSignupData UserSignupFormData
		if err := ctx.ShouldBind(&userSignupData); err != nil {
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
			newUser := models.Users{
				FirstName:  userSignupData.FirstName,
				LastName:   userSignupData.LastName,
				Email:      userSignupData.Email,
				Mobile:     userSignupData.Mobile,
				Password:   pass,
				Isadmin:    false,
				Isverified: false,
				Status:     true,
			}

			// Upload this struct into db
			if response := db.Create(&newUser); response.Error != nil {

				// Incase get error on create user just send bad response
				ctx.JSON(http.StatusInternalServerError, models.Response{
					Status:  false,
					Message: "Signup failed",
					Error:   "Server issue",
				})
				return
			}

			// Seccess response
			ctx.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Signup Success",
				"error":   "Your signup is success !",
				"user": map[string]any{
					"user_id":     newUser.ID,
					"first_name":  newUser.FirstName,
					"last_name":   newUser.LastName,
					"is_verified": newUser.Isverified,
				},
			})

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
