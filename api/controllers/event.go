package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/glssn/scheduler-api/api/models"
	"github.com/glssn/scheduler-api/initializers"
)

type NewEventInput struct {
	Type              string    `binding:"required" json:"type"`
	Title             string    `json:"title"`
	StartDate         time.Time `binding:"required" json:"start_date"`
	EndDate           time.Time `json:"end_date"`
	AllDay            bool      `json:"all_day"`
	RecurringType     string    `json:"recurring_type"`
	RecurringInterval uint32    `json:"recurring_interval"`
}

type PatchEventInput struct {
	Type              string    `binding:"required" json:"type"`
	Title             string    `json:"title"`
	StartDate         time.Time `binding:"required" json:"start_date"`
	EndDate           time.Time `json:"end_date"`
	AllDay            bool      `json:"all_day"`
	RecurringType     string    `json:"recurring_type"`
	RecurringInterval uint32    `json:"recurring_interval"`
}

type APIEvent struct {
	ID                float64   `form:"id" json:"id"`
	Type              string    `form:"type" json:"type"`
	Title             string    `json:"title"`
	StartDate         time.Time `json:"start_date"`
	EndDate           time.Time `json:"end_date"`
	AllDay            bool      `json:"all_day"`
	RecurringType     string    `json:"recurring_type"`
	RecurringInterval uint32    `json:"recurring_interval"`
	// User              APIUser   `json:"user"`
	UserID int `json:"user_id"`
}

// eventToAPIEvent converts a Event struct to an APIEvent struct.
// The APIEvent struct is a subset of the Event struct, containing only the fields that are needed by the API.
func eventToAPIEvent(event models.Event) (APIEvent, error) {
	if event.Type == "" {
		log.Println("found a null Event")
		return APIEvent{}, errors.New("event is null")
	}
	apiEvent := APIEvent{}
	// Marshal the Event struct into a JSON string
	jsonEvent, err := json.Marshal(event)
	if err != nil {
		log.Println(err)
		return apiEvent, errors.New("Could not marshal event")
	}
	// Parse the JSON string into the apiEvent struct
	err = json.Unmarshal(jsonEvent, &apiEvent)
	if err != nil {
		log.Println(err)
		return apiEvent, errors.New("Could not unmarshal into APIEvent")
	}
	return apiEvent, nil
}

// eventsToAPIEvents converts a slice of Events struct to a slice of APIEvents
// The APIEvent struct is a subset of the Event struct, containing only the fields that are needed by the API.
func eventsToAPIEvents(events []models.Event) []APIEvent {
	apiEvents := make([]APIEvent, 0)
	for _, event := range events {
		if event.Type == "" {
			log.Println("found a null Event type")
		}
		apiEvent, err := eventToAPIEvent(event)
		if err != nil {
			log.Println("error converting event to APIEvent:", err)
			continue
		}
		apiEvents = append(apiEvents, apiEvent)
	}
	return apiEvents
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

type EventQuery struct {
	ID        float64   `form:"id"`
	Type      string    `form:"type"`
	Date      time.Time `form:"date"`
	StartDate time.Time `form:"start_date"`
	EndDate   time.Time `form:"end_date"`
	UserID    int       `form:"user_id"`
}

// GET /events
// Get events based on the specified parameters
// If the query string is empty, fetch all events
// If the query string contains an "id" parameter, fetch the event with the specified id
// If the query string contains "start_date" and "end_date" parameters, fetch the events where the start_date and end_dates are within the specified range,
// or where the start_date is between the specified startDate and endDate and the all_day field is true
func GetEvent(c *gin.Context) {
	// Get the query parameters
	// params := c.Request.URL.Query()

	var eventQuery EventQuery
	if err := c.ShouldBind(&eventQuery); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters for eventQuery"})
		return
	}
	log.Println("====== Only Bind By Query String ======")
	log.Println("id: ", eventQuery.ID)
	log.Println("type: ", eventQuery.Type)
	log.Println("date: ", eventQuery.Date)
	log.Println("time is zero?: ", eventQuery.Date.IsZero())
	log.Println("start_date: ", eventQuery.StartDate)
	log.Println("end_date: ", eventQuery.EndDate)
	log.Println("user_id: ", eventQuery.UserID)

	// If event ID is provided, retrieve it
	if eventQuery.ID != float64(0) {
		GetEventById(&eventQuery.ID, c)
		return
	}
	// If type is provided, but no User ID
	if eventQuery.Type != "" && eventQuery.UserID == 0 {
		// if no date or date range is provided
		if eventQuery.Date.IsZero() && eventQuery.StartDate.IsZero() && eventQuery.StartDate.IsZero() {
			GetEventByType(&eventQuery.Type, c)
			return
		}
		// if date is provided
		if !eventQuery.Date.IsZero() {
			GetEventByTypeAndDate(&eventQuery.Type, &eventQuery.Date, c)
			return
		}
		// if date range is provided
		if !eventQuery.StartDate.IsZero() && !eventQuery.StartDate.IsZero() {
			GetEventByTypeAndDateRange(&eventQuery.Type, &eventQuery.StartDate, &eventQuery.EndDate, c)
			return
		}
	}
	// If User ID is provided
	if eventQuery.UserID != 0 {
		// if no type or date or date range is provided
		if eventQuery.Type == "" && eventQuery.Date.IsZero() && eventQuery.StartDate.IsZero() && eventQuery.StartDate.IsZero() {
			GetEventByUserID(&eventQuery.UserID, c)
			return
		}
		// if type is provided, but date and date range is not provided
		if eventQuery.Type != "" && eventQuery.Date.IsZero() && eventQuery.StartDate.IsZero() && eventQuery.StartDate.IsZero() {
			GetEventByUserIdAndType(&eventQuery.UserID, &eventQuery.Type, c)
			return
		}
		// if type and date is provided but not date range
		if eventQuery.Type != "" && !eventQuery.Date.IsZero() && eventQuery.StartDate.IsZero() && eventQuery.StartDate.IsZero() {
			GetEventByUserIdAndTypeAndDate(&eventQuery.UserID, &eventQuery.Type, &eventQuery.Date, c)
			return
		}
		// if type and date range is provided
		if eventQuery.Type != "" && !eventQuery.StartDate.IsZero() && !eventQuery.StartDate.IsZero() {
			GetEventByUserIdAndTypeAndDateRange(&eventQuery.UserID, &eventQuery.Type, &eventQuery.StartDate, &eventQuery.EndDate, c)
			return
		}
	}
	// If date is provided
	if !eventQuery.Date.IsZero() {
		GetEventByDate(&eventQuery.Date, c)
		return
	}
	// If date range is provided
	if !eventQuery.StartDate.IsZero() && !eventQuery.StartDate.IsZero() {
		GetEventByDateRange(&eventQuery.StartDate, &eventQuery.EndDate, c)
		return
	}
	FindEvents(c)
	c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request."})
}

