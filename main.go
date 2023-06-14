package ecom

import "github.com/gin-gonic/gin"

func main() {
	// This is just create defult gin engine
	router := gin.Default()

	// Create user route using "/" this
	user := router.Group("/")

	//
	admin := router.Group("/admin")
	// admin_route.AdminRoutes(admin)

	// server run at default port 8080
	router.Run()
}
