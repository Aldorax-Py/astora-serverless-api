package controllers

import (
	"astora/inits"
	"astora/models"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// UserController is the controller for the user model
func CreateUser(ctx *gin.Context) {
	var body struct {
		Name        string
		CompanyName string
		Password    string
		Email       string
		Phone       string
	}

	ctx.BindJSON(&body)

	// Make sure that the user does not exist
	var user models.User
	result := inits.DB.Where("email = ?", body.Email).First(&user)
	if result.Error == nil {
		ctx.JSON(400, gin.H{
			"error": "User already exists",
		})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})

		return
	}

	user = models.User{
		Name:         body.Name,
		CompanyName:  body.CompanyName,
		Password:     string(hashedPassword),
		Email:        body.Email,
		Phone:        body.Phone,
		TokenBalance: 100000,
	}

	result = inits.DB.Create(&user)
	if result.Error != nil {
		ctx.JSON(400, gin.H{
			"error": result.Error,
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "User created successfully",
		"data":    user,
	})

}

func LoginUser(ctx *gin.Context) {
	var body struct {
		Email    string
		Password string
	}

	ctx.BindJSON(&body)

	user := models.User{
		Email: body.Email,
	}

	result := inits.DB.Where("email = ?", user.Email).First(&user)
	if result.Error != nil {
		ctx.JSON(400, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// Compare hashed password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// Create a token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
	})

	// Sign the token with the secret key from the .env file using godotenv
	secretKey := os.Getenv("SECRET")
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	// Set the token in the cookie
	ctx.SetCookie("token", tokenString, 3600, "/", "localhost", false, true)

	ctx.JSON(200, gin.H{
		"message": "User logged in successfully",
	})
}

func Dashboard(c *gin.Context) {
	userID, exists := c.Get("id")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	// Preload both the api responses and the serviceLogs
	result := inits.DB.Where("id = ?", userID).Preload("ApiResponses").Preload("ServiceLogs").First(&user)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"error": result.Error,
		})
		return
	}

	// Return the dashboard page with the user data
	c.JSON(200, gin.H{
		"message": "Welcome to your dashboard",
		"data":    user,
	})
}
