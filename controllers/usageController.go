package controllers

import (
	"astora/inits"
	"astora/models"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Use services controller to make a request to the service URL and then return the response
func UseService(ctx *gin.Context) {
	var body struct {
		ServiceName   string
		Params        string
		RequestMethod string
	}

	ctx.BindJSON(&body)

	// Make sure that the service exists
	var service models.Services
	result := inits.DB.Where("name = ?", body.ServiceName).First(&service)
	if result.Error != nil {
		ctx.JSON(400, gin.H{
			"error": "Service does not exist",
		})
		return
	}

	// Make sure that the user exists
	var user models.User
	result = inits.DB.Where("id = ?", ctx.Keys["id"]).First(&user)
	if result.Error != nil {
		ctx.JSON(400, gin.H{
			"error": "User does not exist",
		})
		return
	}

	// Make sure that the user has enough tokens
	if user.TokenBalance < service.Price {
		ctx.JSON(400, gin.H{
			"error": "Not enough tokens",
		})
		return
	}

	// Initialize an empty Params JSON if not provided
	var paramsJSON []byte
	if body.Params != "" {
		// Convert body.Params to JSON string
		var err error
		paramsJSON, err = json.Marshal(body.Params)
		if err != nil {
			ctx.JSON(400, gin.H{
				"error": "Error converting params to JSON",
			})
			return
		}
	} else {
		paramsJSON = []byte("{}") // Default to an empty JSON object
	}

	// Create a new request with the specified RequestMethod, service URL, and params in the body
	req, err := http.NewRequest(body.RequestMethod, service.ServiceURL, bytes.NewBuffer(paramsJSON))
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": "Error creating request",
		})
		return
	}

	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": "Error making request",
		})
		return
	}
	defer resp.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": "Error reading response body",
		})
		return
	}

	// Create a new ServiceLog entry
	serviceLog := models.ServiceLogs{
		UserID:    user.ID,
		ServiceID: service.ID,
		Log:       string(responseBody),
	}
	inits.DB.Create(&serviceLog)

	// Create a new ApiResponse entry
	apiResponse := models.ApiResponses{
		RequestHeaders:    "", // Add request headers if needed
		RequestBody:       string(paramsJSON),
		RequestTime:       time.Now().Format(time.RFC3339),
		RequestStatusCode: resp.Status,
		RequestUrl:        service.ServiceURL,
		RequestMethod:     body.RequestMethod,
		RequestResponse:   string(responseBody),
		RequestPrice:      service.Price,
		UserID:            user.ID,
	}
	inits.DB.Create(&apiResponse)

	// Update the user token balance
	user.TokenBalance = user.TokenBalance - service.Price
	result = inits.DB.Save(&user)
	if result.Error != nil {
		ctx.JSON(400, gin.H{
			"error": "Error updating user",
		})
		return
	}

	// Return the ApiResponse data
	ctx.JSON(200, gin.H{
		"message": "Service used successfully",
		"data":    apiResponse,
	})
}

// Get all the api responses for a user
func GetApiResponses(ctx *gin.Context) {
	var apiResponses []models.ApiResponses
	result := inits.DB.Where("user_id = ?", ctx.Keys["id"]).Find(&apiResponses)
	if result.Error != nil {
		ctx.JSON(400, gin.H{
			"error": "Error getting api responses",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "Api responses fetched successfully",
		"data":    apiResponses,
	})
}
