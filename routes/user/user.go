package userroutes

import (
	"github.com/gin-gonic/gin"
	usercontroller "github.com/muhammedarif/Ecomapi/controller/user"
)

func UserRoutes(auth *gin.RouterGroup, unauth *gin.RouterGroup) {

	// Post requsts

	// This is login enpoint.
	unauth.POST("/login", usercontroller.UserLoginController())

	// Signup ..................................................................
	// http://localhost:8080/signup --- post
	/*
			type UserSignupFormData struct {
			FirstName string `form:"first_name" validate:"required"`
			LastName  string `form:"last_name" validate:"required"`
			Email     string `form:"email" validate:"required,email"`
			Mobile    string `form:"mobile" validate:"required"`
			Password  string `form:"password" validate:"required,min=6"`
		}
	*/
	unauth.POST("/signup", usercontroller.UserSignupController())

	// Requst Otp ..............................................................
	/*
		{
			id : your id
		}
	*/
	unauth.POST("/req-otp", usercontroller.UserVerifyOtpController())

	// Validate Otp ............................................................
	/*
		{
			email : your email,
			otp : user's otp
		}
	*/
	unauth.POST("/verify-otp", usercontroller.UserVarifyValidateController())

	// Products section
	unauth.GET("/getall-products/:page", usercontroller.GetallProducts())
	unauth.GET("/getpro-id", usercontroller.UserGetProductByID())

	// Profile Section
	auth.GET("/account", usercontroller.UserGetProfileController())
	auth.POST("/account/edit-details", usercontroller.UserEditProfileController())
	// auth.POST("/account/forgot-pass", usercontroller.UserForgotPassword())

	// Address manage
	auth.GET("/account/addresses", usercontroller.UserGetAddressesController())
	auth.POST("/account/add-address", usercontroller.UserCreateAddressController())
	auth.DELETE("/account/remove-address", usercontroller.RemoveUserAddress())

	// Password manage
	auth.POST("/account/req-forgotpass", usercontroller.UserForgotPasswordAuth())

	// Cart Management
	auth.GET("/cart/getall", usercontroller.GetAllCartProducts())
	auth.GET("/cart/add-pro/:productid", usercontroller.UserAddToCartController())
	auth.DELETE("/cart/delete/:id", usercontroller.UserRemoveCartController())
	auth.GET("/cart/quntity", usercontroller.UserCartQuntityController())

	// Checkout Section
	auth.POST("/single-checkout", usercontroller.UserSingleCheckout())
	auth.POST("/cart-checkout/:method/:coupon", usercontroller.CartCheckout())
	auth.POST("/verify-pay", usercontroller.VerifyOrder())
	// auth.POST("/online-checkout", usercontroller.UserOnlineCheckout())

	// Manage orders section
	auth.GET("/orders/getall", usercontroller.UserGetAllOrders())
	auth.POST("/orders/cancel", usercontroller.CancelOrder())
	// auth.GET("/orders/return", usercontroller.ReturnOrder())

	// Wallet Controll
	auth.GET("/wallet/getbal", usercontroller.GetWalletBalance())

	// Search product
	auth.GET("/searchbyid/:productid", usercontroller.SearchProductUsingID())
	auth.GET("/searchbycato/:catogary", usercontroller.SearchProductsUsingCatogary())
}
