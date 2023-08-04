package admincontroller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/muhammedarif/Ecomapi/config"
	"github.com/muhammedarif/Ecomapi/models"
)

type CatogaryModel struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

func AdminAddNewCatogary() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var CatogaryData CatogaryModel
		if err := ctx.ShouldBindJSON(&CatogaryData); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
				Status:  false,
				Message: "Catogary adding failed",
				Error:   err.Error(),
			})
			return
		}
		db := *config.GetDb()

		//
		var Catogary = models.Catogory{
			Name: CatogaryData.Name,
			Disc: CatogaryData.Desc,
		}
		if res := db.Create(&Catogary); res.Error != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
				Status:  false,
				Message: "Catogary adding failed",
				Error:   res.Error.Error(),
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusOK, models.Response{
			Status:  true,
			Message: "Adding Success",
			Error:   nil,
		})
	}
}

// Get all catogarys

func AdminGetAllCatogarys() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var catogarys []models.Catogory
		db := *config.GetDb()
		if res := db.Find(&catogarys); res.Error != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": res.Error.Error(),
			})
		}
		ctx.JSON(200, catogarys)
	}
}

// Filter Using Catogarys
// func CatogaryFiltering(ctx *gin.Context) {
// 	db := *config.GetDb()
// 	filterData := ctx.Query("filter")
// 	var filters
// 	db.Table("catogaries").First()
// }
