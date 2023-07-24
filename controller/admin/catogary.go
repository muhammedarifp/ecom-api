package admincontroller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/muhammedarif/Ecomapi/config"
	"github.com/muhammedarif/Ecomapi/models"
)

type CatogaryModel struct {
	Name string `form:"name"`
	Desc string `form:"desc"`
}

func AdminAddNewCatogary() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var CatogaryData CatogaryModel
		if err := ctx.ShouldBind(&CatogaryData); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
				Status:  false,
				Message: "Catogary adding failed",
				Error:   err.Error(),
			})
			return
		}
		db := *config.GetDb()

		//
		var Catogary = models.Catogary{
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
			Message: "C-- Please select an object in the tree view.atogary adding success",
			Error:   nil,
		})
	}
}

// Get all catogarys

func AdminGetAllCatogarys() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var catogarys []models.Catogary
		db := *config.GetDb()
		if res := db.Select("id", "name", "disc").Find(&catogarys); res.Error != nil {
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
