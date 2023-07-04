package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/globus303/sportujspolu/middleware"
	"github.com/globus303/sportujspolu/pkg/events"
	adapter "github.com/gwatts/gin-adapter"
	"github.com/joho/godotenv"
	"github.com/jub0bs/fcors"
	"github.com/jub0bs/fcors/risky"
)

func startRouter(db *sql.DB) {
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
	v1 := router.Group("/api/v1")

	// userService := user.NewUserService(db)

	// users := v1.Group("/user")
	// users.POST("/register", userService.Register)
	// users.POST("/login", userService.Login)

	eventsService := events.NewEventsService(db)

	events := v1.Group("/events")
	events.GET("", eventsService.GetEvents)
	events.GET("/:eventId", eventsService.GetSingleEvent)

	protectedEvents := events.Group("")
	protectedEvents.Use(middleware.JwtAuth())

	protectedEvents.POST("", eventsService.CreateEvent)
	protectedEvents.PUT("/:eventId", eventsService.UpdateEvent)
	protectedEvents.DELETE("/:eventId", eventsService.DeleteEvent)

	//	@Summary Health check
	//	@Description Returns the status of the server.
	//	@Tags	health
	//	@Success 200
	//	@Failure 500 "Internal Server Error"
	//	@Router	/health [get]
	v1.GET("/health", func(c *gin.Context) {
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

// @title SportujSpolu API
// @description	This is the API for the SportujSpolu app.
// @version 1.0
// @host sportujspolu-api.onrender.com
// @BasePath /api/v1
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

	startRouter(db)
}
