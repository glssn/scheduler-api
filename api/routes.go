package api

import (
	"github.com/gin-gonic/gin"
	"github.com/glssn/scheduler-api/api/controllers"
	"github.com/glssn/scheduler-api/api/middleware"
)

func Routes(app *gin.Engine) {
	// Event endpoints
	events := app.Group("/api/events")
	events.Use(middleware.RequireAuth)
	events.GET("/all", controllers.FindEvents)
	events.GET("/", controllers.GetEvent)
	events.GET("/:id", controllers.GetEvent)
	events.GET("?start_date=:end_date", controllers.GetEvent)
	events.POST("/", controllers.CreateEvent)
	events.PATCH("/:id", controllers.UpdateEvent) // broken
	events.DELETE("/:id", controllers.DeleteEvent)

	// Auth endpoints
	auth := app.Group("/")
	auth.POST("/login", controllers.Login)
	auth.GET("/validate", controllers.Validate)
	auth.POST("/logout", controllers.Logout)

	// User endpoints
	users := app.Group("/api/users")
	users.Use(middleware.RequireAuth)
	users.GET("/all", controllers.GetAllUsers)
	users.GET("/:id", controllers.GetUserByID)

	// User/event endpoints
	userevents := app.Group("/api/events/user")
	userevents.Use(middleware.RequireAuth)
	userevents.GET("/", controllers.GetUserEventByUserID)
}
