package admincontroller

import (
	"github.com/gin-gonic/gin"
	"github.com/muhammedarif/Ecomapi/config"
	"github.com/muhammedarif/Ecomapi/models"
)

func AddCatogaryOffer() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		db := *config.GetDb()
		type AdminInpStruct struct {
			CatogoryID   uint    `json:"catogary_id"`
			ExpiryDays   uint    `json:"expiry_days"`
			MinCheckount float64 `json:"min_checkout"`
			Discount     float64 `json:"discount"`
			MaxDiscount  float64 `json:"max_discount"`
		}

		response := map[string]any{
			"Status": true,
		}

		var adminInp AdminInpStruct
		if err := ctx.ShouldBindJSON(&adminInp); err != nil {
			response["Status"] = false
			ctx.JSON(400, response)
			return
		}

		newCatogaryOffer := models.CatogoryOffer{
			CatogoryID:           adminInp.CatogoryID,
			Expiry:               adminInp.ExpiryDays,
			MinTransactionAmount: adminInp.MinCheckount,
			Dicount:              adminInp.Discount,
			MaxDiscount:          adminInp.MaxDiscount,
		}

		if err := db.Create(&newCatogaryOffer).Error; err != nil {
			response["Status"] = false
			ctx.JSON(400, response)
			return
		}

		ctx.JSON(200, response)
	}
}

func DeleteCatogoryOffer() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		db := *config.GetDb()
		offerid := ctx.Param("offerid")
		response := map[string]any{
			"Status":  true,
			"Message": "Offer Deleted",
		}
		if err := db.Delete(&models.CatogoryOffer{}, offerid).Error; err != nil {
			response["Status"] = false
			response["Message"] = "Offer not deleted"
			ctx.JSON(200, response)
			return
		}

		response["Status"] = false
		response["Message"] = "Offer not deleted"
		ctx.JSON(200, response)
	}
}
