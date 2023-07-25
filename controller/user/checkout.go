package usercontroller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/muhammedarif/Ecomapi/config"
	"github.com/muhammedarif/Ecomapi/helpers"
	"github.com/muhammedarif/Ecomapi/models"
	"gorm.io/gorm"
)

// var (
// 	RAZ_KEY = "rzp_test_hjvilmQsgbSYDR"
// 	RAZ_SEC = "IQu2ZsJiZBHWi3703mDrRmxd"
// )

// This section focus about to single product checkout
// This fuction also available cod and online
func UserSingleCheckout() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// Get userid from jwt token clims
		userID := helpers.GetUserIDFromJwt(ctx)

		// Check user blocked or not
		if !helpers.IsUserBlocked(userID) {
			ctx.AbortWithStatusJSON(403, gin.H{
				"Error": "Your account was blocked",
			})
			return
		}

		// This structure user enter data model
		var data struct {
			ProductID uint   `json:"productID"`
			Method    string `json:"method"`
			Coupon    string `json:"coupon"`
		}

		// Bind user enter data to data struct
		// for using simplicity
		if err := ctx.ShouldBindJSON(&data); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			return
		}

		// Get user defult address
		address, stwhether := getDefaultAddress(userID)
		if !stwhether {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"Status": false,
				"Error":  "Order not completed | Delivery address not available",
			})
			return
		}

		// If user select cod on single checkot that case perform cod checkout
		//
		if data.Method == "COD" {
			// TODO : PERFORM COD CHECKOUT

			whether, orderData, error := SingleCheckoutWithCod(data.ProductID, userID, data.Coupon)
			final_price := orderData.TottalAmount

			if data.Coupon != "" {
				price, err := UseCoupon(data.Coupon, userID, orderData.TottalAmount)
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
			whether, orderData, transactionData, err := SingleCheckoutWithOnline(fmt.Sprint(data.ProductID), userID, data.Coupon)
			if err != nil && !whether {
				ctx.AbortWithStatusJSON(400, gin.H{
					"Error": err.Error(),
				})
				return
			}

			final_price := orderData.TottalAmount

			if data.Coupon != "" {
				price, err := UseCoupon(data.Coupon, userID, orderData.TottalAmount)
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

			//
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
				price, err := UseCoupon(data.Coupon, userID, float64(final_price))
				if err != nil {
					final_price = price
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

// This is a cart checkout main function
// This function recive paymoth method and coupon code
func CartCheckout() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := helpers.GetUserIDFromJwt(ctx)
		db := *config.GetDb()

		// User input structure model
		var data struct {
			Method string `json:"method"`
			Coupon string `json:"coupon"`
		}

		ctx.ShouldBindJSON(&data)
		method := strings.ToUpper(data.Method)

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
		total_price := 0.0
		var CartItems []models.UserCart
		db.Table("user_carts").
			Joins("JOIN products ON products.id = user_carts.product_id").
			Where("user_carts.user_id = ?", userID).
			Find(&CartItems).
			Select("SUM(price * product_count) as tottal").
			Find(&total_price)

		dicount_price := 0.0
		if data.Coupon != "" {
			price, err := UseCoupon(data.Coupon, userID, float64(total_price))
			if err == nil {
				total_price = price
				dicount_price = price
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

			if res := db.Delete(&models.UserCart{}, `user_id = ?`, userID); res.Error != nil {
				ctx.JSON(400, gin.H{
					"Error": "Internal server error",
				})
				return
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
					"Coupon": map[string]any{
						"Code":     data.Coupon,
						"Discount": total_price - dicount_price,
					},
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

			whether, trData := helpers.CreateOrder(&models.Products{}, float64(total_price))
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

// This is single cod checkout helper function
// This function recive productID, userID and coupon
// after that fetch product details using product id
// Next step is if coupon available apply coupon otherwise nothing
// Next step is create order on orders table
// Next step is create order details on order items table
// Return that order details and items details
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
		price, err := UseCoupon(coupon, userID, float64(final_price))
		if err == nil {
			final_price = price
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

// This function purpus is hep Single checkout section
// This function recive productID, userID and coupon code
// Fetch product details using product id
// coupon pending
// After that Create Order using previusly fetched order details
// Finally create order items
// Return that all details
func SingleCheckoutWithOnline(productID, UserID, coupon string) (bool, *models.Orders, *models.Transactions, error) {
	db := *config.GetDb()
	var proData models.Products
	if res := db.First(&proData, productID); res.Error != nil {
		return false, &models.Orders{}, &models.Transactions{}, errors.New("internal server error or invalid input")
	}

	final_price := proData.Price
	if coupon != "" {
		fprice, err := UseCoupon(coupon, UserID, final_price)
		if err == nil {
			final_price = fprice
		}
	}

	// Create razorpay order
	whether, transactionData := helpers.CreateOrder(&proData, float64(proData.Price))
	if !whether {
		return false, &models.Orders{}, &models.Transactions{}, errors.New("checkout failed, order not created")
	}

	// Create new order
	newOrder := models.Orders{
		UserID:        UserID,
		TottalAmount:  final_price,
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

// This function purpous is use discount and offer coupons
// This function recive 2 params coupon code and current price
// This fuction basic working is fetch coupon details using coupon code
// After that check coupon expired or not, coupon table available coupon validity in days
// Next step is check min-amount is available or not
// Next step is calculate discount amount and that amount greter than max-amount
// Finally return final price -- !!final price means (org-price - discount)!!
func UseCoupon(coupon, userID string, price float64) (float64, error) {
	db := *config.GetDb()
	var couponData models.Coupons
	var couponUsage models.CouponUsages
	if res := db.First(&couponData, `code = ?`, coupon); res.Error != nil {
		return 0, errors.New("invalid coupon code")
	}

	if res := db.First(&couponUsage, `user_id = ?`, userID); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			newUsage := models.CouponUsages{
				UserID:     userID,
				CouponID:   couponData.ID,
				UsageCount: 1,
			}

			if res := db.Create(&newUsage); res.Error != nil {
				return 0, errors.New("Coupon not applied")
			}
		}

		return 0, errors.New("Coupon not applied")
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

// This function usage is fetch default address
// This function recive userID only
// After that fetch that user addresses on database
// Finally return defult address
// This functon main purpous is code reusability
func getDefaultAddress(userID string) (*models.Address, bool) {
	db := *config.GetDb()
	var addressData models.Address
	if res := db.Where("user_id = ? AND is_default = true", userID).First(&addressData); res.Error != nil {
		return &addressData, false
	}

	return &addressData, true
}
