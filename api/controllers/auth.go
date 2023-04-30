package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/glssn/scheduler-api/api/models"
	"github.com/glssn/scheduler-api/initializers"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nerney/dappy"
)

var db = initializers.DB

// Login logs the user in and returns a JWT token.
func Login(c *gin.Context) {
	// get the user and pass from the request body
	var body struct {
		User     string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to connect to LDAP backend",
		})
	}

	if err := client.Auth(body.User, body.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid username and password",
		})
		return
	}

	// Create the user in the database
	user := models.User{Username: body.User, Role: "Viewer"}
	result := db.FirstOrCreate(&user, models.User{Username: body.User})

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})
	}

	apiUser := userToAPIUser(user)

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
	c.JSON(http.StatusOK, apiUser)
}

// Validate validates the user and returns an APIUser struct.
// If the user is not valid, it returns a 401 Unauthorized status code.
func Validate(c *gin.Context) {
	// get the user from the middleware
	user, _ := c.Get("user")

	if _, ok := user.(models.User); !ok {
		// Handle the case where the user is not of type models.User
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Create the return struct
	apiUser := userToAPIUser(user.(models.User))

	// Return the user in the response
	c.JSON(http.StatusOK, apiUser)
}

func Logout(c *gin.Context) {
	c.SetCookie("Authorization", "", -1, "", "", false, true)
	c.Status(http.StatusOK)
}
