package usercontroller

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/muhammedarif/Ecomapi/config"
	"github.com/muhammedarif/Ecomapi/helpers"
	"github.com/muhammedarif/Ecomapi/models"
	"github.com/razorpay/razorpay-go"
)

var (
	RAZ_KEY = "rzp_test_hjvilmQsgbSYDR"
	RAZ_SEC = "IQu2ZsJiZBHWi3703mDrRmxd"
)

type CheckoutRequest struct {
	ProductID      string `json:"productID"`
	Method         string `json:"method"`
	IsCartCheckout bool   `json:"iscart"`
	Coupon         string `json:"coupon"`
}

func UserSingleCheckout() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// mode := ctx.Query("mode")
		userID := helpers.GetUserIDFromJwt(ctx)

		// Check user blocked or not
		if !helpers.IsUserBlocked(userID) {
			ctx.AbortWithStatusJSON(403, gin.H{
				"Error": "Your account was blocked",
			})
			return
		}

		// This is a user inp struct
		var data struct {
			ProductID uint   `json:"productID"`
			Method    string `json:"method"`
			Coupon    string `json:"coupon"`
		}

		if err := ctx.ShouldBindJSON(&data); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			return
		}

		address, stwhether := getDefaultAddress(userID)
		if !stwhether {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"Status": false,
				"Error":  "Order not completed | Add Delivery address",
			})
			return
		}

		if data.Method == "COD" {
			// TODO : PERFORM COD CHECKOUT

			whether, orderData, error := SingleCheckoutWithCod(data.ProductID, userID, data.Coupon)
			final_price := orderData.TottalAmount

			if data.Coupon != "" {
				price, err := UseCoupon(data.Coupon, orderData.TottalAmount)
				if err != nil {
					final_price = price
				}
			}

			// If order not placed that time server throw this exeption message
			if !whether && error != nil {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"Status": false,
					"Error":  "Order not completed",
				})
				return
			}

			// If Order is success.
			ctx.JSON(http.StatusOK, gin.H{
				"Message": "Order success",
				"Error":   nil,
				"Order": map[string]any{
					"OrderID":     orderData.ID,
					"Create":      orderData.CreatedAt,
					"Tottal":      final_price,
					"Status":      orderData.Status,
					"PayMethod":   orderData.PayMethod,
					"Transaction": nil,
				},
				"User": map[string]any{
					"UserID": userID,
					"Delivery": map[string]any{
						"Name":     address.Name,
						"Mobile":   address.Mobile,
						"Pincode":  address.Pincode,
						"State":    address.State,
						"City":     address.City,
						"Address":  address.Address,
						"Landmark": address.Landmark,
					},
				},
			})

			// If user select online payment
			// This condition i integrate razorpay
			//
		} else if data.Method == "ONLINE" {
			// TODO : PERFORM ONLINE PAY CHECKOUT
			whether, orderData, transactionData, err := SingleCheckoutWithOnline(fmt.Sprint(data.ProductID), userID)
			if err != nil && !whether {
				ctx.AbortWithStatusJSON(400, gin.H{
					"Error": err.Error(),
				})
				return
			}

			final_price := orderData.TottalAmount

			if data.Coupon != "" {
				price, err := UseCoupon(data.Coupon, orderData.TottalAmount)
				if err != nil {
					final_price = price
				}
			}

			ctx.JSON(http.StatusOK, gin.H{
				"Message": "Order success",
				"Error":   nil,
				"Order": map[string]any{
					"OrderID":     orderData.ID,
					"Create":      orderData.CreatedAt,
					"Tottal":      final_price,
					"Status":      orderData.Status,
					"PayMethod":   orderData.PayMethod,
					"Transaction": transactionData,
				},
				"User": map[string]any{
					"UserID": userID,
					"Delivery": map[string]any{
						"Name":     address.Name,
						"Mobile":   address.Mobile,
						"Pincode":  address.Pincode,
						"State":    address.State,
						"City":     address.City,
						"Address":  address.Address,
						"Landmark": address.Landmark,
					},
				},
			})

		} else if data.Method == "WALLET" {
			db := *config.GetDb()
			var proData models.Products
			if res := db.First(&proData, data.ProductID); res.Error != nil {
				ctx.AbortWithStatusJSON(400, gin.H{
					"Error": "Checkout failed",
				})
				return
			}
			if !useWallet(float64(proData.Price), ctx) {
				ctx.AbortWithStatusJSON(400, gin.H{
					"Error": "Checkout failed",
				})
				return
			}

			final_price := proData.Price
			if data.Coupon != "" {
				price, err := UseCoupon(data.Coupon, float64(final_price))
				if err != nil {
					final_price = uint(price)
				}
			}

			newOrder := models.Orders{
				UserID:       userID,
				TottalAmount: float64(final_price),
				Status:       "Success",
				PayMethod:    "WALLET",
				IsSuccess:    true,
			}
			newOrderItem := models.OrdersItems{
				OrderID:   newOrder.ID,
				ProductID: proData.ID,
				Quntity:   1,
				Price:     float64(final_price),
			}
			tx := db.Begin()
			if res := tx.Create(&newOrder); res.Error != nil {
				ctx.JSON(400, gin.H{
					"Error": res.Error.Error(),
				})
				return
			}
			if res := tx.Create(&newOrderItem); res.Error != nil {
				ctx.JSON(400, gin.H{
					"Error": res.Error.Error(),
				})
				return
			}
			if res := tx.Commit(); res.Error != nil {
				ctx.JSON(400, gin.H{
					"Error": res.Error.Error(),
				})
				return
			}

			ctx.JSON(http.StatusOK, gin.H{
				"Message": "Order success",
				"Error":   nil,
				"Order": map[string]any{
					"OrderID":     newOrder.ID,
					"Create":      newOrder.CreatedAt,
					"Tottal":      final_price,
					"Status":      newOrder.Status,
					"PayMethod":   newOrder.PayMethod,
					"Transaction": newOrder,
				},
				"User": map[string]any{
					"UserID": userID,
					"Delivery": map[string]any{
						"Name":     address.Name,
						"Mobile":   address.Mobile,
						"Pincode":  address.Pincode,
						"State":    address.State,
						"City":     address.City,
						"Address":  address.Address,
						"Landmark": address.Landmark,
					},
				},
			})
		} else {
			// This case user enter invalid input
			// So This case return error / invalid message from a user
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"Error": "Invalid input",
			})
		}
	}
}

