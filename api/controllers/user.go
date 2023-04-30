package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/glssn/scheduler-api/api/models"
)

type APIUser struct {
	ID       float64 `json:"id" gorm:"foreignKey:UserID;references:ID"`
	Username string  `json:"username"`
	Role     string  `json:"role"`
}

// userToAPIUser converts a User struct to an APIUser struct.
// The APIUser struct is a subset of the User struct, containing only the fields that are needed by the API.
func userToAPIUser(user models.User) APIUser {
	apiUser := APIUser{}
	// Marshal the User struct into a JSON string
	jsonUser, err := json.Marshal(user)
	if err != nil {
		log.Println(err)
	}
	// Parse the JSON string into the apiUser struct
	err = json.Unmarshal(jsonUser, &apiUser)
	if err != nil {
		log.Println(err)
	}
	return apiUser
}

// GET /api/user/:id
// Get user by ID
func GetUserByID(c *gin.Context) {
	var user models.User

	if err := db.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

// GET /api/users
// Get all users
func GetAllUsers(c *gin.Context) {
	var users []models.User

	db.Find(&users)
	c.JSON(http.StatusOK, users)
}

// POST /api/user/:id
// Update a user by ID
func UpdateUserByID(c *gin.Context) {
	var user models.User

	if err := db.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}
