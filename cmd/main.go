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
		&models.Products{},
		&models.Catogary{},
		&models.ProductImages{},
		&models.Address{},
		&models.UserCart{},
		&models.Orders{},
		&models.OrdersItems{},
		&models.Transactions{},
		&models.Coupons{},
		&models.Wallets{},
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
