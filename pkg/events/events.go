package events

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/globus303/sportujspolu/constants"
	"github.com/globus303/sportujspolu/models"
	"github.com/globus303/sportujspolu/utils"
)

const columns = "name, sport, date, location, price, description, level, public_id, created_at, owner_id"

func getColumnForEvent(event *models.EventWithOwner) []interface{} {
	return []interface{}{&event.Name, &event.Sport, &event.Date, &event.Location, &event.Price, &event.Description, &event.Level, &event.Public_ID, &event.Created_At, &event.Owner_ID}
}

type Service struct {
	db *sql.DB
}

func NewEventsService(db *sql.DB) *Service {
	return &Service{db}
}

func (s *Service) includeOwner(event *models.EventWithOwner, c *gin.Context) error {
	includes := c.Query("includes")
	if includes != "owner" {
		event.Owner = nil

		return nil
	}

	ownerID := event.Owner_ID

	var owner models.PublicUser
	err := s.db.QueryRow("SELECT name, email, rating FROM users WHERE public_id = ?", ownerID).
		Scan(&owner.Name, &owner.Email, &owner.Rating)

	if err != nil {
		return err
	}

	event.Owner = &owner

	return nil
}

// @Summary Get all events
// @Description Retrieve all events from the database
// @Tags events
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of events per page" default(12)
// @Param includes query string false "Include additional details" Enums(owner)
// @Success 200 {array} models.EventWithOwner
// @Failure 400 {object} models.ErrorResponse
// @Router /events [get]
func (s *Service) GetEvents(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		log.Println("(GetEvents)", err)
		c.JSON(http.StatusBadRequest, utils.GetError("Invalid page parameter"))

		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "12"))
	if err != nil || limit < 1 {
		log.Println("(GetEvents)", err)
		c.JSON(http.StatusBadRequest, utils.GetError("Invalid limit parameter"))

		return
	}

	offset := (page - 1) * limit
	res, err := s.db.Query("SELECT "+columns+" FROM events LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		log.Println("(GetEvents) db.Query", err)
	}

	defer res.Close()

	events := []models.EventWithOwner{}
	for res.Next() {
		var event models.EventWithOwner
		err := res.Scan(getColumnForEvent(&event)...)
		if err != nil {
			log.Println("(GetEvents) res.Scan", err)
		}
		events = append(events, event)
	}

	for i := range events {
		if err := s.includeOwner(&events[i], c); err != nil {
			log.Println("(GetEvents) includeOwner", err)
		}
	}

	c.JSON(http.StatusOK, events)
}

// @Summary Get a single event
// @Description Retrieves a single event from the database
// @Tags events
// @Produce json
// @Param eventId path string true "Event ID" example(q76j5d1a3xtn)
// @Param includes query string false "Include additional details" Enums(owner)
// @Success 200 {object} models.EventWithOwner
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /events/{eventId} [get]
func (s *Service) GetSingleEvent(c *gin.Context) {
	eventId := c.Param("eventId")

	var event models.EventWithOwner
	query := "SELECT " + columns + " FROM events WHERE public_id = $1"
	err := s.db.QueryRow(query, eventId).Scan(getColumnForEvent(&event)...)
	if err != nil {
		log.Println("(GetSingleEvent) db.Exec", err)
		c.JSON(http.StatusNotFound, utils.GetError("Event not found"))

		return
	}

	if err := s.includeOwner(&event, c); err != nil {
		log.Println("(GetEvents) includeOwner", err)
	}

	c.JSON(http.StatusOK, event)
}

type EventInput struct {
	Name        string    `json:"name" example:"Basketball Match at Park"`
	Sport       string    `json:"sport" example:"Basketball"`
	Date        time.Time `json:"date" example:"2023-11-03T10:15:30Z"`
	Location    string    `json:"location" example:"Central Park"`
	Price       uint16    `json:"price" example:"123"`
	Description string    `json:"description" example:"Example Description"`
	Level       string    `json:"level" example:"Any"`
}

