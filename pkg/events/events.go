package events

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/globus303/sportujspolu/models"
	"github.com/globus303/sportujspolu/utils"
	_ "github.com/go-sql-driver/mysql"
)

type Service struct {
	db *sql.DB
}

func NewEventsService(db *sql.DB) *Service {
	return &Service{db}
}

// @Summary Get all events
// @Description Retrieve all events from the database
// @Tags events
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of events per page" default(10)
// @Success 200 {array} models.Event
// @Failure 400 {object} models.ErrorResponse
// @Router /events [get]
func (s *Service) GetEvents(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		log.Println("(GetEvents)", err)
		c.JSON(http.StatusBadRequest, utils.GetError("Invalid page parameter"))

		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		log.Println("(GetEvents)", err)
		c.JSON(http.StatusBadRequest, utils.GetError("Invalid limit parameter"))

		return
	}

	offset := (page - 1) * limit
	res, err := s.db.Query("SELECT * FROM events LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		log.Println("(GetEvents) db.Query", err)
	}

	defer res.Close()

	events := []models.Event{}
	for res.Next() {
		var event models.Event
		err := res.Scan(&event.ID, &event.Name, &event.Sport, &event.Date, &event.Location, &event.Price, &event.Description, &event.Level, &event.Public_ID)
		if err != nil {
			log.Println("(GetEvents) res.Scan", err)

		}
		events = append(events, event)
	}

	c.JSON(http.StatusOK, events)
}

// @Summary Get a single event
// @Description Retrieves a single event from the database
// @Tags events
// @Produce json
// @Param eventId path string true "Event ID" example("q76j5d1a3xtn")
// @Success 200 {object} models.Event
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /events/{eventId} [get]
func (s *Service) GetSingleEvent(c *gin.Context) {
	eventId := c.Param("eventId")

	var event models.Event
	query := `SELECT * FROM events WHERE public_id = ?`
	err := s.db.QueryRow(query, eventId).Scan(&event.ID, &event.Name, &event.Sport, &event.Date, &event.Location, &event.Price, &event.Description, &event.Level, &event.Public_ID)
	if err != nil {
		log.Println("(GetSingleEvent) db.Exec", err)
		c.JSON(http.StatusNotFound, utils.GetError("Event not found"))

		return
	}

	c.JSON(http.StatusOK, event)
}

// @Summary Create a new event
// @Description Creates a new event in the database
// @Tags events
// @Accept json
// @Produce json
// @Param newEvent body models.Event true "Event object"
// @Success 200 {object} models.Event
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Security BearerAuth
// @Router /events [post]
func (s *Service) CreateEvent(c *gin.Context) {
	var newEvent models.Event
	err := c.BindJSON(&newEvent)
	if err != nil {
		log.Println("(CreateEvent) c.BindJSON", err)
	}

	query := `INSERT INTO events (name, sport) VALUES (?, ?)`
	res, err := s.db.Exec(query, newEvent.Name, newEvent.Sport)
	if err != nil {
		log.Println("(CreateEvent) db.Exec", err)
		c.JSON(http.StatusBadRequest, utils.GetError("Error creating event"))

		return
	}

	newEvent.ID, err = res.LastInsertId()
	if err != nil {
		log.Println("(CreateEvent) res.LastInsertId", err)

		return
	}

	c.JSON(http.StatusOK, newEvent)
}

// @Summary Update an event
// @Description Update an existing event with the given event ID
// @Tags events
// @Accept json
// @Produce json
// @Param eventId path int true "Event ID"
// @Param event body models.Event true "Event object that needs to be updated"
// @Success 200 {object} models.Event
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Security BearerAuth
// @Router /events/{eventId} [put]
func (s *Service) UpdateEvent(c *gin.Context) {
	var updates models.Event
	err := c.BindJSON(&updates)
	if err != nil {
		log.Println("(UpdateEvent) c.BindJSON", err)
		c.JSON(http.StatusBadRequest, utils.GetError("Error while parsing request body"))

		return
	}

	eventId := c.Param("eventId")

	query := `UPDATE events SET name = ?, sport = ? WHERE public_id = ?`
	_, err = s.db.Exec(query, updates.Name, updates.Sport, eventId)
	if err != nil {
		log.Println("(UpdateEvent) db.Exec", err)
		c.JSON(http.StatusBadRequest, utils.GetError("Error updating event"))

		return
	}

	c.Status(http.StatusOK)
}

// @Summary Delete an event
// @Description Delete an existing event with the given event ID
// @Tags events
// @Param eventId path int true "Event ID"
// @Success 200 {object} models.Event
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Security BearerAuth
// @Router /events/{eventId} [delete]
func (s *Service) DeleteEvent(c *gin.Context) {
	eventId := c.Param("eventId")

	query := `DELETE FROM events WHERE public_id = ?`
	_, err := s.db.Exec(query, eventId)
	if err != nil {
		log.Println("(DeleteEvent) db.Exec", err)
		c.JSON(http.StatusBadRequest, utils.GetError("Error deleting event"))

		return
	}

	c.Status(http.StatusOK)
}
