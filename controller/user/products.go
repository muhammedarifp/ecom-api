package usercontroller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/muhammedarif/Ecomapi/config"
	"github.com/muhammedarif/Ecomapi/models"
	"gorm.io/gorm"
)

// Get product by id
func UserGetProductByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		proid := ctx.Query("id")
		db := *config.GetDb()
		var proDeta models.Products
		var imgDeta []models.ProductImages

		if res := db.First(&proDeta, proid); res.Error != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"Status":  false,
				"Message": "Product searching failed",
				"Error":   res.Error.Error(),
			})
			return
		}

		if res := db.Find(&imgDeta, `product_id = ?`, proid); res.Error != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"Status":  false,
				"Message": "Product Image searching failed",
				"Error":   res.Error.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"Status":  true,
			"Message": "Search success",
			"Error":   nil,
			"Result":  proDeta,
			"Image":   imgDeta,
		})
	}
}

// Get all products

type Product struct {
	gorm.Model
	Name       string `gorm:"name"`
	Disc       string `gorm:"disc"`
	Price      uint   `gorm:"price"`
	Quntity    uint   `gorm:"quntity"`
	CatogaryID uint   `gorm:"catogary_id"`
	Images     []models.ProductImages
}

func GetallProducts() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// userID := helpers.GetUserIDFromJwt(ctx)
		page := ctx.Param("page")
		db := *config.GetDb()
		var products []Product

		pageint, err := strconv.Atoi(page)
		if err != nil {
			ctx.AbortWithStatusJSON(400, gin.H{
				"Error": "Invalid page number",
			})
			return
		}
		offset := (pageint - 1) * 10
		if res := db.Table("products").Offset(offset).Limit(10).Find(&products); res.Error != nil {
			ctx.AbortWithStatusJSON(400, gin.H{
				"Error": res.Error.Error(),
			})
			return
		}

		result := []Product{}
		for _, product := range products {
			var images []models.ProductImages
			if res := db.Find(&images, `product_id = ?`, product.ID); res.Error != nil {
				ctx.JSON(400, gin.H{"Error": "Internal server error"})
				ctx.Abort()
				return
			}

			product.Images = images
			result = append(result, product)
		}

		ctx.JSON(200, gin.H{
			"Page":     page,
			"Limit":    10,
			"Products": result,
		})
	}
}
