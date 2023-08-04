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

		method := strings.ToUpper(data.Method)
		fmt.Println(method)

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
		if method == "COD" {
			// TODO : PERFORM COD CHECKOUT
			SingleCheckoutWithCod(data.ProductID, userID, data.Coupon, ctx)

		} else if method == "ONLINE" {
			// TODO : PERFORM ONLINE PAY CHECKOUT
			SingleCheckoutWithOnline(fmt.Sprint(data.ProductID), userID, data.Coupon, ctx)

		} else if method == "WALLET" {
			db := *config.GetDb()
			var proData models.Product
			if res := db.First(&proData, data.ProductID); res.Error != nil {
				ctx.AbortWithStatusJSON(400, gin.H{
					"Error": "Checkout failed",
				})
				return
			}
			if !useWallet(float64(proData.Price), ctx) {
				ctx.AbortWithStatusJSON(400, gin.H{
					"Error": "Balance is over. topup your wallet",
				})
				return
			}

			final_price := proData.Price
			if data.Coupon != "" {
				price, err := UseCoupon(data.Coupon, userID, float64(final_price))
				if err == nil {
					final_price = price
				}
			}

			userIDint, _ := strconv.Atoi(userID)
			newOrder := models.Order{
				UserID:       uint(userIDint),
				TottalAmount: float64(final_price),
				Status:       "Success",
				PayMethod:    "WALLET",
				IsSuccess:    true,
				Items: []models.OrderItem{
					{ProductID: proData.ID, Quntity: 1, Price: final_price},
				},
			}

			tx := db.Begin()
			if res := tx.Create(&newOrder); res.Error != nil {
				ctx.JSON(400, gin.H{
					"Error": res.Error.Error(),
					"Num":   1,
				})
				return
			}

			if res := tx.Commit(); res.Error != nil {
				ctx.JSON(400, gin.H{
					"Error": res.Error.Error(),
					"Num":   3,
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
					"Coupon": map[string]any{
						"Coupon":   data.Coupon,
						"Discount": proData.Price - final_price,
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
		} else {
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

		if total_price <= 0 {
			ctx.JSON(400, gin.H{
				"Error": "Cart is empty",
			})
			ctx.Abort()
			return
		}

		final_price := total_price
		if data.Coupon != "" {
			price, err := UseCoupon(data.Coupon, userID, float64(total_price))
			if err == nil {
				fmt.Println(price)
				final_price = price
			} else {
				fmt.Println(err.Error())
			}
		}

		if method == "COD" {
			// TODO : Cod Cart checkout section

			// //? First step is create new order
			userIDint, _ := strconv.Atoi(userID)

			// New OrderItems
			var items []models.OrderItem
			for _, value := range CartItems {
				var _productDeta models.Product
				if res := db.First(&_productDeta, value.ProductID); res.Error != nil {
					ctx.AbortWithStatusJSON(400, gin.H{
						"Error": "Checkout failed",
					})
					return
				}
				// intUserID, _ := strconv.Atoi(userID)
				newItem := models.OrderItem{
					ProductID: value.ProductID,
					Quntity:   uint(value.ProductCount),
					Price:     float64(_productDeta.Price) * float64(value.ProductCount),
				}

				items = append(items, newItem)
			}

			newOrder := models.Order{
				UserID:       uint(userIDint),
				TottalAmount: float64(final_price),
				Status:       "Success",
				PayMethod:    "COD",
				IsSuccess:    true,
				Items:        items,
			}

			if res := db.Create(&newOrder); res.Error != nil {
				ctx.JSON(400, gin.H{
					"Error": "Checkout not completed",
				})
				return
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
						"Code":         data.Coupon,
						"OrgPrice":     total_price,
						"DiscoutPrice": final_price,
						"Discount":     total_price - final_price,
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

			whether, trData := helpers.CreateOrder(&models.Product{}, final_price*100)
			if !whether {
				ctx.AbortWithStatusJSON(400, gin.H{
					"Error": "Checkout failed",
				})
				return
			}

			var items []models.OrderItem
			for _, value := range CartItems {
				var _productDeta models.Product
				if res := db.First(&_productDeta, value.ProductID); res.Error != nil {
					ctx.AbortWithStatusJSON(400, gin.H{
						"Error": "Checkout failed",
					})
					return
				}
				newItem := models.OrderItem{
					ProductID: value.ProductID,
					Quntity:   uint(value.ProductCount),
					Price:     float64(_productDeta.Price) * float64(value.ProductCount),
				}

				items = append(items, newItem)
			}

			// Creete new order
			userIDint, _ := strconv.Atoi(userID)
			newOrder := models.Order{
				UserID:        uint(userIDint),
				TottalAmount:  float64(final_price),
				Status:        "Pending",
				PayMethod:     method,
				TransactionID: trData.ID,
				IsSuccess:     false,
				Items:         items,
			}

			if res := db.Create(&newOrder); res.Error != nil {
				ctx.JSON(400, gin.H{
					"Error": "Checkout not completed",
				})
				return
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
func SingleCheckoutWithCod(productID uint, userID, coupon string, ctx *gin.Context) bool {
	fmt.Println("Iam")
	db := *config.GetDb()

	address, isFetch := getDefaultAddress(userID)

	if !isFetch {
		ctx.AbortWithStatusJSON(400, gin.H{"Error": "Delivery address not found"})
	}

	var proData models.Product
	if res := db.First(&proData, productID); res.Error != nil {
		return false
	}

	final_price := proData.Price
	if coupon != "" {
		price, err := UseCoupon(coupon, userID, float64(final_price))
		if err == nil {
			final_price = price
		}
	}

	var catogoryOffer models.CatogoryOffer
	if res := db.Preload("Catogory").Find(&catogoryOffer, proData.CatogaryID); res.Error == nil {
		percentage := (catogoryOffer.Dicount / final_price) * final_price
		if percentage > catogoryOffer.MaxDiscount {
			final_price = final_price - catogoryOffer.MaxDiscount
		} else {
			final_price = final_price - percentage
		}
	}

	// Create order using user id
	userIDint, _ := strconv.Atoi(userID)
	newOrder := models.Order{
		UserID:        uint(userIDint),
		TottalAmount:  float64(final_price),
		Status:        "Success",
		PayMethod:     "COD",
		IsSuccess:     true,
		TransactionID: 0,
		Items: []models.OrderItem{
			{
				ProductID: productID,
				Quntity:   1,
				Price:     final_price,
			},
		},
	}
	if res := db.Create(&newOrder); res.Error != nil {
		return false
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
				"Coupon":   coupon,
				"Cashback": proData.Price - final_price,
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

	return true
}

// This function purpus is hep Single checkout section
// This function recive productID, userID and coupon code
// Fetch product details using product id
// coupon pending
// After that Create Order using previusly fetched order details
// Finally create order items
// Return that all details
func SingleCheckoutWithOnline(productID, UserID, coupon string, ctx *gin.Context) bool {
	db := *config.GetDb()
	address, isFetch := getDefaultAddress(UserID)
	if !isFetch {
		ctx.AbortWithStatusJSON(400, gin.H{"Error": "Delivery address not found"})
	}
	var proData models.Product
	if res := db.First(&proData, productID); res.Error != nil {
		return false
	}

	final_price := proData.Price
	if coupon != "" {
		fprice, err := UseCoupon(coupon, UserID, final_price)
		if err == nil {
			final_price = fprice
		}
	}

	var catogoryOffer models.CatogoryOffer
	if res := db.Preload("Catogory").Find(&catogoryOffer, proData.CatogaryID); res.Error == nil {
		percentage := (catogoryOffer.Dicount / final_price) * final_price
		fmt.Println(percentage)
		if percentage > catogoryOffer.MaxDiscount {
			final_price = final_price - catogoryOffer.MaxDiscount
		} else {
			final_price = final_price - percentage
		}
	}

	// Create razorpay order
	whether, transactionData := helpers.CreateOrder(&proData, final_price*100)
	if !whether {
		return false
	}

	// Create new order
	userIDint, _ := strconv.Atoi(UserID)
	newOrder := models.Order{
		UserID:        uint(userIDint),
		TottalAmount:  final_price,
		Status:        "PENDING",
		PayMethod:     "ONLINE",
		TransactionID: transactionData.ID,
		IsSuccess:     false,
		Items: []models.OrderItem{
			{ProductID: proData.ID, Quntity: 1, Price: final_price},
		},
	}

	if res := db.Create(&newOrder); res.Error != nil {
		return false
	}

	// If This all are success. Your order is success
	// That case server return success message
	ctx.JSON(http.StatusOK, gin.H{
		"Message": "Order success",
		"Error":   nil,
		"Order": map[string]any{
			"OrderID":     newOrder.ID,
			"Create":      newOrder.CreatedAt,
			"Tottal":      final_price,
			"T_Test":      newOrder.TottalAmount,
			"Status":      newOrder.Status,
			"PayMethod":   newOrder.PayMethod,
			"Transaction": transactionData,
			"Coupon": map[string]any{
				"Coupon":   coupon,
				"Cashback": proData.Price - final_price,
			},
		},
		"User": map[string]any{
			"UserID": UserID,
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

	return true
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
	var couponData models.Coupon
	var couponUsage models.CouponUsage
	if res := db.First(&couponData, `code = ?`, coupon); res.Error != nil {
		return 0, errors.New("invalid coupon code")
	}

	duration := couponData.CreatedAt.Sub(time.Now())

	if duration.Hours() > float64(couponData.ExpiryDate)*24 {
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

	couponUsage.UsageCount++
	db.Save(&couponUsage)
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
