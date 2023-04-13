package events

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Event struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Sport string `json:"sport"`
}

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
// @Success 200 {array} Event
// @Router /events [get]
func (s *Service) GetEvents(c *gin.Context) {
	query := "SELECT * FROM events"
	res, err := s.db.Query(query)

	defer res.Close()

	if err != nil {
		log.Fatal("(GetEvents) db.Query", err)
	}

	events := []Event{}
	for res.Next() {
		var event Event
		err := res.Scan(&event.Id, &event.Name, &event.Sport)
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
// @Failure 404 {object} string
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
	query := `SELECT * FROM events WHERE id = ?`
	err = s.db.QueryRow(query, eventIdInt).Scan(&event.Id, &event.Name, &event.Sport)
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
// @Failure 500 {object} string
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
	newEvent.Id, err = res.LastInsertId()
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
// @Failure 404 {object} string
// @Failure 500 {object} string
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
// @Failure 404 {object} string
// @Failure 500 {object} string
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
