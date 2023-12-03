package controllers

import (
	"astora/inits"
	"astora/models"

	"github.com/gin-gonic/gin"
)

// sercice controller

func CreateService(ctx *gin.Context) {
	var body struct {
		Name        string
		Description string
		Price       int
		ServiceURL  string
	}

	ctx.BindJSON(&body)

	// Make sure that the service does not exist
	var service models.Services
	result := inits.DB.Where("name = ?", body.Name).First(&service)
	if result.Error == nil {
		ctx.JSON(400, gin.H{
			"error": "Service already exists",
		})
		return
	}

	service = models.Services{
		Name:        body.Name,
		Description: body.Description,
		Price:       body.Price,
		ServiceURL:  body.ServiceURL,
	}

	result = inits.DB.Create(&service)
	if result.Error != nil {
		ctx.JSON(400, gin.H{
			"error": result.Error,
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "Service created successfully",
		"data":    service,
	})
}

func GetServices(ctx *gin.Context) {
	var services []models.Services
	result := inits.DB.Find(&services)

	if result.Error != nil {
		ctx.JSON(400, gin.H{
			"error": result.Error,
		})
		return
	}

	ctx.JSON(200, gin.H{
		"data": services,
	})
}

func GetService(ctx *gin.Context) {
	var service models.Services
	serviceID := ctx.Param("id")
	result := inits.DB.Where("id = ?", serviceID).First(&service)

	if result.Error != nil {
		ctx.JSON(400, gin.H{
			"error": result.Error,
		})
		return
	}

	ctx.JSON(200, gin.H{
		"data": service,
	})
}

func UpdateService(ctx *gin.Context) {
	var body struct {
		Name        string
		Description string
		Price       int
		ServiceURL  string
	}

	ctx.BindJSON(&body)

	var service models.Services
	serviceID := ctx.Param("id")
	result := inits.DB.Where("id = ?", serviceID).First(&service)

	if result.Error != nil {
		ctx.JSON(400, gin.H{
			"error": "Service does not exist",
		})
		return
	}

	service.Name = body.Name
	service.Description = body.Description
	service.Price = body.Price
	service.ServiceURL = body.ServiceURL

	result = inits.DB.Save(&service)
	if result.Error != nil {
		ctx.JSON(400, gin.H{
			"error": result.Error,
		})
		return
	}

	ctx.JSON(200, gin.H{
		"data": service,
	})
}

func DeleteService(ctx *gin.Context) {
	var service models.Services
	serviceID := ctx.Param("id")
	result := inits.DB.Where("id = ?", serviceID).First(&service)

	if result.Error != nil {
		ctx.JSON(400, gin.H{
			"error": "Service does not exist",
		})
		return
	}

	result = inits.DB.Delete(&service)
	if result.Error != nil {
		ctx.JSON(400, gin.H{
			"error": result.Error,
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "Service deleted successfully",
	})
}
