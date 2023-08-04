package admincontroller

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/muhammedarif/Ecomapi/config"
	"github.com/muhammedarif/Ecomapi/models"
)

func SalesGraph() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		db := *config.GetDb()
		var orders []models.Order
		current_time := time.Now()
		prev_date := current_time.AddDate(0, 0, -30)

		var monthlySales []struct {
			Day        int
			Month      int
			Year       int
			TotalSales int
		}

		query := `
        SELECT
            TO_CHAR(created_at, 'DD') AS day,
			TO_CHAR(created_at, 'MM') AS month,
            TO_CHAR(created_at, 'YYYY') AS year,
            SUM(tottal_amount) AS total_sales
        FROM
			orders
        GROUP BY
            TO_CHAR(created_at, 'DD'),
			TO_CHAR(created_at, 'MM'),
            TO_CHAR(created_at, 'YYYY')
    `

		db.Raw(query).Scan(&monthlySales)

		db.Find(&orders, "created_at >= ?", prev_date)
		ctx.JSON(200, monthlySales)

	}
}
