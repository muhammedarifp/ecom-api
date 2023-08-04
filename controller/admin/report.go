package admincontroller

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	"github.com/muhammedarif/Ecomapi/config"
	"github.com/muhammedarif/Ecomapi/models"
)

type Sale struct {
	Date        string
	ProductName string
	Quantity    int
	TotalAmount float64
}

func DownloadSalesReport() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		db := *config.GetDb()
		var orders []models.Order
		var tottal_income float64
		db.Table("orders").Where("status = 'Success'").Select("SUM(tottal_amount)").Take(&tottal_income)
		db.Preload("Items").Find(&orders)

		pdf := gofpdf.New("P", "mm", "A4", "")
		pdf.AddPage()

		pdf.SetFont("Arial", "B", 20)
		pdf.Cell(100, 10, "SALES REPORT")
		pdf.Ln(15)

		//
		pdf.SetFont("Arial", "B", 13)
		pdf.SetFillColor(51, 255, 51)
		pdf.CellFormat(60, 10, fmt.Sprintf("Total Income : %.2f", tottal_income), "", 0, "", true, 0, "")
		pdf.SetXY(73, 25)
		pdf.CellFormat(50, 10, fmt.Sprintf("New Customers : %d", 100), "", 0, "", true, 0, "")
		pdf.Ln(10)

		//
		pdf.Cell(100, 20, "Top Orders ")
		pdf.Ln(15)

		pdf.SetFont("Arial", "", 10)
		pdf.Cell(30, 10, "OrderID")
		pdf.Cell(30, 10, "Tottal Price")
		pdf.Cell(30, 10, "Name")
		pdf.Cell(30, 10, "Status")
		pdf.Ln(10)

		pdf.SetFont("Arial", "", 10)
		for _, order := range orders {
			if order.Status == "Success" {
				pdf.SetFillColor(153, 255, 153)
			} else {
				pdf.SetFillColor(255, 204, 153)
			}
			pdf.CellFormat(30, 5, fmt.Sprint(order.ID), "", 0, "", true, 0, "")
			pdf.CellFormat(30, 5, fmt.Sprint(order.TottalAmount), "", 0, "", true, 0, "")
			for k, item := range order.Items {
				db.First(&item.Product, item.ProductID)
				if k == 0 {
					pdf.CellFormat(30, 5, fmt.Sprint(item.Product.Name), "", 0, "", true, 0, "")
					pdf.CellFormat(30, 5, fmt.Sprint(order.Status), "", 0, "", true, 0, "")
					pdf.Ln(6)
				} else {
					pdf.CellFormat(30, 5, "sub", "", 0, "", false, 0, "")
					pdf.CellFormat(30, 5, "sub", "", 0, "", false, 0, "")
					pdf.CellFormat(30, 5, fmt.Sprint(item.Product.Name), "", 0, "", true, 0, "")
					pdf.CellFormat(30, 5, fmt.Sprint(order.Status), "", 0, "", true, 0, "")
					pdf.Ln(6)
				}

			}
			pdf.Ln(6)
		}

		pdf.Output(ctx.Writer)
	}
}
