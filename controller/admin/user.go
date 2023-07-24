package admincontroller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/muhammedarif/Ecomapi/config"
	"github.com/muhammedarif/Ecomapi/models"
)

func AdminBlockUserController() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userid := ctx.Query("id")
		fmt.Println(userid)
		if userid == "" {
			ctx.JSON(http.StatusBadRequest, models.Response{
				Status:  false,
				Message: "Block user attempt failed",
				Error:   "Invalid user id",
			})
			return
		}

		db := *config.GetDb()
		var userDetails models.Users
		if res := db.First(&userDetails, userid); res.Error != nil {
			ctx.JSON(http.StatusInternalServerError, models.Response{
				Status:  false,
				Message: "Block user attempt failed",
				Error:   res.Error.Error(),
			})
			return
		}

		userDetails.Status = false
		if res := db.Save(&userDetails); res.Error != nil {
			ctx.JSON(http.StatusInternalServerError, models.Response{
				Status:  false,
				Message: "Block user attempt failed",
				Error:   res.Error.Error(),
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, models.Response{
			Status:  true,
			Message: "Block user attempt success || user blocked",
			Error:   nil,
		})
	}
}

func AdminUnblockUserController() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userid := ctx.Query("id")
		if userid == "" {
			ctx.JSON(http.StatusBadRequest, models.Response{
				Status:  false,
				Message: "Unblock user attempt failed",
				Error:   "Invalid user id",
			})
			return
		}

		db := *config.GetDb()
		var userDetails models.Users
		if res := db.First(&userDetails, userid); res.Error != nil {
			ctx.JSON(http.StatusInternalServerError, models.Response{
				Status:  false,
				Message: "Unblock user attempt failed",
				Error:   res.Error.Error(),
			})
			return
		}

		userDetails.Status = true
		if res := db.Save(&userDetails); res.Error != nil {
			ctx.JSON(http.StatusInternalServerError, models.Response{
				Status:  false,
				Message: "Unblock user attempt failed",
				Error:   res.Error.Error(),
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, models.Response{
			Status:  true,
			Message: "Unblock user attempt success || user blocked",
			Error:   nil,
		})
	}
}

// Get all users

func AdminGetAllUsersController() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		db := *config.GetDb()
		var users []models.Users
		if res := db.Find(&users); res.Error != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":  res.Error.Error(),
				"result": nil,
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":  nil,
			"result": users,
		})
	}
}