// GetEventById returns the event from the database where the ID is
func GetEventById(id *float64, c *gin.Context) {
	var event models.Event
	if err := initializers.DB.Where("id = ?", *id).First(&event).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found."})
		return
	}
	// convert the event to apiEvent
	var apiEvent APIEvent
	apiEvent, err := eventToAPIEvent(event)
	if err != nil {
		log.Println("error converting event to APIEvent:", err)
	}

	// Return the event
	c.JSON(http.StatusOK, apiEvent)
	return
}

func GetEventByType(t *string, c *gin.Context) {
	var events []models.Event
	if err := initializers.DB.Where("type = ?", &t).Find(&events).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No events found."})
		return
	}
	// convert the event to apiEvent
	var apiEvents []APIEvent
	apiEvents = eventsToAPIEvents(events)

	// Return the event
	c.JSON(http.StatusOK, apiEvents)
}

func GetEventByTypeAndDate(t *string, date *time.Time, c *gin.Context) {
	var events []models.Event
	if err := initializers.DB.Where("type = ? & start_date = ?", &t, &date).Find(&events).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No events found."})
		return
	}
	// convert the event to apiEvent
	var apiEvents []APIEvent
	apiEvents = eventsToAPIEvents(events)

	// Return the event
	c.JSON(http.StatusOK, apiEvents)
}

func GetEventByTypeAndDateRange(t *string, startDate *time.Time, endDate *time.Time, c *gin.Context) {
	var events []models.Event
	if err := initializers.DB.Where("(type = ? AND start_date >= ? AND end_date <= ?) OR (type = ? AND start_date BETWEEN ? AND ? AND all_day = true)", &t, &startDate, &endDate, &t, &startDate, &endDate).Find(&events).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No events found."})
		return
	}
	// convert the event to apiEvent
	var apiEvents []APIEvent
	apiEvents = eventsToAPIEvents(events)

	// Return the event
	c.JSON(http.StatusOK, apiEvents)
}

func GetEventByDate(date *time.Time, c *gin.Context) {
	var events []models.Event
	if err := initializers.DB.Where("start_date = ?", &date).Find(&events).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No events found."})
		return
	}
	// convert the event to apiEvent
	var apiEvents []APIEvent
	apiEvents = eventsToAPIEvents(events)

	// Return the event
	c.JSON(http.StatusOK, apiEvents)
}

func GetEventByDateRange(startDate *time.Time, endDate *time.Time, c *gin.Context) {
	var events []models.Event
	if err := initializers.DB.Where("(start_date >= ? AND end_date <= ?) OR (start_date BETWEEN ? AND ? AND all_day = true)", &startDate, &endDate, &startDate, &endDate).Find(&events).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No events found."})
		return
	}
	// convert the event to apiEvent
	var apiEvents []APIEvent
	apiEvents = eventsToAPIEvents(events)

	// Return the event
	c.JSON(http.StatusOK, apiEvents)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}
	// Get user ID from middleware
	userFromCookie, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
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
	apiEvent, err := eventToAPIEvent(event)
	if err != nil {
		log.Println("error converting event to APIEvent:", err)
	}
	c.JSON(http.StatusCreated, apiEvent)
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
	var event models.Event
	if err := initializers.DB.Where("id = ?", id).First(&event).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found."})
		return
	}
	// Validate input
	var apiEvent APIEvent
	if err := c.ShouldBindJSON(&apiEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the event in the database
	if err := initializers.DB.Model(&event).Where("id = ?", id).Updates(&apiEvent).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	apiEvent, err := eventToAPIEvent(event)
	if err != nil {
		log.Println("error converting event to APIEvent:", err)
	}

	c.JSON(http.StatusOK, &apiEvent)
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
