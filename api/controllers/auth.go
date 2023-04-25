package controllers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/glssn/scheduler-api/api/models"
	"github.com/glssn/scheduler-api/initializers"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nerney/dappy"
)

func Login(c *gin.Context) {
	// get the user and pass from the request body
	var body struct {
		User     string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body",
		})
		return
	}

	// Attempt to bind the requested user to LDAP
	client, err := dappy.New(dappy.Config{
		BaseDN: "dc=example,dc=com",
		Filter: "uid",
		ROUser: dappy.User{Name: "cn=read-only-admin,dc=example,dc=com", Pass: "password"},
		Host:   "ldap.forumsys.com:389",
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to connect to LDAP backend",
		})
	}

	if err := client.Auth(body.User, body.Password); err != nil {
		fmt.Println(body.User, body.Password)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid username and password",
		})
		return
	}

	// Create the user in the database
	user := models.User{Username: body.User, Role: "Viewer"}
	result := initializers.DB.FirstOrCreate(&user, models.User{Username: body.User})

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})
	}

	// Generate the JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 90).Unix(),
	})

	// Sign the token
	tokenStringSigned, err := token.SignedString([]byte(os.Getenv("JWT_AUTH_SECRET_KEY")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
		return
	}

	// Send the token back
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenStringSigned, 3600*24*90, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func Validate(c *gin.Context) {
	// get the user from the middleware
	user, _ := c.Get("user")

	if user == nil {
		// Handle the case where the user is nil
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Create the return struct
	type userReturn struct {
		ID       uint
		Role     string
		Username string
	}
	// Map the user to the return struct
	u := &userReturn{
		ID:       user.(models.User).ID,
		Role:     user.(models.User).Role,
		Username: user.(models.User).Username,
	}
	// Return the user in the response
	c.JSON(http.StatusOK, gin.H{
		"user": u,
	})
}

func Logout(c *gin.Context) {
	c.SetCookie("Authorization", "", -1, "", "", false, true)
	c.Status(http.StatusOK)
}
