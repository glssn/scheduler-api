package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/glssn/scheduler-api/api/models"
	"github.com/glssn/scheduler-api/initializers"
)

// GET /api/events/user
// Retrieves events for the specified user ID and type (if provided).
// If the "type" parameter is not provided, the function returns all events for the specified user ID.
// If the "type" parameter is provided, the function returns events of the specified type for the specified user ID.
// If no events are found for the specified user ID and type, the function returns a 404 Not Found response.
// Otherwise, the function returns a 200 OK response with the found events.
func GetUserEventByUserID(c *gin.Context) {
	// Return a 400 status code if the id or query is not provided
	if c.Query("id") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request. ID is required."})
		return
	}

	var events []models.Event

	// Set the query to search for all events if the "type" parameter is not provided
	query := initializers.DB.Where("user_id = ?", c.Query("id"))
	if c.Query("type") != "" {
		// handle "type" parameter if present
		switch c.Query("type") {
		case "duty_tech":
			query = query.Where("type IN ?", []string{"DutyTech1", "DutyTech2"})
		default:
			query = query.Where("type = ?", c.Query("type"))
		}
	}

	// Find the events
	if err := query.Find(&events).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Events not found."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"events": events})
}