// ! ================================================================================================

func CartCheckout() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := helpers.GetUserIDFromJwt(ctx)
		db := *config.GetDb()
		method := ctx.Param("method")
		coupon := ctx.Param("coupon")
		method = strings.ToUpper(method)

		address, stwhether := getDefaultAddress(userID)
		if !stwhether {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"Status": false,
				"Error":  "Order not completed | Add Delivery address",
			})
			return
		}

		// Check user blocked or not
		if !helpers.IsUserBlocked(userID) {
			ctx.AbortWithStatusJSON(403, gin.H{
				"Error": "Your account was blocked",
			})
			return
		}

		// Find user products and tottal price
		// This all write using sql quries
		total_price := 0
		var CartItems []models.UserCart
		db.Table("user_carts").
			Joins("JOIN products ON products.id = user_carts.product_id").
			Where("user_carts.user_id = ?", userID).
			Find(&CartItems).
			Select("SUM(price * product_count) as tottal").
			Find(&total_price)

		if coupon != "" {
			price, err := UseCoupon(coupon, float64(total_price))
			if err == nil {
				total_price = int(price)
			}
		}

		if method == "COD" {
			// TODO : Cod Cart checkout section

			// //? First step is create new order
			newOrder := models.Orders{
				UserID:       userID,
				TottalAmount: float64(total_price),
				Status:       "Success",
				PayMethod:    "COD",
				IsSuccess:    true,
			}

			if res := db.Create(&newOrder); res.Error != nil {
				ctx.JSON(400, gin.H{
					"Error": "Checkout not completed",
				})
				return
			}

			// New OrderItems
			for _, value := range CartItems {
				var _productDeta models.Products
				if res := db.First(&_productDeta, value.ProductID); res.Error != nil {
					ctx.AbortWithStatusJSON(400, gin.H{
						"Error": "Checkout failed",
					})
					return
				}
				intUserID, _ := strconv.Atoi(userID)
				newItem := models.OrdersItems{
					OrderID:   uint(intUserID),
					ProductID: value.ProductID,
					Quntity:   uint(value.ProductCount),
					Price:     float64(_productDeta.Price) * float64(value.ProductCount),
				}

				if res := db.Create(&newItem); res.Error != nil {
					ctx.AbortWithStatusJSON(400, gin.H{
						"Error": "Checkout failed",
					})
					return
				}
			}

			ctx.JSON(http.StatusOK, gin.H{
				"Message": "Order success",
				"Error":   nil,
				"Order": map[string]any{
					"OrderID":     newOrder.ID,
					"Create":      newOrder.CreatedAt,
					"Tottal":      newOrder.TottalAmount,
					"Status":      newOrder.Status,
					"PayMethod":   newOrder.PayMethod,
					"Transaction": nil,
				},
				"User": map[string]any{
					"UserID": userID,
					"Delivery": map[string]any{
						"Name":     address.Name,
						"Mobile":   address.Mobile,
						"Pincode":  address.Pincode,
						"State":    address.State,
						"City":     address.City,
						"Address":  address.Address,
						"Landmark": address.Landmark,
					},
				},
			})

		} else if method == "ONLINE" {

			// TODO : Online Cart checkout section

			whether, trData := CreateOrder(&models.Products{}, float64(total_price))
			if !whether {
				ctx.AbortWithStatusJSON(400, gin.H{
					"Error": "Checkout failed",
				})
				return
			}

			// Creete new order
			newOrder := models.Orders{
				UserID:        userID,
				TottalAmount:  float64(total_price),
				Status:        "Pending",
				PayMethod:     method,
				TransactionID: trData.ID,
				IsSuccess:     false,
			}

			if res := db.Create(&newOrder); res.Error != nil {
				ctx.JSON(400, gin.H{
					"Error": "Checkout not completed",
				})
				return
			}

			for _, value := range CartItems {
				var _productDeta models.Products
				if res := db.First(&_productDeta, value.ProductID); res.Error != nil {
					ctx.AbortWithStatusJSON(400, gin.H{
						"Error": "Checkout failed",
					})
					return
				}
				intUserID, _ := strconv.Atoi(userID)
				newItem := models.OrdersItems{
					OrderID:   uint(intUserID),
					ProductID: value.ProductID,
					Quntity:   uint(value.ProductCount),
					Price:     float64(_productDeta.Price) * float64(value.ProductCount),
				}

				if res := db.Create(&newItem); res.Error != nil {
					ctx.AbortWithStatusJSON(400, gin.H{
						"Error": "Checkout failed",
					})
					return
				}
			}

			ctx.JSON(http.StatusOK, gin.H{
				"Message": "Order success",
				"Error":   nil,
				"Order": map[string]any{
					"OrderID":     newOrder.ID,
					"Create":      newOrder.CreatedAt,
					"Tottal":      newOrder.TottalAmount,
					"Status":      newOrder.Status,
					"PayMethod":   newOrder.PayMethod,
					"Transaction": trData,
				},
				"User": map[string]any{
					"UserID": userID,
					"Address": map[string]any{
						"Name":     address.Name,
						"Mobile":   address.Mobile,
						"Pincode":  address.Pincode,
						"State":    address.State,
						"City":     address.City,
						"Address":  address.Address,
						"Landmark": address.Landmark,
					},
				},
			})

			//
		} else {
			// ! [error] Invalid input section
			ctx.AbortWithStatusJSON(400, gin.H{
				"Error": "Invalid Input",
			})
		}
	}
}

