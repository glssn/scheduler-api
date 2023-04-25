package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/glssn/scheduler-api/api/models"
	"github.com/glssn/scheduler-api/initializers"
)

type APIUser struct {
	ID       float64 `json:"id" gorm:"foreignKey:UserID;references:ID"`
	Username string  `json:"username"`
	Role     string  `json:"role"`
}

func UserRoutes(incomingRoutes *gin.Engine) {
}

// GET /api/user/:id
// Get user by ID
func GetUserByID(c *gin.Context) {
	var user models.User

	if err := initializers.DB.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

// GET /api/users
// Get all users
func GetAllUsers(c *gin.Context) {
	var users []models.User

	initializers.DB.Find(&users)
	c.JSON(http.StatusOK, users)
}
