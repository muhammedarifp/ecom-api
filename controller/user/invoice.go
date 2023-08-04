package usercontroller

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/muhammedarif/Ecomapi/config"
	"github.com/muhammedarif/Ecomapi/helpers"
	"github.com/muhammedarif/Ecomapi/models"
)

func DownloadInvoice() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		orderid := ctx.Param("orderid")
		db := *config.GetDb()
		userID := helpers.GetUserIDFromJwt(ctx)
		var orders models.Order
		var address models.Address
		var userData models.Users
		db.First(&userData, userID)
		db.Preload("Items").Where("id = ? AND user_id = ?", orderid, userID).Find(&orders)

		for j := range orders.Items {
			db.Preload("Product").First(&orders.Items[j].Product, orders.Items[j].ProductID)
		}

		// Fetch address
		db.First(&address, "user_id = ? AND is_default = true", userID)

		pdf := *helpers.CreateInvoice(orders, address, userData)

		ctx.Header("Content-Disposition", "attachment; filename=OR"+fmt.Sprint(orders.ID)+".pdf")
		ctx.Header("Content-Type", "application/octet-stream")

		pdf.Output(ctx.Writer)
		ctx.JSON(200, orders)
	}
}
