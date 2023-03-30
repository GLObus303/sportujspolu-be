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

func (service *Service) GetEvents(c *gin.Context) {
	query := "SELECT * FROM events"
	res, err := service.db.Query(query)

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

func (service *Service) GetSingleEvent(c *gin.Context) {
	eventId := c.Param("eventId")
	eventId = strings.ReplaceAll(eventId, "/", "")
	eventIdInt, err := strconv.Atoi(eventId)
	if err != nil {
		log.Fatal("(GetSingleEvent) strconv.Atoi", err)
	}

	var event Event
	query := `SELECT * FROM events WHERE id = ?`
	err = service.db.QueryRow(query, eventIdInt).Scan(&event.Id, &event.Name, &event.Sport)
	if err != nil {
		log.Fatal("(GetSingleEvent) db.Exec", err)
	}

	c.JSON(http.StatusOK, event)
}

func (service *Service) CreateEvent(c *gin.Context) {
	var newEvent Event
	err := c.BindJSON(&newEvent)
	if err != nil {
		log.Fatal("(CreateEvent) c.BindJSON", err)
	}

	query := `INSERT INTO events (name, sport) VALUES (?, ?)`
	res, err := service.db.Exec(query, newEvent.Name, newEvent.Sport)
	if err != nil {
		log.Fatal("(CreateEvent) db.Exec", err)
	}
	newEvent.Id, err = res.LastInsertId()
	if err != nil {
		log.Fatal("(CreateEvent) res.LastInsertId", err)
	}

	c.JSON(http.StatusOK, newEvent)
}

func (service *Service) UpdateEvent(c *gin.Context) {
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
	_, err = service.db.Exec(query, updates.Name, updates.Sport, eventIdInt)
	if err != nil {
		log.Fatal("(UpdateEvent) db.Exec", err)
	}

	c.Status(http.StatusOK)
}

func (service *Service) DeleteEvent(c *gin.Context) {
	eventId := c.Param("eventId")

	eventId = strings.ReplaceAll(eventId, "/", "")
	eventIdInt, err := strconv.Atoi(eventId)
	if err != nil {
		log.Fatal("(DeleteEvent) strconv.Atoi", err)
	}
	query := `DELETE FROM events WHERE id = ?`
	_, err = service.db.Exec(query, eventIdInt)
	if err != nil {
		log.Fatal("(DeleteEvent) db.Exec", err)
	}

	c.Status(http.StatusOK)
}
