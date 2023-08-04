package admincontroller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/muhammedarif/Ecomapi/config"
	"github.com/muhammedarif/Ecomapi/models"
	"gorm.io/gorm"
)

// Add new product controller. This function create new product on product table

// New product model sample struct
type NewProductFormData struct {
	Name       string  `form:"name" validate:"required"`
	Disc       string  `form:"disc" validate:"required"`
	Price      float64 `form:"price" validate:"required"`
	Quntity    uint    `form:"quntity" validate:"required"`
	CatogaryID uint    `form:"catogary_id" validate:"required"`
}

func AdminAddNewProductController() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var enterData NewProductFormData

		form, _ := ctx.MultipartForm()
		images, whether := form.File["images"]

		if len(images) <= 2 {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
				Status:  false,
				Message: "Add product attempt failed",
				Error:   "min image count is 2",
			})
			return
		}

		// This Statement Throw MultiformData Image Error
		// Send Error Status,Message and Error json model
		//
		if !whether {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
				Status:  false,
				Message: "Add product attempt failed",
				Error:   "Image not accessible",
			})
			return
		}

		// Bind user enterd data using a struct
		// if any case product data not binded error throw
		if err := ctx.ShouldBind(&enterData); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
				Status:  false,
				Message: "Add product attempt failed 2",
				Error:   err.Error(),
			})
			return
		}

		// Validate user data using validator package.
		validate := validator.New()
		if err := validate.Struct(&enterData); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
				Status:  false,
				Message: "Add product attempt failed 3",
				Error:   err.Error(),
			})
			return
		}

		// Create database instance
		db := *config.GetDb()

		// Check Catogary is is valid or not
		if res := db.First(&models.Catogory{}, enterData.CatogaryID); res != nil {
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
					Status:  false,
					Message: "Product adding failed 4",
					Error:   "Your enterd catogary is is not valid !",
				})
				return
			}
		}

		var newProduct = models.Product{
			Name:       enterData.Name,
			Disc:       enterData.Disc,
			Price:      enterData.Price,
			Quntity:    enterData.Quntity,
			CatogaryID: enterData.CatogaryID,
			IsActive:   true,
			IsDeleted:  false,
		}

		if res := db.Create(&newProduct); res.Error != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
				Status:  false,
				Message: "Product adding failed 5",
				Error:   res.Error.Error(),
			})
			return
		}

		if len(images) > 4 || len(images) < 2 {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
				Status:  false,
				Message: "Product adding failed",
				Error:   "Your product count images follow this rool - count > 4 || count < 2",
			})
			return
		}

		for _, image := range images {
			filenamesplit := strings.Split(image.Filename, ".")
			ext := filenamesplit[len(filenamesplit)-1]
			newnameTemp := uuid.New()
			newname := fmt.Sprint(newnameTemp)

			// destinationPath := filepath.Join("helloo", newname+"."+ext)

			err := ctx.SaveUploadedFile(image, "/home/arifu/Desktop/ecom/public/images/"+newname+"."+ext)
			if err != nil {
				fmt.Println(err)
			}

			newImage := models.ProductImages{
				ProductID: newProduct.ID,
				ImageName: newname + "." + ext,
				IsDefault: false,
			}

			if res := db.Create(&newImage); res.Error != nil {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
					Status:  false,
					Message: "Product not added",
					Error:   res.Error.Error(),
				})
				return
			}
		}

		ctx.JSON(http.StatusOK, gin.H{
			"Status":    true,
			"Message":   "Product added success",
			"Error":     nil,
			"ProductID": newProduct.ID,
		})
	}
}

// Edit old product controller. This function just update product details
type NewDetails struct {
	ProductID   uint    `form:"product_id" validate:"required"`
	NewName     string  `form:"name" validate:"required"`
	NewDisc     string  `form:"disc" validate:"required"`
	NewPrice    float64 `form:"price" validate:"required"`
	NewQuntity  uint    `form:"quntity" validate:"required"`
	NewCatogary uint    `form:"catogary" validate:"required"`
}

