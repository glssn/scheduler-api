package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"
	"github.com/glssn/scheduler-api/api/controllers"
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func TestGetEvent(t *testing.T) {
	mockResponse := `{"data":[{"ID":7,"CreatedAt":"2022-08-21T22:05:43.723079+01:00","UpdatedAt":"2022-08-21T22:05:43.723079+01:00","DeletedAt":null,"Type":"DutyTech1","Title":"","StartDate":"2050-04-23T00:00:00Z","EndDate":"0001-01-01T00:00:00Z","AllDay":false,"RecurringType":"","RecurringInterval":0,"User":{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"Username":"","Role":""},"UserID":0}]}`
	r := SetUpRouter()
	r.GET("/api/events/:id", controllers.GetEvent)
	req, _ := http.NewRequest("GET", "/api/events/7", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	responseData, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, mockResponse, string(responseData))
	assert.Equal(t, http.StatusOK, w.Code)
}
