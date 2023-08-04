package admincontroller

import (
	"github.com/gin-gonic/gin"
	"github.com/muhammedarif/Ecomapi/config"
	"github.com/muhammedarif/Ecomapi/models"
)

type CreateCouponInp struct {
	Code        string  `json:"code"`
	Discount    float64 `json:"discount"` // in persentage
	ExpiryDate  int     `json:"expiry"`   // In days
	MinAmount   float64 `json:"min"`      // min purchase amount
	MaxDiscount float64 `json:"maxdisc"`  // max discount amount
}

func CreateNewCoupon() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var InpData CreateCouponInp
		db := *config.GetDb()
		if err := ctx.ShouldBindJSON(&InpData); err != nil {
			ctx.AbortWithStatusJSON(400, gin.H{
				"Error": err.Error(),
			})
			return
		}

		// Create new token
		newToken := models.Coupon{
			Code:              InpData.Code,
			Discount:          InpData.Discount,
			ExpiryDate:        InpData.ExpiryDate,
			MinAmount:         InpData.MinAmount,
			MaxDiscountAmount: InpData.MaxDiscount,
		}
		if res := db.Create(&newToken); res.Error != nil {
			ctx.AbortWithStatusJSON(400, gin.H{
				"Error": res.Error.Error(),
			})
			return
		}

		ctx.JSON(200, gin.H{
			"Status":  true,
			"Message": "Coupon added success",
			"coupon": map[string]any{
				"Create":            newToken.CreatedAt,
				"CouponCode":        newToken.Code,
				"Discount":          newToken.Discount,
				"Expiry":            newToken.ExpiryDate,
				"MinPurchaseAmount": newToken.MinAmount,
				"MaxDiscountAmount": newToken.MaxDiscountAmount,
			},
		})
	}
}

// Delete Coupon
func DeleteCoupon() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		couponid := ctx.Query("couponid")
		db := *config.GetDb()
		if res := db.Delete(&models.Coupon{}, couponid); res.Error != nil {
			ctx.AbortWithStatusJSON(400, gin.H{
				"Error": res.Error.Error(),
			})
			return
		}

		ctx.JSON(200, gin.H{
			"Message": "Coupon Deleted",
		})
	}
}
