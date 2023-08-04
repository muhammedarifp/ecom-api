package usercontroller

import (
	"github.com/gin-gonic/gin"
	"github.com/muhammedarif/Ecomapi/config"
	"github.com/muhammedarif/Ecomapi/models"
)

type ProductImages struct {
	ImageName string `json:"image_name"`
}

func SearchProductUsingID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		productid := ctx.Param("productid")
		db := *config.GetDb()
		var proData models.Product
		var proImages []ProductImages
		if res := db.First(&proData, productid); res.Error != nil {
			ctx.AbortWithStatusJSON(400, gin.H{
				"Error":   "Product not fetched",
				"Message": res.Error.Error(),
			})
			return
		}

		if res := db.Find(&proImages, `product_id = ?`, productid); res.Error != nil {
			ctx.AbortWithStatusJSON(400, gin.H{
				"Error":   "Product not fetched",
				"Message": res.Error.Error(),
			})
			return
		}

		ctx.JSON(200, gin.H{
			"Product": map[string]any{
				"ProductID": proData.ID,
				"Create":    proData.CreatedAt,
				"Name":      proData.Name,
				"Disc":      proData.Disc,
				"Price":     proData.Price,
				"Quntity":   proData.Quntity,
				"Images":    proImages,
			},
		})
	}
}

func SearchProductsUsingCatogary() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		name := ctx.Param("catogary")
		var products []models.Product
		var catogary models.Catogory
		db := *config.GetDb()
		if res := db.Find(&catogary, `name = ?`, name); res.Error != nil {
			ctx.AbortWithStatusJSON(400, gin.H{
				"Products": products,
			})
			return
		}

		if res := db.Find(&products, `catogary_id = ?`, catogary.ID); res.Error != nil {
			ctx.AbortWithStatusJSON(400, gin.H{
				"Products": products,
			})
			return
		}

		ctx.JSON(200, gin.H{
			"Products": products,
		})
	}
}
