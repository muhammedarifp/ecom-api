package helpers

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/muhammedarif/Ecomapi/config"
	"github.com/muhammedarif/Ecomapi/models"
	"github.com/muhammedarif/Ecomapi/private"
	"github.com/razorpay/razorpay-go"
)

// This is a -- CreateOrder -- helper function
// In razorpay we will create a order first and send order details on server
// So this helper function used for create order
func CreateOrder(productDatas *models.Product, total_amount float64) (bool, *models.Transactions) {
	db := *config.GetDb()
	client := razorpay.NewClient(private.RAZORPAY_KEY, private.RAZORPAY_SECRET)
	// tottal_mount := productData.Price
	recipt_id := uuid.New()
	newOrder := map[string]interface{}{
		"amount":   total_amount,
		"currency": "INR",
		"receipt":  recipt_id,
	}

	// Create order using razorpay api
	body, err := client.Order.Create(newOrder, nil)
	if err != nil {
		return false, &models.Transactions{}
	}

	// Create new transaction on database
	newTransaction := models.Transactions{
		OrderID:  body["id"].(string),
		Amount:   total_amount,
		Attempts: body["attempts"].(float64),
		Currency: body["currency"].(string),
		Receipt:  body["receipt"].(string),
		Status:   "Pending",
	}

	fmt.Println(newTransaction)
	if res := db.Create(&newTransaction); res.Error != nil {
		return false, &models.Transactions{}
	}
	return true, &newTransaction
}