// @Summary Create a new event
// @Description Creates a new event in the database
// @Tags events
// @Accept json
// @Produce json
// @Param newEvent body EventInput true "Event object"
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

	userID, _ := c.Get(constants.UserID_key)

	newEvent.Owner_ID = userID.(string)
	newEvent.Public_ID = utils.GenerateUUID()
	newEvent.Created_At = time.Now()

	query := "INSERT INTO events (name, sport, date, location, description, level, public_id, created_at, owner_id"

	values := []interface{}{newEvent.Name, newEvent.Sport, newEvent.Date, newEvent.Location, newEvent.Description, newEvent.Level, newEvent.Public_ID, newEvent.Created_At, newEvent.Owner_ID}

	if newEvent.Price != 0 {
		query += ", price"
		values = append(values, newEvent.Price)
	}

	query += ") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9"
	if newEvent.Price != 0 {
		query += ",$10"
	}
	query += ")"

	fmt.Println(query, values)
	_, err = s.db.Exec(query, values...)
	if err != nil {
		log.Println("(CreateEvent) db.Exec", err)
		c.JSON(http.StatusBadRequest, utils.GetError("Error creating event"))

		return
	}

	c.JSON(http.StatusOK, newEvent)
}

func (s *Service) validateUserIsOwnerOfEvent(c *gin.Context, eventId string) bool {
	userID, _ := c.Get(constants.UserID_key)

	var ownerID string
	err := s.db.QueryRow("SELECT owner_id FROM events WHERE public_id = $1", eventId).Scan(&ownerID)
	if err != nil {
		log.Println("(UpdateEvent) db.QueryRow", err)
		c.JSON(http.StatusBadRequest, utils.GetError("Error updating event"))

		return false
	}

	if ownerID != userID.(string) {
		c.JSON(http.StatusForbidden, utils.GetError("You are not the owner of this event"))

		return false
	}

	return true
}

// @Summary Update an event
// @Description Update an existing event with the given event ID
// @Tags events
// @Accept json
// @Produce json
// @Param eventId path string true "Event ID" example(q76j5d1a3xtn)
// @Param event body EventInput true "Event object that needs to be updated"
// @Success 200 {object} models.Event
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Security BearerAuth
// @Router /events/{eventId} [put]
func (s *Service) UpdateEvent(c *gin.Context) {
	var updates EventInput
	err := c.BindJSON(&updates)
	if err != nil {
		log.Println("(UpdateEvent) c.BindJSON", err)
		c.JSON(http.StatusBadRequest, utils.GetError("Error while parsing request body"))

		return
	}

	eventId := c.Param("eventId")

	if !s.validateUserIsOwnerOfEvent(c, eventId) {
		return
	}

	query := "UPDATE events SET name = $1, sport = $2, date = $3, location = $4, price = $5, description = $6, level = $7"
	values := []interface{}{updates.Name, updates.Sport, updates.Date, updates.Location, updates.Price, updates.Description, updates.Level}

	query += " WHERE public_id = $8"
	values = append(values, eventId)

	_, err = s.db.Exec(query, values...)
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
// @Param eventId path string true "Event ID" example(q76j5d1a3xtn)
// @Success 200 {object} models.Event
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Security BearerAuth
// @Router /events/{eventId} [delete]
func (s *Service) DeleteEvent(c *gin.Context) {
	eventId := c.Param("eventId")

	if !s.validateUserIsOwnerOfEvent(c, eventId) {
		return
	}

	query := `DELETE FROM events WHERE public_id = $1`
	_, err := s.db.Exec(query, eventId)
	if err != nil {
		log.Println("(DeleteEvent) db.Exec", err)
		c.JSON(http.StatusBadRequest, utils.GetError("Error deleting event"))

		return
	}

	c.Status(http.StatusOK)
}
