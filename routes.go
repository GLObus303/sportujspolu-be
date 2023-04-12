package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/globus303/sportujspolu/pkg/events"
	adapter "github.com/gwatts/gin-adapter"
	"github.com/joho/godotenv"
	"github.com/jub0bs/fcors"
	"github.com/jub0bs/fcors/risky"
)

func startGin(db *sql.DB) {
	allowedOrigins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")

	cors, err := fcors.AllowAccess(
		fcors.FromOrigins(allowedOrigins[0], allowedOrigins[1:]...),
		fcors.WithMethods(
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
		),
		fcors.WithRequestHeaders("Authorization"),
		risky.SkipPublicSuffixCheck(),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	router := gin.Default()
	router.Use(adapter.Wrap(cors))

	eventsService := events.NewEventsService(db)

	router.GET("/events", eventsService.GetEvents)
	router.GET("/events/:eventId", eventsService.GetSingleEvent)
	router.POST("/events", eventsService.CreateEvent)
	router.PUT("/events/:eventId", eventsService.UpdateEvent)
	router.DELETE("/events/:eventId", eventsService.DeleteEvent)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}
	if err := router.Run(":" + port); err != nil {
		log.Panicf("error: %s", err)
	}

	fmt.Println("Server is running on port " + port)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load env", err)
	}

	db, err := sql.Open("mysql", os.Getenv("DSN"))
	if err != nil {
		log.Fatal("failed to open db connection", err)
	}

	defer db.Close()

	startGin(db)

}