// ! =============================================================================================

// ? Helper Functions Section -----

func SingleCheckoutWithCod(productID uint, userID, coupon string) (bool, *models.Orders, error) {
	db := *config.GetDb()

	// Fetch product using product ID
	//
	var proData models.Products
	if res := db.First(&proData, productID); res.Error != nil {
		return false, &models.Orders{}, errors.New("internal servor error")
	}

	final_price := proData.Price
	if coupon != "" {
		price, err := UseCoupon(coupon, float64(final_price))
		if err == nil {
			final_price = uint(price)
		}
	}

	// Create order using user id
	newOrder := models.Orders{
		UserID:       userID,
		TottalAmount: float64(final_price),
		Status:       "Success",
		PayMethod:    "COD",
		IsSuccess:    true,
	}
	if res := db.Create(&newOrder); res.Error != nil {
		return false, &models.Orders{}, errors.New("internal servor error")
	}

	// Create order item
	newOrderItem := models.OrdersItems{
		OrderID:   newOrder.ID,
		ProductID: productID,
		Quntity:   1,
		Price:     float64(final_price),
	}
	if res := db.Create(&newOrderItem); res.Error != nil {
		return false, &models.Orders{}, errors.New("internal servor error")
	}

	return false, &newOrder, nil
}

// Online single checkout helper function
func SingleCheckoutWithOnline(productID, UserID string) (bool, *models.Orders, *models.Transactions, error) {
	db := *config.GetDb()
	var proData models.Products
	if res := db.First(&proData, productID); res.Error != nil {
		return false, &models.Orders{}, &models.Transactions{}, errors.New("internal server error or invalid input")
	}

	// Create razorpay order
	whether, transactionData := CreateOrder(&proData, float64(proData.Price))
	if !whether {
		return false, &models.Orders{}, &models.Transactions{}, errors.New("checkout failed, order not created")
	}

	// Create new order
	newOrder := models.Orders{
		UserID:        UserID,
		TottalAmount:  float64(transactionData.Amount),
		Status:        "PENDING",
		PayMethod:     "ONLINE",
		TransactionID: transactionData.ID,
		IsSuccess:     false,
	}

	if res := db.Create(&newOrder); res.Error != nil {
		return false, &models.Orders{}, &models.Transactions{}, errors.New("checkout failed")
	}
	// New Order Item.

	newOrderItem := models.OrdersItems{
		OrderID:   newOrder.ID,
		ProductID: proData.ID,
		Quntity:   1,
		Price:     float64(proData.Price),
	}

	if res := db.Create(&newOrderItem); res.Error != nil {
		return false, &models.Orders{}, &models.Transactions{}, errors.New("checkout failed")
	}

	// If This all are success. Your order is success
	// That case server return success message
	return true, &newOrder, transactionData, nil
}

