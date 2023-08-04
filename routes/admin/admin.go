package adminroutes

import (
	"github.com/gin-gonic/gin"
	admincontroller "github.com/muhammedarif/Ecomapi/controller/admin"
)

func AdminRoutes(auth *gin.RouterGroup, unauth *gin.RouterGroup) {
	// Admin Post requests

	//
	unauth.POST("admin/login", admincontroller.AdminLoginController())

	// User Controll
	auth.GET("admin/getall-users", admincontroller.AdminGetAllUsersController())
	auth.POST("admin/block", admincontroller.AdminBlockUserController())
	auth.POST("/admin/unblock", admincontroller.AdminUnblockUserController())

	// Product
	auth.POST("admin/product/add", admincontroller.AdminAddNewProductController())
	auth.POST("admin/product/edit", admincontroller.AdminEditProductController())
	auth.GET("admin/product/getbyid/:id", admincontroller.AdminGetProductUsingIDController())
	auth.GET("admin/product/getall/:page", admincontroller.AdminGetAllProductController())
	auth.DELETE("admin/product/delete/:id", admincontroller.AdminDeleteProductController())
	auth.Group("admin/product")

	// Catogary
	auth.POST("admin/add-catogary", admincontroller.AdminAddNewCatogary())
	auth.GET("admin/getall-catogarys", admincontroller.AdminGetAllCatogarys())

	// Admin order management
	auth.GET("admin/orders/getall/:page", admincontroller.GetallOrders())

	// Coupon management
	auth.POST("admin/coupon/add", admincontroller.CreateNewCoupon())
	auth.DELETE("admin/coupon/delete", admincontroller.DeleteCoupon())

	// Graph view controller
	auth.GET("admin/sales-graph", admincontroller.SalesGraph())

	// Offer management
	auth.POST("admin/add-cat-offer", admincontroller.AddCatogaryOffer())
	auth.DELETE("admin/delete-offer/:offerid", admincontroller.DeleteCatogoryOffer())

	// Report
	auth.POST("admin/download-sales-repo/:duration", admincontroller.DownloadSalesReport())
}
