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

// This Func uses get user addresses
func UserGetAddressesController() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userid := helpers.GetUserIDFromJwt(ctx)
		if userid == "" {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
				Status:  false,
				Message: "Address fetch failed",
				Error:   "Authentication error",
			})
			return
		}

		helpers.GetUserIDFromJwt(ctx)
		db := *config.GetDb()
		var userAddresses []models.Address
		if res := db.Find(&userAddresses, `user_id = ?`, userid); res.Error != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
				Status:  false,
				Message: "Address fetch failed",
				Error:   res.Error.Error(),
			})
			return
		}

		if len(userAddresses) <= 0 {
			ctx.AbortWithStatusJSON(201, models.Response{
				Status:  true,
				Message: "Address is empty",
				Error:   "Address is empty",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"Status":    true,
			"Message":   "Address fetch successfully",
			"Error":     nil,
			"Count":     len(userAddresses),
			"Addresses": userAddresses,
		})
	}
}

// This function usage is create new address. Iam applied logic is create address used address predifigned model
// if user address is empty i created a default address, otherwise I create non default address
type UserAddressInp struct {

	// User name. This name use on delivery time
	Name string `json:"name" validate:"required"`

	// Mobile number
	Mobile string `json:"mobile" validate:"required,min=10,max=10"`

	// User pincode
	Pincode string `json:"pincode" validate:"required"`

	// Pincode Location
	Locality string `json:"locality" validate:"required"`

	// Main address box, this field area and street details
	Address string `json:"address" validate:"required,min=15"`

	// Users city/Distinct/Town name. but this field called city. for simplicity purpous
	City string `json:"city" validate:"required"`

	// State like kerala or more
	State string `json:"state" validate:"required"`

	// Landmark. This is a optional field
	Landmark string `json:"landmark"`

	// AlternatePhone number. optional
	AlternatePhone string `json:"alternate_phone"`
}

func UserCreateAddressController() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var userInpData UserAddressInp
		db := *config.GetDb()
		userID := helpers.GetUserIDFromJwt(ctx)

		// Bind user enter data into json
		if err := ctx.ShouldBindJSON(&userInpData); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
				Status:  false,
				Message: "Address Not Added",
				Error:   err.Error(),
			})
			return
		}

		// Validate user enter details..!
		validator := validator.New()
		if err := validator.Struct(&userInpData); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
				Status:  false,
				Message: "Address not added",
				Error:   err.Error(),
			})
			return
		}

		// Check user address aldredy alive or not
		var userAddressModel models.Address
		if res := db.First(&userAddressModel); res.Error == nil {
			newAddress := models.Address{
				UserID:         userID,
				Name:           userInpData.Name,
				Mobile:         userInpData.Mobile,
				Pincode:        userInpData.Pincode,
				Locality:       userInpData.Locality,
				Address:        userInpData.Address,
				City:           userInpData.City,
				State:          userInpData.State,
				Landmark:       userInpData.Landmark,
				AlternatePhone: userInpData.AlternatePhone,
				IsDefault:      false,
			}

			if result := db.Create(&newAddress); result.Error != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
					Status:  false,
					Message: "Address not added",
					Error:   result.Error.Error(),
				})
				return
			}

			ctx.AbortWithStatusJSON(http.StatusOK, models.Response{
				Status:  true,
				Message: "Address added successfully",
				Error:   nil,
			})

		} else if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			newAddress := models.Address{
				UserID:         userID,
				Name:           userInpData.Name,
				Mobile:         userInpData.Mobile,
				Pincode:        userInpData.Pincode,
				Locality:       userInpData.Locality,
				Address:        userInpData.Address,
				City:           userInpData.City,
				State:          userInpData.State,
				Landmark:       userInpData.Landmark,
				AlternatePhone: userInpData.AlternatePhone,
				IsDefault:      true,
			}

			if result := db.Create(&newAddress); result.Error != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
					Status:  false,
					Message: "Address not added",
					Error:   result.Error.Error(),
				})
				return
			}

			ctx.AbortWithStatusJSON(http.StatusOK, models.Response{
				Status:  false,
				Message: "Address added successfully",
				Error:   nil,
			})

		} else {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
				Status:  false,
				Message: "Address not added",
				Error:   "Internal server error",
			})
		}
	}
}

// Remove user address

func RemoveUserAddress() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := helpers.GetUserIDFromJwt(ctx)
		addressID := ctx.Query("add-id")
		db := *config.GetDb()
		var addressData models.Address
		var secondaryAddress models.Address

		if res := db.First(&addressData, `id = ? AND user_id = ?`, addressID, userID); res.Error != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
				Status:  false,
				Message: "Address not deleted",
				Error:   "Invalid status id",
			})

			return
		}

		if res := db.Delete(&addressData, addressID); res.Error != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
				Status:  false,
				Message: "Address not deleted",
				Error:   "Invalid status id",
			})
			return
		}

		if addressData.IsDefault {
			if res := db.First(&secondaryAddress, `user_id = ?`, userID); res.Error == nil {
				secondaryAddress.IsDefault = true
				db.Save(&secondaryAddress)
			}
		}

		ctx.JSON(http.StatusOK, gin.H{
			"Status":  true,
			"Message": "Address deleted",
			"Error":   nil,
		})
	}
}
