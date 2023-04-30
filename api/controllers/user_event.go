package controllers

import (
	"net/http"
	"time"

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
func GetEventByUserID(id *int, c *gin.Context) {
	var events []models.Event
	if err := initializers.DB.Where("user_id = ?", &id).Find(&events).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No events found."})
		return
	}
	// convert the event to apiEvent
	var apiEvents []APIEvent
	apiEvents = eventsToAPIEvents(events)

	// Return the event
	c.JSON(http.StatusOK, apiEvents)
}

func GetEventByUserIdAndType(id *int, t *string, c *gin.Context) {
	var events []models.Event
	if err := initializers.DB.Where("user_id = ? AND type = ?", &id, &t).Find(&events).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No events found."})
		return
	}
	// convert the event to apiEvent
	var apiEvents []APIEvent
	apiEvents = eventsToAPIEvents(events)

	// Return the event
	c.JSON(http.StatusOK, apiEvents)
}

func GetEventByUserIdAndTypeAndDate(id *int, t *string, date *time.Time, c *gin.Context) {
	var events []models.Event
	if err := initializers.DB.Where("user_id = ? AND type = ? AND start_date = ?", &id, &t, &date).Find(&events).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No events found."})
		return
	}
	// convert the event to apiEvent
	var apiEvents []APIEvent
	apiEvents = eventsToAPIEvents(events)

	// Return the event
	c.JSON(http.StatusOK, apiEvents)
}
func GetEventByUserIdAndTypeAndDateRange(id *int, t *string, startDate *time.Time, endDate *time.Time, c *gin.Context) {
	var events []models.Event
	if err := initializers.DB.Where("user_id = ? AND type = ? AND start_date BETWEEN ? AND ?", &id, &t, &startDate, &endDate).Find(&events).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No events found."})
		return
	}
	// convert the event to apiEvent
	var apiEvents []APIEvent
	apiEvents = eventsToAPIEvents(events)

	// Return the event
	c.JSON(http.StatusOK, apiEvents)
}