func AdminEditProductController() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var EnterData NewDetails
		if err := ctx.ShouldBind(&EnterData); err != nil {
			ctx.JSON(400, gin.H{
				"Error": "Product edit error",
			})
			return
		}
		db := *config.GetDb()
		var proDeta models.Product
		if res := db.First(&proDeta, EnterData.ProductID); res.Error != nil {
			ctx.JSON(400, gin.H{
				"Error": "Product edit error",
			})
			return
		}

		file, _ := ctx.MultipartForm()
		var oldImages []models.ProductImages
		images, _ := file.File["images"]

		proDeta.Name = EnterData.NewName
		proDeta.Disc = EnterData.NewDisc
		proDeta.Price = EnterData.NewPrice
		proDeta.Quntity = EnterData.NewQuntity
		proDeta.CatogaryID = EnterData.NewCatogary
		if res := db.Save(&proDeta); res.Error != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
				Status:  false,
				Message: "Product updation failed 1",
				Error:   res.Error.Error(),
			})
			return
		}

		if res := db.Delete(&oldImages, `product_id = ?`, EnterData.ProductID); res.Error != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"Error": "Product updation failed 2",
				"E":     res.Error.Error(),
			})
			return
		}

		for _, image := range images {
			newuuid := uuid.New()
			newname := fmt.Sprint(newuuid)
			splitedname := strings.Split(image.Filename, ".")
			ext := splitedname[len(splitedname)-1]
			ctx.SaveUploadedFile(image, "/home/arifu/Desktop/ecom/public/images/"+newname+"."+ext)

			newImage := models.ProductImages{
				ProductID: EnterData.ProductID,
				ImageName: newname + "." + ext,
				IsDefault: true,
			}
			db.Create(&newImage)
		}

		ctx.AbortWithStatusJSON(http.StatusOK, models.Response{
			Status:  true,
			Message: "Product updation success",
			Error:   nil,
		})
	}
}

// finally end our Admin edit product controller

// -------------------------------------

// Delete old product controller. This function delete product etails into database
func AdminDeleteProductController() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		productid := ctx.Param("id")
		fmt.Println(productid)

		// Database model
		db := *config.GetDb()

		// Product and product images tables structure
		var Product models.Product

		if res := db.First(&Product, productid); res.Error != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
				Status:  false,
				Message: "Product not available",
				Error:   res.Error.Error(),
			})
			return
		}

		if Product.IsDeleted || !Product.IsActive {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
				Status:  false,
				Message: "Product not available",
				Error:   "Product deleted or non active",
			})
			return
		}

		Product.IsDeleted = true
		if res := db.Save(&Product); res.Error != nil {
			ctx.AbortWithStatusJSON(400, gin.H{
				"Error": "Product not deleted",
			})
			return
		}

		// Delete Success Message
		ctx.JSON(http.StatusOK, models.Response{
			Status:  true,
			Message: "Product deleted success",
			Error:   nil,
		})
	}
}

//
// Finally end Delete product section

// -----------------------------------------

// Get all products
type Response struct {
	ID      uint
	Create  time.Time
	Name    string
	Disc    string
	Price   float64
	Quntity uint
	Images  []models.ProductImages
}

func AdminGetAllProductController() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		page := ctx.Param("page")
		pageToInt, _ := strconv.Atoi(page)
		pageSize := 10

		db := *config.GetDb()
		var Products []models.Product
		var Images []models.ProductImages
		var result []Response

		// Create offset limit
		offset := (pageToInt - 1) * pageSize

		// This
		if res := db.Offset(offset).Limit(pageSize).Find(&Products); res.Error != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.Response{
				Status:  false,
				Message: "ailed",
				Error:   res.Error.Error(),
			})
		}

		for _, val := range Products {
			db.Joins("JOIN products ON products.id = product_images.product_id AND product_images.product_id = ?", val.ID).Find(&Images)
			newResp := Response{
				ID:      val.ID,
				Create:  val.CreatedAt,
				Name:    val.Name,
				Disc:    val.Disc,
				Price:   val.Price,
				Quntity: val.Quntity,
				Images:  Images,
			}

			// After I get all result. I append result on result array
			// This step mainly write for simplisity purpous
			result = append(result, newResp)
		}

		// This step is return success response
		ctx.JSON(200, gin.H{
			"status":    true,
			"page":      1,
			"pagelimit": pageSize,
			"Count":     len(result),
			"products":  result,
		})
	}
}

// Finally End Get all product controller
//

// ------------------------------------

// Get Product using product id
// This Controller Return a map[string]any model
// product name,disc,price and almost all details send on json model
func AdminGetProductUsingIDController() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		productid := ctx.Param("id")

		db := *config.GetDb()
		var Product models.Product
		var Images []models.ProductImages

		if res := db.First(&Product, productid); res.Error != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
				Status:  false,
				Message: "Product not found",
				Error:   res.Error.Error(),
			})
			return
		}

		fmt.Println(Product)

		if res := db.Joins("JOIN products ON products.id = product_images.product_id AND product_images.product_id = ?", Product.ID).Find(&Images); res.Error != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, models.Response{
				Status:  false,
				Message: "Product not found",
				Error:   res.Error.Error(),
			})
			return
		}

		ctx.JSON(200, gin.H{
			"Status": true,
			"Product": map[string]any{
				"ProductID": Product.ID,
				"Created":   Product.CreatedAt,
				"Name":      Product.Name,
				"Disc":      Product.Disc,
				"Price":     Product.Price,
				"Quntity":   Product.Quntity,
				"Images":    Images,
			},
			"Error": nil,
		})
	}
}
