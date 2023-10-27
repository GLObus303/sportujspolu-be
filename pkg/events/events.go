package events

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Service struct {
	db *sql.DB
}

func NewEventsService(db *sql.DB) *Service {
	return &Service{db}
}

type Event struct {
	ID          int64     `json:"id" example:"24"`
	Name        string    `json:"name" example:"Basketball Match at Park"`
	Sport       string    `json:"sport" example:"Basketball"`
	Date        time.Time `json:"date" example:"2023-07-10"`
	Location    string    `json:"location" example:"Central Park"`
	Price       float64   `json:"price" example:"0.00"`
	Description string    `json:"description" example:"Example Description"`
	Level       string    `json:"level" example:"Any"`
}

// @Summary Get all events
// @Description Retrieve all events from the database
// @Tags events
// @Produce json
// @Success 200 {array} Event
// @Router /events [get]
func (s *Service) GetEvents(c *gin.Context) {
	res, err := s.db.Query("SELECT * FROM events")
	if err != nil {
		log.Fatal("(GetEvents) db.Query", err)
	}

	defer res.Close()

	events := []Event{}
	for res.Next() {
		var event Event
		err := res.Scan(&event.ID, &event.Name, &event.Sport, &event.Date, &event.Location, &event.Price, &event.Description, &event.Level)
		if err != nil {
			log.Fatal("(GetEvents) res.Scan", err)
		}
		events = append(events, event)
	}

	c.JSON(http.StatusOK, events)
}

// @Summary Get a single event
// @Description Retrieves a single event from the database
// @Tags events
// @Produce json
// @Param eventId path int true "Event ID"
// @Success 200 {object} Event
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /events/{eventId} [get]
func (s *Service) GetSingleEvent(c *gin.Context) {
	eventId := c.Param("eventId")
	eventId = strings.ReplaceAll(eventId, "/", "")
	eventIdInt, err := strconv.Atoi(eventId)
	if err != nil {
		log.Fatal("(GetSingleEvent) strconv.Atoi", err)
	}

	var event Event
	err = s.db.QueryRow(`SELECT * FROM events WHERE id = ?`, eventIdInt).Scan(&event.ID, &event.Name, &event.Sport, &event.Date, &event.Location, &event.Price, &event.Description, &event.Level)
	if err != nil {
		log.Fatal("(GetSingleEvent) db.Exec", err)
	}

	c.JSON(http.StatusOK, event)
}

// @Summary Create a new event
// @Description Creates a new event in the database
// @Tags events
// @Accept json
// @Produce json
// @Param newEvent body Event true "Event object"
// @Success 200 {object} Event
// @Failure 400 {object} string
// @Failure 401 {object} string "Unauthorized"
// @Failure 500 {object} string
// @Security BearerAuth
// @Router /events [post]
func (s *Service) CreateEvent(c *gin.Context) {
	var newEvent Event
	err := c.BindJSON(&newEvent)
	if err != nil {
		log.Fatal("(CreateEvent) c.BindJSON", err)
	}

	query := `INSERT INTO events (name, sport) VALUES (?, ?)`
	res, err := s.db.Exec(query, newEvent.Name, newEvent.Sport)
	if err != nil {
		log.Fatal("(CreateEvent) db.Exec", err)
	}
	newEvent.ID, err = res.LastInsertId()
	if err != nil {
		log.Fatal("(CreateEvent) res.LastInsertId", err)
	}

	c.JSON(http.StatusOK, newEvent)
}

// @Summary Update an event
// @Description Update an existing event with the given event ID
// @Tags events
// @Accept json
// @Produce json
// @Param eventId path int true "Event ID"
// @Param event body Event true "Event object that needs to be updated"
// @Success 200
// @Failure 400 {object} string
// @Failure 401 {object} string "Unauthorized"
// @Failure 500 {object} string
// @Security BearerAuth
// @Router /events/{eventId} [put]
func (s *Service) UpdateEvent(c *gin.Context) {
	var updates Event
	err := c.BindJSON(&updates)
	if err != nil {
		log.Fatal("(UpdateEvent) c.BindJSON", err)
	}

	eventId := c.Param("eventId")
	eventId = strings.ReplaceAll(eventId, "/", "")
	eventIdInt, err := strconv.Atoi(eventId)
	if err != nil {
		log.Fatal("(UpdateEvent) strconv.Atoi", err)
	}

	query := `UPDATE events SET name = ?, sport = ? WHERE id = ?`
	_, err = s.db.Exec(query, updates.Name, updates.Sport, eventIdInt)
	if err != nil {
		log.Fatal("(UpdateEvent) db.Exec", err)
	}

	c.Status(http.StatusOK)
}

// @Summary Delete an event
// @Description Delete an existing event with the given event ID
// @Tags events
// @Param eventId path int true "Event ID"
// @Success 200
// @Failure 400 {object} string
// @Failure 401 {object} string "Unauthorized"
// @Failure 500 {object} string
// @Security BearerAuth
// @Router /events/{eventId} [delete]
func (s *Service) DeleteEvent(c *gin.Context) {
	eventId := c.Param("eventId")

	eventId = strings.ReplaceAll(eventId, "/", "")
	eventIdInt, err := strconv.Atoi(eventId)
	if err != nil {
		log.Fatal("(DeleteEvent) strconv.Atoi", err)
	}
	query := `DELETE FROM events WHERE id = ?`
	_, err = s.db.Exec(query, eventIdInt)
	if err != nil {
		log.Fatal("(DeleteEvent) db.Exec", err)
	}

	c.Status(http.StatusOK)
}
