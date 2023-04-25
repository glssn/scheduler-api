package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/glssn/scheduler-api/api/models"
	"github.com/glssn/scheduler-api/initializers"
)

type NewEventInput struct {
	Type              string    `binding:"required"`
	StartDate         time.Time `binding:"required"`
	EndDate           time.Time
	AllDay            bool
	RecurringType     string
	RecurringInterval uint32
}

type UpdateEventInput struct {
	Type              string
	StartDate         time.Time
	EndDate           time.Time
	AllDay            bool
	RecurringType     string
	RecurringInterval uint32
}

type APIEvent struct {
	ID                float64   `json:"id"`
	Type              string    `json:"type"`
	Title             string    `json:"title"`
	StartDate         time.Time `json:"start_date"`
	EndDate           time.Time `json:"end_date"`
	AllDay            bool      `json:"all_day"`
	RecurringType     string    `json:"recurring_type"`
	RecurringInterval uint32    `json:"recurring_interval"`
	User              APIUser   `json:"user"`
	UserID            int       `json:"user_id"`
}

// ParseTime takes a string representing a date and time and attempts to parse it using a list of supported formats.
// If the string can be parsed successfully, the corresponding time.Time value is returned.
// If the string cannot be parsed, an error is returned.
func ParseTime(input string) (time.Time, error) {
	var formats = []string{"2006-01-02", "2006-01-02T15:04:05", "2006-01-02T15:04:05.000Z"}
	for _, format := range formats {
		t, err := time.Parse(format, input)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, errors.New("unrecognised time format")
}

// ConvertDateRange takes two date strings and returns the corresponding time.Time values.
// If either of the date strings cannot be parsed, an error is returned.
func ConvertDateRange(min string, max string) (time.Time, time.Time, error) {
	parsedMin, errMin := ParseTime(min)
	parsedMax, errMax := ParseTime(max)
	if errMin == nil && errMax == nil {
		return parsedMin, parsedMax, nil
	}
	return time.Time{}, time.Time{}, errors.New("could not convert date range")
}

// ParseRecurringInterval takes a string representing a recurring interval and returns the corresponding timestamp value.
// Supported intervals are "daily", "weekly", "fortnightly", "monthly", and "yearly".
// If the provided interval is not supported, an error is returned.
func ParseRecurringInterval(interval string) (uint32, error) {
	var intervals = map[string]uint32{
		"daily":       86400,
		"weekly":      604800,
		"fortnightly": 1.21e+6,
		"monthly":     2.628e+6,
		"yearly":      3.154e+7,
	}
	timestampInterval, ok := intervals[interval]
	if ok {
		return timestampInterval, nil
	} else {
		return 0, errors.New("could not parse interval")
	}
}

// GET /events/all
// Get all events
func FindEvents(c *gin.Context) {
	var events []APIEvent
	initializers.DB.Model(&models.Event{}).Find(&events)
	c.JSON(http.StatusOK, events)
}

// GET /events
// Get events based on the specified parameters
// If the query string is empty, fetch all events
// If the query string contains an "id" parameter, fetch the event with the specified id
// If the query string contains "start_date" and "end_date" parameters, fetch the events where the start_date and end_dates are within the specified range,
// or where the start_date is between the specified startDate and endDate and the all_day field is true
func GetEvent(c *gin.Context) {

	params := c.Request.URL.Query()

	// no parameters
	if len(params) == 0 {
		// query string is empty, so fetch all results
		var events []APIEvent
		if err := initializers.DB.Model(&models.Event{}).Find(&events).Error; err != nil {
			c.AbortWithStatus(404)
			fmt.Println(err)
		} else {
			c.JSON(http.StatusOK, events)
			return
		}
	}

	// id
	if len(params) == 1 && params.Get("id") != "" {
		// id parameter is present, so fetch the event with the specified id
		id := params.Get("id")
		var event APIEvent
		if err := initializers.DB.Model(&models.Event{}).Where("id = ?", id).First(&event).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
			return
		} else {
			c.JSON(http.StatusOK, event)
			return
		}
	}

	// type
	if len(params) == 1 && params.Get("type") != "" {
		// type parameter is present, so fetch all events with the specified type
		var events []APIEvent
		t := params.Get("type")
		if err := initializers.DB.Model(&models.Event{}).Where("type = ?", t).Find(&events).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "No events found for event type" + t})
			return
		} else {
			c.JSON(http.StatusOK, events)
			return
		}
	}

	// user_id
	if len(params) == 1 && params.Get("user_id") != "" {
		// user_id parameter is present, so fetch all events with the specified user_id
		var events []APIEvent
		userID := params.Get("user_id")
		if err := initializers.DB.Model(&models.Event{}).Where("user_id = ?", userID).Find(&events).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "No events found for user" + userID})
			return
		} else {
			c.JSON(http.StatusOK, events)
			return
		}
	}

	// start_date and end_date
	if len(params) == 2 && params.Get("start_date") != "" && params.Get("end_date") != "" {
		// Fetch the events where the start_date and end_dates are within the specified range,
		// or where the start_date is between the specified startDate and endDate and the end_date is NULL or the all_day field is true
		var events []APIEvent
		startDate := params.Get("start_date")
		endDate := params.Get("end_date")
		if err := initializers.DB.Model(&models.Event{}).Where("(start_date >= ? AND end_date <= ?) OR (start_date BETWEEN ? AND ? AND all_day = true)", startDate, endDate, startDate, endDate).Find(&events).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
			return
		} else {
			c.JSON(http.StatusOK, events)
			return
		}
	}

	// start_date and end_date and user_id
	if len(params) == 3 && params.Get("start_date") != "" && params.Get("end_date") != "" && params.Get("user_id") != "" {
		// handle request with the "start_date" and "end_date" and "user_id" parameters
		var events []APIEvent
		startDate := params.Get("start_date")
		endDate := params.Get("end_date")
		userID := params.Get("user_id")
		if err := initializers.DB.Model(&models.Event{}).Where("user_id = ? AND ((start_date >= ? AND end_date <= ?) OR (start_date BETWEEN ? AND ? AND all_day = true))", userID, startDate, endDate, startDate, endDate).Find(&events).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Events not found"})
			return
		} else {
			c.JSON(http.StatusOK, events)
			return
		}
	}

	// start_date and end_date and type
	if len(params) == 3 && params.Get("start_date") != "" && params.Get("end_date") != "" && params.Get("type") != "" {
		// handle request with the "start_date" and "end_date" and "type" parameters
		var events []APIEvent
		startDate := params.Get("start_date")
		endDate := params.Get("end_date")
		t := params.Get("type")
		if err := initializers.DB.Model(&models.Event{}).Where("type = ? AND ((start_date >= ? AND end_date <= ?) OR (start_date BETWEEN ? AND ? AND all_day = true))", t, startDate, endDate, startDate, endDate).Find(&events).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Events not found"})
			return
		} else {
			c.JSON(http.StatusOK, events)
			return
		}
	}

	// start_date and end_date and type and user_id
	if len(params) == 4 && params.Get("start_date") != "" && params.Get("end_date") != "" && params.Get("type") != "" && params.Get("user_id") != "" {
		// handle request with the "start_date" and "end_date" and "type" and "user_id" parameters
		var events []APIEvent
		startDate := params.Get("start_date")
		endDate := params.Get("end_date")
		t := params.Get("type")
		userID := params.Get("user_id")
		if err := initializers.DB.Model(&models.Event{}).Where("type = ? AND user_id = ? AND ((start_date >= ? AND end_date <= ?) OR (start_date BETWEEN ? AND ? AND all_day = true))", t, userID, startDate, endDate, startDate, endDate).Find(&events).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Events not found"})
			return
		} else {
			c.JSON(http.StatusOK, events)
			return
		}
	}

	c.JSON(http.StatusBadRequest, nil)
	return

}

