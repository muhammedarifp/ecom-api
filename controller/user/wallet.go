package usercontroller

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/muhammedarif/Ecomapi/config"
	"github.com/muhammedarif/Ecomapi/helpers"
	"github.com/muhammedarif/Ecomapi/models"
	"gorm.io/gorm"
)

// In initial case user have no wallet
// So this function is focus to create wallet
func CreateWallet(userID uint) bool {
	db := *config.GetDb()
	newWallet := models.Wallets{
		UserID:  userID,
		Balance: 0,
	}

	if res := db.Create(&newWallet); res.Error != nil {
		return false
	}

	return true
}

func GetWalletBalance() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		db := *config.GetDb()
		userID := helpers.GetUserIDFromJwt(ctx)
		var walletData models.Wallets
		if res := db.First(&walletData, `user_id = ?`, userID); res.Error != nil {

			// user have no wallet
			// Create new wallet
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				userIDint, _ := strconv.Atoi(userID)
				if whether := CreateWallet(uint(userIDint)); !whether {
					ctx.AbortWithStatusJSON(400, gin.H{
						"Error": "Fetch balance failed",
					})
					return
				}

				ctx.JSON(200, gin.H{
					"Message": "New wallet created",
					"Balance": 0.0,
				})
				return
			}

			ctx.AbortWithStatusJSON(400, gin.H{
				"Error": res.Error.Error(),
			})

			return
		}

		ctx.JSON(200, gin.H{
			"Error": nil,
			"Wallet": map[string]any{
				"Create":        walletData.CreatedAt,
				"WalletID":      walletData.ID,
				"WalletBalance": walletData.Balance,
			},
		})
	}
}

func useWallet(price float64, ctx *gin.Context) bool {
	userID := helpers.GetUserIDFromJwt(ctx)
	db := *config.GetDb()
	var walletData models.Wallets
	if res := db.First(&walletData, `user_id = ?`, userID); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			intuserID, err := strconv.Atoi(userID)
			if err != nil {
				return false
			}
			CreateWallet(uint(intuserID))
		}
		return false
	}

	if price <= walletData.Balance {
		walletData.Balance = walletData.Balance - price
		db.Save(&walletData)
		return true
	}
	return false
}
