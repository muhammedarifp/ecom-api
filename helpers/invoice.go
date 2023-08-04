package helpers

import (
	"fmt"

	"github.com/jung-kurt/gofpdf"
	"github.com/muhammedarif/Ecomapi/models"
)

func CreateInvoice(order models.Order, address models.Address, user models.Users) **gofpdf.Fpdf {
	pdf := gofpdf.New("P", "mm", "A4", "")

	pdf.AddPage()
	// Set the document title
	pdf.SetTitle(fmt.Sprintf("ORDER_%d", order.ID), false)

	// Add a logo
	// pdf.AddImage("logo.png", 10, 10, 50, 50)

	// Add the header
	pdf.SetFont("Arial", "B", 18)
	pdf.Cell(40, 10, "INVOICE")

	pdf.Ln(5)

	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(255, 0, 0)
	pdf.Cell(10, 10, "Invoice Id : 12345")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(0, 0, 0)

	pdf.Cell(40, 10, "Order Date:")
	pdf.Cell(100, 10, order.CreatedAt.Format("2015/02/25"))
	pdf.Ln(5)

	// Add the client information
	pdf.Cell(40, 10, "Client Name:")
	pdf.Cell(100, 10, fmt.Sprintf("%s %s", user.FirstName, user.LastName))
	pdf.Ln(5)
	pdf.Cell(40, 10, "Delivery address:")
	pdf.Cell(100, 10, fmt.Sprintf("%s %s", address.Address, address.Locality))
	pdf.Ln(5)
	pdf.Cell(40, 10, "Phone:")
	pdf.Cell(100, 10, address.Mobile)
	pdf.Ln(5)
	pdf.Cell(40, 10, "Email:")
	pdf.Cell(100, 10, user.Email)

	pdf.Ln(5)

	// Add the date
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(40, 10, "Date:")
	pdf.Cell(100, 10, "2023-07-25")
	pdf.Ln(10)

	// Add the line items
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(40, 10, "Product")
	pdf.Cell(20, 10, "Qty")
	pdf.Cell(20, 10, "Unit Cost")
	pdf.Cell(20, 10, "Discount")
	pdf.Cell(20, 10, "Total")

	pdf.SetFont("Arial", "", 10)
	for _, item := range order.Items {
		pdf.Ln(5)
		//Name
		pdf.Cell(40, 10, item.Product.Name)
		//Description
		pdf.Cell(20, 10, fmt.Sprintf("%d", item.Quntity))
		// Unit price
		pdf.Cell(20, 10, fmt.Sprintf("%.2f", item.Product.Price))
		// unit price
		pdf.Cell(20, 10, "0.00")
		//
		pdf.Cell(20, 10, fmt.Sprintf("%.2f", item.Product.Price*float64(item.Quntity)))
	}
	pdf.Ln(7)

	// Add the subtotal
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(40, 10, "Subtotal:")
	pdf.Cell(20, 10, "")
	pdf.Cell(20, 10, "")
	pdf.Cell(20, 10, "")
	pdf.Cell(20, 10, fmt.Sprintf("%.2f", order.TottalAmount))
	pdf.Ln(20)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(1, 5, "Thanks For Buying")

	return &pdf
}
