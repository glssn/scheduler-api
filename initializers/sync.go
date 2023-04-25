package initializers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/glssn/scheduler-api/api/models"
)

type Response struct {
	EnglandAndWales UKDivision `json:"england-and-wales"`
	Scotland        UKDivision `json:"scotland"`
	NorthernIreland UKDivision `json:"northern-ireland"`
}

type UKDivision struct {
	Division string    `json:"division"`
	Holidays []Holiday `json:"events"`
}

type Holiday struct {
	Title   string `json:"title"`
	Date    string `json:"date"`
	Notes   string `json:"notes"`
	Bunting bool   `json:"bunting"`
}

type EventList struct {
	Events []models.Event
}

func (eventlist *EventList) AddEvent(event models.Event) {
	eventlist.Events = append(eventlist.Events, event)
}

func convertToEvent(holidays Response, bot models.User) []models.Event {
	// create a new logger instance
	logger := Logger()

	eventlist := EventList{}

	for _, hol := range holidays.EnglandAndWales.Holidays {
		date, err := time.Parse("2006-01-02", hol.Date)
		if err != nil {
			// log the error message
			logger.Println(err)
			continue
		}

		event := models.Event{
			Type:              "bank_holiday",
			Title:             hol.Title,
			StartDate:         date,
			AllDay:            true,
			User:              bot,
			RecurringType:     "None",
			RecurringInterval: 0,
		}
		eventlist.AddEvent(event)
	}
	return eventlist.Events
}

func retrieveBankHolidays() (bank_holidays Response, err error) {
	// https://www.gov.uk/bank-holidays.json
	// create a new logger instance
	logger := Logger()

	var holidays Response

	// Get the JSON data
	r, err := http.Get("https://www.gov.uk/bank-holidays.json")
	if err != nil {
		// log the error message
		logger.Println(err)
		return holidays, err
	}
	defer r.Body.Close()

	// Decode and unmarshal the response into structs
	err = json.NewDecoder(r.Body).Decode(&holidays)
	if err != nil {
		// log the error message
		logger.Println(err)
		return holidays, err
	}
	return holidays, nil
}

func PopulateBankHolidays() {
	// create a new logger instance
	logger := Logger()

	holidays, err := retrieveBankHolidays()
	if err != nil {
		// log the error message
		logger.Println(err)
		return
	}
	// log the number of bank holidays retrieved
	logger.Printf("Retrieved %d bank holidays from gov.uk/bank-holidays", len(holidays.EnglandAndWales.Holidays))

	// create bank holiday bot user
	bankHolidayBotUser := models.User{
		Username: "bank-holiday-bot",
		Role:     "bot",
	}
	// Create the bot user if it doesn't already exist
	DB.Where(&bankHolidayBotUser).FirstOrCreate(&bankHolidayBotUser)
	// convert the holidays into models.Event
	events := convertToEvent(holidays, bankHolidayBotUser)
	// log the number of bank holiday events added to the database
	logger.Printf("Adding %d bank holiday events to the database", len(events))
	// add events to database, currently sequentially
	for _, event := range events {
		// only create if there isn't a Type: 'bank_holiday' event on this StartDate
		DB.FirstOrCreate(&event, models.Event{Type: event.Type, StartDate: event.StartDate})
	}
	logger.Printf("Successfully synchronised %d bank holiday events with the database", len(events))
}

func SyncBankHolidays() {
	// Create a ticker that will fire every 6 hours.
	ticker := time.NewTicker(time.Hour * 6)

	// Loop forever and sync the db every time the ticker fires
	for {
		select {
		case <-ticker.C:
			Logger().Println("Synchronising database with new bank holidays...")
			PopulateBankHolidays()
		}
	}
}