// Finally End Cart Checkout Section
// After i continue cart checkout section
// This  Function Recive array of cart ids
// TODO : Cart checkout section

// ! ==========================================================================================

// ! ====================================================================================

// TODO : Verify checkout / payment section
type VerifyInps struct {
	OrderID   string `json:"order_id"`
	PaymentID string `json:"payment_id"`
	Signature string `json:"signature"`
}

func VerifyOrder() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		db := *config.GetDb()
		userID := helpers.GetUserIDFromJwt(ctx)
		var inpData VerifyInps
		if err := ctx.ShouldBindJSON(&inpData); err != nil {
			ctx.JSON(400, gin.H{
				"Error": "Order not success",
			})
			return
		}

		var transactionData models.Transactions
		if res := db.First(&transactionData, inpData.OrderID); res.Error != nil {
			ctx.JSON(400, gin.H{
				"Error": "Transaction not found",
			})
			return
		}

		var data = transactionData.OrderID + "|" + inpData.PaymentID

		//
		hash := hmac.New(sha256.New, []byte(RAZ_SEC))
		hash.Write([]byte(data))
		genarated_signature := hex.EncodeToString(hash.Sum(nil))

		if inpData.Signature != genarated_signature {
			ChangeOrderStatus(userID, inpData.OrderID)
			transactionData.Attempts++
			if res := db.Save(&transactionData); res.Error != nil {
				ctx.JSON(200, gin.H{
					"Error": "Transaction not verified | Internal server error",
				})
			}
			ctx.JSON(200, gin.H{
				"Error": "Transaction not verified",
			})
		} else {
			ctx.JSON(200, gin.H{
				"Message": "Transaction verified",
			})
		}
	}
}

// ! ==================================================================================================

func ChangeOrderStatus(userID, orderID string) bool {
	db := *config.GetDb()
	var orderData models.Orders
	if res := db.Where("id = ? AND user_id = ?", orderID, userID).First(&orderData); res.Error != nil {
		return false
	}

	orderData.IsSuccess = true
	orderData.Status = "Success"

	if res := db.Save(&orderData); res.Error != nil {
		return false
	}

	return true
}

// func

// ! =====================================================================================

// This is a -- CreateOrder -- helper function
// In razorpay we will create a order first and send order details on server
// So this helper function used for create order
func CreateOrder(productDatas *models.Products, total_amount float64) (bool, *models.Transactions) {
	db := *config.GetDb()
	client := razorpay.NewClient(RAZ_KEY, RAZ_SEC)
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

// ? ----------------------------------------------------------------------------------
// Coupon Using

func UseCoupon(coupon string, price float64) (float64, error) {
	db := *config.GetDb()
	var couponData models.Coupons
	if res := db.First(&couponData, `code = ?`, coupon); res.Error != nil {
		return 0, errors.New("invalid coupon code")
	}

	duration := time.Now().Sub(couponData.CreatedAt)
	days := int(duration.Hours() / 24)
	if days > couponData.ExpiryDate {
		return 0, errors.New("coupon expired")
	}

	if price < couponData.MinAmount {
		return 0, errors.New("minimum amount required")
	}

	discount := price * (couponData.Discount / 100)
	if discount > couponData.MaxDiscountAmount {
		discount = couponData.MaxDiscountAmount
	}

	final_price := price - discount
	return final_price, nil

}

// ? ----------------------------------------------------------------------------------
// This Function maily focus about get user user default address

func getDefaultAddress(userID string) (*models.Address, bool) {
	db := *config.GetDb()
	var addressData models.Address
	if res := db.Where("user_id = ? AND is_default = true", userID).First(&addressData); res.Error != nil {
		return &addressData, false
	}

	return &addressData, true
}
