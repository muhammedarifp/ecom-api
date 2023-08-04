package tests

// import (
// 	"bytes"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/gin-gonic/gin"
// 	usercontroller "github.com/muhammedarif/Ecomapi/controller/user"
// 	"github.com/muhammedarif/Ecomapi/models"
// )

// func TestUserLoginSuccess(t *testing.T) {
// 	router := gin.Default()
// 	router.GET("/login", usercontroller.UserLoginController())

// 	userData := models.UserLoginForm{
// 		Email:    "hello@gmail.com",
// 		Password: "12345",
// 	}
// 	payload, _ := json.Marshal(userData)

// 	req, _ := http.NewRequest("POST", "/login", bytes.NewReader(payload))
// 	req.Header.Set("Content-Type", "application/json")

// 	rec := httptest.NewRecorder()
// 	rec.
// }
