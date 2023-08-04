package main

import (
	"github.com/gin-gonic/gin"

	"github.com/muhammedarif/Ecomapi/config"
	"github.com/muhammedarif/Ecomapi/middleware"
	"github.com/muhammedarif/Ecomapi/models"
	adminroutes "github.com/muhammedarif/Ecomapi/routes/admin"
	userroutes "github.com/muhammedarif/Ecomapi/routes/user"
)

func init() {
	config.InitDataBase()
}

func main() {

	// Migrate all database models using gorm auto migrate
	config.DataBase.AutoMigrate(
		&models.Users{},
		&models.Product{},
		&models.ProductImages{},
		&models.Address{},
		&models.UserCart{},
		&models.Order{},
		&models.OrderItem{},
		&models.Transactions{},
		&models.Coupon{},
		&models.CouponUsage{},
		&models.Catogory{},
		&models.Wallets{},
		&models.CatogoryOffer{},
	)

	// Create Default gin engine
	router := gin.Default()

	adminauth := router.Group("/")
	userauth := router.Group("/")
	adminauth.Use(middleware.AuthAdminMiddleWare())
	userauth.Use(middleware.AuthUserMiddleWare())
	unauth := router.Group("/")

	adminroutes.AdminRoutes(adminauth, unauth)
	userroutes.UserRoutes(userauth, unauth)

	// App Run default localhost post 8080
	router.Run(":8090")
}
