package helpers

import (
	"errors"
	"fmt"
	"math/rand"
	"net/smtp"
	"sync"

	"github.com/muhammedarif/Ecomapi/config"
	"github.com/muhammedarif/Ecomapi/constants"
	"github.com/muhammedarif/Ecomapi/models"
)

// Initialize the cache
var (
	otpCache = sync.Map{}
)

func SendOtpEmail(userid string, otp int) (bool, error) {
	db := *config.GetDb()
	var userData models.Users
	if res := db.First(&userData, userid); res.Error != nil {
		return false, errors.New("user not found")
	}

	_username := "muhammedarif0100@gmail.com"
	_password := ""
	host := "smtp.gmail.com"
	auth := smtp.PlainAuth(
		"",
		_username,
		_password,
		host,
	)

	msg := fmt.Sprintf("From: muhammedarif0100@gmail.com\n"+
		"To: %s \n"+
		"Subject: Email Verification \n\r"+
		"Your One Time Password Is : %d \n\n Enter and submit this. This otp valid only next 2 minuts", userData.Email, otp)

	if err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		_username,
		[]string{userData.Email},
		[]byte(msg),
	); err != nil {
		return false, errors.New("otp not send")
	}

	otpCache.Store(userid, otp)

	return true, nil
}

func ValidateOtp(userotp, userid string) bool {

	otp, ok := otpCache.Load(userid)

	fmt.Println("okcc : ", otp, " ", userotp)
	if !ok {
		return false
	}

	otpString := fmt.Sprint(otp)

	if userotp == otpString {
		otpCache.Delete(userid)
		return true
	}

	return false
}

// Genarate One time pass

func GenarateOtp() int {
	min := 10000
	max := 99999
	otp := rand.Intn(max-min+1) + min
	return otp
}

func GetOtp(userid string) any {
	val, _ := constants.CACHE.Get(userid)

	return val
}
