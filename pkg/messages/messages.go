package messages

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/globus303/sportujspolu/constants"
	"github.com/globus303/sportujspolu/models"
	"github.com/globus303/sportujspolu/utils"
)

type Service struct {
	db *sql.DB
}

func NewMessagesService(db *sql.DB) *Service {
	return &Service{db}
}

// @Summary Send an email request
// @Description Sends an email request to join an event
// @Tags messages
// @Accept json
// @Produce json
// @Param newEmailRequest body models.EmailRequestInput true "Email Request object"
// @Success 200 {object} models.EmailRequest
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /messages/email/request [post]
func (s *Service) SendEmailRequest(c *gin.Context) {
	var inputEmailRequest models.EmailRequestInput
	if err := c.BindJSON(&inputEmailRequest); err != nil {
		log.Println("(SendEmailRequest) c.BindJSON", err)
		c.JSON(http.StatusBadRequest, utils.GetError("Invalid request data"))
		return
	}

	var eventOwnerID string
	query := "SELECT owner_id FROM events WHERE public_id = $1"
	err := s.db.QueryRow(query, inputEmailRequest.EventID).Scan(&eventOwnerID)
	if err != nil {
		log.Println("(SendEmailRequest) db.QueryRow", err)
		c.JSON(http.StatusNotFound, utils.GetError("Event not found"))

		return
	}

	requesterID := c.GetString(constants.UserID_key)

	query = `
    SELECT id
    FROM email_requests
    WHERE event_id = $1 AND requester_id = $2
`
	var existingRequestID int
	err = s.db.QueryRow(query, inputEmailRequest.EventID, requesterID).Scan(&existingRequestID)

	if err == nil {
		log.Println("(SendEmailRequest) db.QueryRow (request already exists)")
		c.JSON(http.StatusConflict, utils.GetError("Request already exists for this event and requester"))

		return
	}

	if err != sql.ErrNoRows {
		log.Println("(SendEmailRequest) db.QueryRow", err)
		c.JSON(http.StatusInternalServerError, utils.GetError("Error checking existing email request"))

		return
	}

	newEmailRequest := models.EmailRequest{
		Text:         inputEmailRequest.Text,
		EventID:      inputEmailRequest.EventID,
		EventOwnerID: eventOwnerID,
		RequesterID:  requesterID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	query = `
		INSERT INTO email_requests (text, event_id, event_owner_id, requester_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	err = s.db.QueryRow(query, newEmailRequest.Text, newEmailRequest.EventID, newEmailRequest.EventOwnerID, newEmailRequest.RequesterID, newEmailRequest.CreatedAt, newEmailRequest.UpdatedAt).Scan(&newEmailRequest.ID)
	if err != nil {
		log.Println("(SendEmailRequest) db.QueryRow", err)
		c.JSON(http.StatusInternalServerError, utils.GetError("Error sending email request"))
		return
	}

	c.JSON(http.StatusOK, newEmailRequest)
}

// @Summary Approve an email request
// @Description Approves an email request for a given ID
// @Tags messages
// @Accept json
// @Produce json
// @Param id path int true "Email Request ID" default(1)
// @Param approveInput body models.EmailRequestApproveInput true "Approval status"
// @Success 200 {object} models.EmailRequest
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /messages/email/{id}/approve [patch]
func (s *Service) ApproveEmailRequest(c *gin.Context) {
	requestId := c.Param("requestId")

	var approveInput models.EmailRequestApproveInput
	if err := c.BindJSON(&approveInput); err != nil {
		log.Println("(UpdateEmailRequest) c.BindJSON", err)
		c.JSON(http.StatusBadRequest, utils.GetError("Invalid request data"))
		return
	}

	query := `
		UPDATE email_requests
		SET approved = $1, approved_at = $2, updated_at = $3
		WHERE id = $4 AND event_owner_id = $5 AND approved IS NULL
		RETURNING id, text, event_id, event_owner_id, requester_id, approved, approved_at, created_at, updated_at
		`

	userID := c.GetString(constants.UserID_key)

	var emailRequest models.EmailRequest
	err := s.db.QueryRow(query, approveInput.Approved, time.Now(), time.Now(), requestId, userID).Scan(
		&emailRequest.ID, &emailRequest.Text, &emailRequest.EventID, &emailRequest.EventOwnerID,
		&emailRequest.RequesterID, &emailRequest.Approved, &emailRequest.ApprovedAt, &emailRequest.CreatedAt, &emailRequest.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("(UpdateEmailRequest) No record found", err)
			c.JSON(http.StatusNotFound, utils.GetError("Email request not found"))

			return
		}

		log.Println("(UpdateEmailRequest) db.QueryRow", err)
		c.JSON(http.StatusInternalServerError, utils.GetError("Error approving email request"))

		return
	}

	c.JSON(http.StatusOK, emailRequest)
}

// @Summary Get all email requests
// @Description Retrieve all email requests from the database
// @Tags messages
// @Produce json
// @Success 200 {array} models.EmailRequest "List of email requests"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /messages/email/user-requests [get]
func (s *Service) GetAllEmailRequests(c *gin.Context) {
	query := `
			SELECT id, text, event_id, event_owner_id, requester_id, approved, approved_at, created_at, updated_at
			FROM email_requests
			WHERE approved IS NULL AND event_owner_id = $1
	`

	userID := c.GetString(constants.UserID_key)

	rows, err := s.db.Query(query, userID)
	if err != nil {
		log.Println("(GetAllEmailRequests) Error querying database:", err)
		c.JSON(http.StatusInternalServerError, utils.GetError("Failed to retrieve email requests"))

		return
	}
	defer rows.Close()

	var emailRequests []models.EmailRequest

	for rows.Next() {
		var emailRequest models.EmailRequest
		err := rows.Scan(
			&emailRequest.ID,
			&emailRequest.Text,
			&emailRequest.EventID,
			&emailRequest.EventOwnerID,
			&emailRequest.RequesterID,
			&emailRequest.Approved,
			&emailRequest.ApprovedAt,
			&emailRequest.CreatedAt,
			&emailRequest.UpdatedAt,
		)
		if err != nil {
			log.Println("(GetAllEmailRequests) Error scanning row:", err)
			c.JSON(http.StatusInternalServerError, utils.GetError("Failed to parse email requests"))

			return
		}

		emailRequests = append(emailRequests, emailRequest)
	}

	if len(emailRequests) == 0 {
		c.JSON(http.StatusOK, []models.EmailRequest{})
	} else {
		c.JSON(http.StatusOK, emailRequests)
	}
}
