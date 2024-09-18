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
	"github.com/globus303/sportujspolu/pkg/messages"
	"github.com/globus303/sportujspolu/pkg/references"
	"github.com/globus303/sportujspolu/pkg/user"
	adapter "github.com/gwatts/gin-adapter"
	"github.com/joho/godotenv"
	"github.com/jub0bs/fcors"
	"github.com/jub0bs/fcors/risky"
	_ "github.com/lib/pq"
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
			http.MethodPatch,
		),
		fcors.WithRequestHeaders("Authorization", "Content-Type", "cache"),
		risky.SkipPublicSuffixCheck(),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	router := gin.Default()
	router.Use(adapter.Wrap(cors))
	v1 := router.Group("/api/v1")

	userService := user.NewUserService(db)

	user := v1.Group("/user")
	user.POST("/register", userService.Register)
	user.POST("/login", userService.Login)

	protectedUser := user.Group("").Use(middleware.JwtAuth())
	protectedUser.GET("/me", userService.GetMe)
	protectedUser.DELETE("/me", userService.DeleteMe)

	referencesService := references.NewReferencesService(db)

	levels := v1.Group("/references")
	levels.GET("/levels", referencesService.GetAllLevels)

	eventsService := events.NewEventsService(db)

	events := v1.Group("/events")
	events.GET("", eventsService.GetAllEvents)
	events.GET("/:eventId", eventsService.GetSingleEvent)

	protectedEvents := events.Group("")
	protectedEvents.Use(middleware.JwtAuth())

	protectedEvents.POST("", eventsService.CreateEvent)
	protectedEvents.PUT("/:eventId", eventsService.UpdateEvent)
	protectedEvents.DELETE("/:eventId", eventsService.DeleteEvent)

	messagesService := messages.NewMessagesService(db)

	protectedMessages := v1.Group("/messages").Use(middleware.JwtAuth())
	protectedMessages.POST("/email/request", messagesService.SendEmailRequest)
	protectedMessages.PATCH("/email/:requestId/approve", messagesService.ApproveEmailRequest)
	protectedMessages.GET("/email/sent-user-requests", messagesService.GetAllSentEmailRequests)
	protectedMessages.GET("/email/received-owner-requests", messagesService.GetAllReceivedOwnerEmailRequests)

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
// @schemes https
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load env", err)
	}

	db, err := sql.Open("postgres", os.Getenv("DB_CONNECTION"))
	if err != nil {
		log.Fatal("failed to open db connection", err)
	}
	defer db.Close()

	rows, err := db.Query("select version()")
	if err != nil {
		log.Fatal("failed to open db connection", err)
	}
	defer rows.Close()

	var version string
	for rows.Next() {
		err := rows.Scan(&version)
		if err != nil {
			log.Fatal("failed to open db connection", err)
		}
	}
	fmt.Printf("version=%s\n", version)

	startRouter(db)
}
