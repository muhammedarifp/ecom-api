package usercontroller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/muhammedarif/Ecomapi/config"
	"github.com/muhammedarif/Ecomapi/helpers"
	"github.com/muhammedarif/Ecomapi/models"
	"gorm.io/gorm"
)

// Get all cart products

type cartProductResponse struct {
	CartID  uint
	Product struct {
		ProductID uint
		Name      string
		Disc      string
		Price     float64
		Quntity   uint
		Image     string
	}
}

func GetAllCartProducts() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		db := *config.GetDb()
		var cartProducts []models.UserCart
		userID := helpers.GetUserIDFromJwt(ctx)
		var Images []models.ProductImages
		var result []cartProductResponse
		if res := db.Find(&cartProducts, `user_id = ?`, userID); res.Error != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
				Status:  false,
				Message: "Cart Products not found",
				Error:   res.Error.Error(),
			})
			return
		}

		for _, val := range cartProducts {
			var ProData models.Products
			db.First(&ProData, val.ProductID)
			db.Joins("JOIN products ON products.id = product_images.product_id").First(&Images)
			newRes := cartProductResponse{
				CartID: val.ID,
				Product: struct {
					ProductID uint
					Name      string
					Disc      string
					Price     float64
					Quntity   uint
					Image     string
				}{
					ProductID: ProData.ID,
					Name:      ProData.Name,
					Disc:      ProData.Disc,
					Price:     ProData.Price,
					Quntity:   uint(val.ProductCount),
					Image:     Images[0].ImageName,
				},
			}

			result = append(result, newRes)
		}

		ctx.JSON(http.StatusOK, gin.H{
			"ProductsCount": len(result),
			"UserID":        userID,
			"Products":      result,
		})
	}
}

// Add To Cart Controller
func UserAddToCartController() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userid := helpers.GetUserIDFromJwt(ctx)
		proid := ctx.Param("productid")
		db := *config.GetDb()

		var productDeta models.Products
		if res := db.First(&productDeta, proid); res.Error != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
				Status:  false,
				Message: "Product not added on cart",
				Error:   res.Error.Error(),
			})
			return
		}

		var cartProduct models.UserCart
		if res := db.First(&cartProduct, `product_id = ? AND user_id = ?`, proid, userid); res.Error != nil {
			fmt.Println(res.Error.Error())
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				newCartPro := models.UserCart{
					UserID:       userid,
					ProductID:    productDeta.ID,
					ProductCount: 1,
				}

				if res := db.Create(&newCartPro); res.Error != nil {
					ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
						Status:  false,
						Message: "Product not added on cart",
						Error:   res.Error.Error(),
					})
					return
				}

				ctx.JSON(http.StatusOK, models.Response{
					Status:  false,
					Message: "Product added on cart",
					Error:   nil,
				})
			} else {
				ctx.AbortWithStatusJSON(400, gin.H{
					"Error": "Internal server error",
				})
			}

			return
		}

		cartProduct.ProductCount = cartProduct.ProductCount + 1
		if res := db.Save(&cartProduct); res.Error != nil {
			ctx.AbortWithStatusJSON(400, gin.H{"Error": "Add to cart attempt failed"})
			return
		}

		ctx.JSON(200, gin.H{
			"Message": "Product Count updated",
			"Error":   nil,
		})
	}
}

// Quntity Inc Or Dec in Cart controller
func UserCartQuntityController() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		db := *config.GetDb()
		cartid := ctx.Param("cartid")
		change := ctx.Param("change")

		fmt.Println(cartid + " " + change)

		var cartData models.UserCart
		db.First(&cartData, `id = ?`, cartid)
		if change == "inc" {
			cartData.ProductCount = cartData.ProductCount + 1
			if res := db.Save(&cartData); res.Error != nil {
				ctx.JSON(http.StatusOK, gin.H{
					"Status":  false,
					"Message": "Cart count not changed",
					"Error":   res.Error.Error(),
				})
				return
			}

			ctx.JSON(http.StatusOK, gin.H{
				"Status":   true,
				"Message":  "Cart count changed",
				"Error":    nil,
				"NewCount": cartData.ProductCount,
			})

		} else if change == "dec" {
			if cartData.ProductCount <= 1 {
				if res := db.Delete(cartData, cartid); res.Error != nil {
					ctx.JSON(http.StatusOK, gin.H{
						"Status":  false,
						"Message": "Cart count not changed",
						"Error":   res.Error.Error(),
					})
					return
				}

				ctx.JSON(http.StatusOK, gin.H{
					"Status":   true,
					"Message":  "Product Deleted",
					"Error":    nil,
					"NewCount": 0,
				})

			} else {
				cartData.ProductCount = cartData.ProductCount - 1
				if res := db.Save(&cartData); res.Error != nil {
					ctx.JSON(http.StatusOK, gin.H{
						"Status":  false,
						"Message": "Cart count not changed",
						"Error":   res.Error.Error(),
					})

					return
				}

				ctx.JSON(http.StatusOK, gin.H{
					"Status":   true,
					"Message":  "Cart count changed",
					"Error":    nil,
					"NewCount": cartData.ProductCount,
				})
			}
		} else {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
				Status:  false,
				Message: "Cart count not changed",
				Error:   "Invalid input",
			})
		}
	}
}

// Remove Cart Controller
func UserRemoveCartController() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Create database instance
		cartid := ctx.Query("cartid")
		db := *config.GetDb()
		if res := db.Table("user_carts").Delete(&models.UserCart{}, `id = ?`, cartid); res.Error != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
				Status:  false,
				Message: "Cart Item not deleted",
				Error:   res.Error.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, models.Response{
			Status:  true,
			Message: "Cart item removed",
			Error:   nil,
		})
	}
}
