package main

import (
	"astora/controllers"
	"astora/inits"
	"astora/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	// LoadEnv()
	inits.LoadEnv()
	inits.DBInit()
}

func main() {
	router := gin.Default()
	config := cors.Config{
		AllowOrigins:     []string{"*"}, // Allow any origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}
	router.Use(cors.New(config))
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})
	router.POST("/register", controllers.CreateUser)
	router.POST("/login", controllers.LoginUser)
	router.GET("/dashboard", middlewares.AuthMiddleware(), controllers.Dashboard)

	// Service routes. group them under /service
	serviceRoutes := router.Group("/service")
	serviceRoutes.Use(middlewares.AuthMiddleware())
	{
		serviceRoutes.POST("/create", controllers.CreateService)
		serviceRoutes.GET("/list", controllers.GetServices)
		serviceRoutes.POST("/use", controllers.UseService)
		serviceRoutes.GET("/get", controllers.GetService)
		serviceRoutes.POST("/update", controllers.UpdateService)
		serviceRoutes.POST("/delete", controllers.DeleteService)
	}

	// Usage routes. group them under /usage
	usageRoutes := router.Group("/usage")
	usageRoutes.Use(middlewares.AuthMiddleware())
	{
		usageRoutes.POST("/use", controllers.UseService)
		usageRoutes.GET("/used", controllers.GetApiResponses)
	}
	router.Run(":8080")
}
