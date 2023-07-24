package usercontroller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/muhammedarif/Ecomapi/config"
	"github.com/muhammedarif/Ecomapi/helpers"
	"github.com/muhammedarif/Ecomapi/models"
)

type ResponseModel struct {
	OrderID   uint
	Ordered   time.Time
	Status    string
	PayMethod string
	Product   struct {
		Name  string
		Price int
		Image string
	}
}

func UserGetAllOrders() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userid := helpers.GetUserIDFromJwt(ctx)
		db := *config.GetDb()
		var userData models.Users

		// Get user data using user id
		if res := db.First(&userData, userid); res.Error != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
				Status:  false,
				Message: "Orders not fetched",
				Error:   "Invalid user details - " + res.Error.Error(),
			})

			return
		}

		var userOrders []models.Orders
		if res := db.Find(&userOrders, `user_id = ?`, userid); res.Error != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
				Status:  false,
				Message: "Orders not fetched",
				Error:   res.Error.Error(),
			})

			return
		}

		// var result []ResponseModel
		var orderProductDeta []models.OrdersItems

		// Loop
		for _, val := range userOrders {
			if res := db.Find(&orderProductDeta, `order_id = ?`, val.ID); res.Error != nil {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"Error": res.Error.Error(),
				})
				return
			}

			// newItem := ResponseModel{
			// 	OrderID:   userData.ID,
			// 	Ordered:   val.CreatedAt,
			// 	Status:    val.Status,
			// 	PayMethod: val.PayMethod,
			// 	Product: struct {
			// 		Name  string
			// 		Price int
			// 		Image string
			// 	}{
			// 		Name: orderProductDeta[0].N,
			// 	},
			// }
		}

		ctx.JSON(http.StatusOK, gin.H{
			"Status":  true,
			"Message": "Orders fetched success",
			"Orders":  userOrders,
		})
	}
}

type CancelInps struct {
	OrderID uint   `json:"orderid"`
	Message string `json:"message"`
}

func CancelOrder() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var cancelInps CancelInps
		ctx.ShouldBindJSON(&cancelInps)
		db := *config.GetDb()
		userID := helpers.GetUserIDFromJwt(ctx)
		var orderData models.Orders
		// get order details on database
		if res := db.First(&orderData, cancelInps.OrderID); res.Error != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"Error": res.Error.Error(),
			})
			return
		}

		if !InitRefundToWallet(orderData.ID, userID) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"Error": "Refund not initialized",
			})
			return
		}

		// After fetch order details
		orderData.Status = "Cancelled"
		if res := db.Save(&orderData); res.Error != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"Error": res.Error.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"Message": "Order Cancelled",
			"Error":   nil,
			"Order":   orderData,
		})
	}
}

func InitRefundToWallet(orderID uint, userID string) bool {
	db := *config.GetDb()
	var OrderData models.Orders
	var walletData models.Wallets
	tx := db.Begin()
	if res := tx.First(&OrderData, orderID); res.Error != nil {
		return false
	}
	if OrderData.PayMethod == "COD" {
		return false
	}

	if res := tx.First(&walletData, `user_id = ?`, userID); res.Error != nil {
		return false
	}

	walletData.Balance = OrderData.TottalAmount
	if res := tx.Save(&walletData); res.Error != nil {
		return false
	}

	if res := tx.Commit(); res.Error != nil {
		return false
	}

	return true
}

// Return

// func ReturnOrder() gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		// Create database instance
// 		orderid := ctx.Query("orderid")
// 		db := *config.GetDb()
// 		var orderData models.Orders
// 		db.First(&orderData, orderid)
// 		orderDeff := time.Now().Sub(orderData.CreatedAt).Hours()
// 		if orderDeff > 168 {
// 			ctx.AbortWithStatusJSON(http.StatusNotAcceptable, models.Response{
// 				Status:  false,
// 				Message: "Return order not allowed",
// 				Error:   nil,
// 			})
// 		} else {
// 			orderData.Status = "Returned"
// 			if res := db.Save(orderData); res.Error != nil {
// 				ctx.AbortWithStatusJSON(http.StatusNotAcceptable, models.Response{
// 					Status:  false,
// 					Message: "Return order not allowed",
// 					Error:   nil,
// 				})
// 				return
// 			}

// 			ctx.JSON(200, gin.H{
// 				"Message": "Your Order Returned Success",
// 				"Error":   nil,
// 				"Order":   orderData,
// 			})
// 		}
// 	}
// }