// POST /events
// Create a new event with the provided details
// The request must include a valid JSON object with the event details
// If the input is invalid, return a 400 status code
// If the user is not authenticated, return a 401 status code
// Otherwise, return the created event object and a 200 status code
func CreateEvent(c *gin.Context) {
	// Validate input
	var input NewEventInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Get user ID from middleware
	userFromCookie, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{})
		return
	}
	user := userFromCookie.(models.User)

	// Create event
	event := models.Event{
		Type:              input.Type,
		StartDate:         input.StartDate,
		EndDate:           input.EndDate,
		AllDay:            input.AllDay,
		RecurringType:     input.RecurringType,
		RecurringInterval: input.RecurringInterval,
		User:              user,
	}
	initializers.DB.Save(&event)
	c.JSON(http.StatusOK, gin.H{"data": event})

}

// PATCH /events/:id
// Update an existing event with the specified id
// The request must include a valid JSON object with the updated event details
// If the id is not provided, return a 400 status code
// If the event with the specified id does not exist, return a 404 status code
// If the input is invalid, return a 400 status code
// Otherwise, return the updated event object and a 200 status code
func UpdateEvent(c *gin.Context) {
	// Get the id parameter from the request
	id := c.Param("id")

	// Return a 400 status code if the id is not provided
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request. ID is required."})
		return
	}

	// Check if event exists
	var event []models.Event
	if err := initializers.DB.Where("id = ?", id).First(&event).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found."})
		return
	}
	// Validate input
	var input UpdateEventInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newevent := models.Event{
		Type:              input.Type,
		StartDate:         input.StartDate,
		EndDate:           input.EndDate,
		AllDay:            input.AllDay,
		RecurringType:     input.RecurringType,
		RecurringInterval: input.RecurringInterval,
	}
	initializers.DB.Model(&event).Where("id = ?", id).Updates(newevent)
	c.JSON(http.StatusOK, gin.H{"data": newevent})
}

// DELETE /events/:id
// Delete a event
func DeleteEvent(c *gin.Context) {
	// Get the id parameter from the request
	id := c.Param("id")

	// Return a 400 status code if the id is not provided
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request. ID is required."})
		return
	}

	// Get model if exist
	var event models.Event
	if err := initializers.DB.Where("id = ?", id).First(&event).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Event not found."})
		return
	}

	initializers.DB.Delete(&event)

	c.JSON(http.StatusOK, gin.H{"data": true})
}
