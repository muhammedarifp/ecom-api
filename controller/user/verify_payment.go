package usercontroller

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"

	"github.com/gin-gonic/gin"
	"github.com/muhammedarif/Ecomapi/config"
	"github.com/muhammedarif/Ecomapi/helpers"
	"github.com/muhammedarif/Ecomapi/models"
	"github.com/muhammedarif/Ecomapi/private"
)

type VerifyInps struct {
	OrderID   string `json:"order_id"`
	PaymentID string `json:"payment_id"`
	Signature string `json:"signature"`
}

func VerifyOrder() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		db := *config.GetDb()
		userID := helpers.GetUserIDFromJwt(ctx)
		var inpData VerifyInps
		if err := ctx.ShouldBindJSON(&inpData); err != nil {
			ctx.JSON(400, gin.H{
				"Error": "Order not success",
			})
			return
		}

		var transactionData models.Transactions
		if res := db.First(&transactionData, inpData.OrderID); res.Error != nil {
			ctx.JSON(400, gin.H{
				"Error": "Transaction not found",
			})
			return
		}

		var data = transactionData.OrderID + "|" + inpData.PaymentID

		//
		hash := hmac.New(sha256.New, []byte(private.RAZORPAY_SECRET))
		hash.Write([]byte(data))
		genarated_signature := hex.EncodeToString(hash.Sum(nil))

		if inpData.Signature != genarated_signature {
			ChangeOrderStatus(userID, inpData.OrderID)
			transactionData.Attempts++
			if res := db.Save(&transactionData); res.Error != nil {
				ctx.JSON(200, gin.H{
					"Error": "Transaction not verified | Internal server error",
				})
			}
			ctx.JSON(200, gin.H{
				"Error": "Transaction not verified",
			})
		} else {
			ctx.JSON(200, gin.H{
				"Message": "Transaction verified",
			})
		}
	}
}

func ChangeOrderStatus(userID, orderID string) bool {
	db := *config.GetDb()
	var orderData models.Orders
	if res := db.Where("id = ? AND user_id = ?", orderID, userID).First(&orderData); res.Error != nil {
		return false
	}

	orderData.IsSuccess = true
	orderData.Status = "Success"

	if res := db.Save(&orderData); res.Error != nil {
		return false
	}

	return true
}
