package admincontroller

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/muhammedarif/Ecomapi/config"
	"github.com/muhammedarif/Ecomapi/models"
	"gorm.io/gorm"
)

type orders struct {
	gorm.Model
	Create    time.Time `gorm:"created_at"`
	UserID    string    `gorm:"user_id"`
	PayMethod string    `gorm:"pay_method"`
	IsSuccess string    `gorm:"is_success"`
}

type Products struct {
	Name string `gorm:"name"`
}

type OrdersResponse struct {
	Create        time.Time
	OrderID       uint
	UserID        string `gorm:"user_id"`
	PayMethod     string
	IsSuccess     string
	OrderProducts []models.OrderItem
	// Transaction
}

func GetallOrders() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		db := *config.GetDb()
		page := ctx.Param("page")
		pageint, err := strconv.Atoi(page)
		if err != nil {
			ctx.AbortWithStatusJSON(400, gin.H{
				"Error": "invalid input or invalid page number",
			})
			return
		}
		var orders []orders
		var result []OrdersResponse
		offset := (pageint - 1) * 10
		if res := db.Offset(offset).Limit(10).Find(&orders); res.Error != nil {
			ctx.AbortWithStatusJSON(400, gin.H{
				"Error": "order not fetched, internal server error",
			})
			return
		}

		for _, val := range orders {
			var ordersItems []models.OrderItem
			db.Where("order_id = ?", val.ID).Find(&ordersItems)
			newResponse := OrdersResponse{
				Create:        val.Create,
				OrderID:       val.ID,
				PayMethod:     val.PayMethod,
				IsSuccess:     val.IsSuccess,
				UserID:        val.UserID,
				OrderProducts: ordersItems,
			}

			result = append(result, newResponse)

		}

		ctx.JSON(200, gin.H{
			"Page":   pageint,
			"Limit":  10,
			"Length": len(orders),
			"Orders": result,
		})
	}
}

// func CancelOrder() gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		db :=
// 	}
// }
