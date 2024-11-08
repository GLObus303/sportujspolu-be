package messages

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/globus303/sportujspolu/constants"
	"github.com/globus303/sportujspolu/models"
	"github.com/globus303/sportujspolu/utils"
)

type MessageService struct {
	db *sql.DB
}

func NewMessagesService(db *sql.DB) *MessageService {
	return &MessageService{db}
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
func (s *MessageService) SendEmailRequest(c *gin.Context) {
	var inputEmailRequest models.EmailRequestInput
	if err := c.BindJSON(&inputEmailRequest); err != nil {
		log.Println("(SendEmailRequest) c.BindJSON", err)
		c.JSON(http.StatusBadRequest, utils.GetError("Invalid request data"))

		return
	}

	requesterID := c.GetString(constants.UserID_key)

	var eventOwnerID string
	query := "SELECT owner_id FROM events WHERE public_id = $1 AND owner_id != $2"
	err := s.db.QueryRow(query, inputEmailRequest.EventID, requesterID).Scan(&eventOwnerID)
	if err != nil {
		log.Println("(SendEmailRequest) db.QueryRow", err)
		c.JSON(http.StatusNotFound, utils.GetError("Event not found"))

		return
	}

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
// @Success 200 {object} models.EmailRequestApproveResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /messages/email/{id}/approve [patch]
func (s *MessageService) ApproveEmailRequest(c *gin.Context) {
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
		WHERE id = $4 AND event_owner_id = $5 AND approved IS NULL AND requester_id != $5
		RETURNING
      id,
      text,
      event_id,
      event_owner_id,
      requester_id,
      approved,
      approved_at,
      created_at,
      updated_at,
      (SELECT users.email FROM users WHERE users.id = email_requests.requester_id) AS requester_email
		`

	userID := c.GetString(constants.UserID_key)

	var emailRequest models.EmailRequestApproveResponse
	err := s.db.QueryRow(query, approveInput.Approved, time.Now(), time.Now(), requestId, userID).Scan(
		&emailRequest.ID, &emailRequest.Text, &emailRequest.EventID, &emailRequest.EventOwnerID,
		&emailRequest.RequesterID, &emailRequest.Approved, &emailRequest.ApprovedAt, &emailRequest.CreatedAt, &emailRequest.UpdatedAt, &emailRequest.RequesterEmail,
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

func getEmailRequests(c *gin.Context, s *MessageService, query string) error {
	approvedFilter := c.Query("approvedFilter")

	fmt.Println("approvedFilter", approvedFilter)

	switch approvedFilter {
	case "true":
		query += " AND approved = true"
	case "false":
		query += " AND approved = false"
	case "null":
		query += " AND approved IS NULL"
	}

	userID := c.GetString(constants.UserID_key)

	rows, err := s.db.Query(query, userID)
	if err != nil {
		log.Println("(GetAllEmailRequests) Error querying database:", err)
		c.JSON(http.StatusInternalServerError, utils.GetError("Failed to retrieve email requests"))

		return err
	}
	defer rows.Close()

	var emailRequests []models.EmailRequestResponse

	for rows.Next() {
		var emailRequest models.EmailRequestResponse
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
			&emailRequest.RequesterName,
			&emailRequest.RequesterEmail,
			&emailRequest.EventName,
			&emailRequest.EventLocation,
			&emailRequest.EventLevel,
			&emailRequest.EventOwnerName,
			&emailRequest.EventOwnerEmail,
		)
		if err != nil {
			log.Println("(GetAllEmailRequests) Error scanning row:", err)
			c.JSON(http.StatusInternalServerError, utils.GetError("Failed to parse email requests"))

			return err
		}

		emailRequests = append(emailRequests, emailRequest)
	}

	if len(emailRequests) == 0 {
		c.JSON(http.StatusOK, []models.EmailRequestResponse{})
	} else {
		c.JSON(http.StatusOK, emailRequests)
	}

	return nil
}

// @Summary Get all email requests send as user
// @Description Retrieve all email requests for user from the database
// @Tags messages
// @Produce json
// @Param approvedFilter query string false "Approved filter" Enums(true, false, null) default(null)
// @Success 200 {array} models.EmailRequestResponse "List of email requests"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /messages/email/sent-user-requests [get]
func (s *MessageService) GetAllSentEmailRequests(c *gin.Context) {
	query := `
        SELECT
          email_requests.id,
          email_requests.text,
          email_requests.event_id,
          email_requests.event_owner_id,
          email_requests.requester_id,
          email_requests.approved,
          email_requests.approved_at,
          email_requests.created_at,
          email_requests.updated_at,
          NULL AS requester_name,
          NULL AS requester_email,
          events.name AS event_name,
          events.location AS event_location,
          events.level AS event_level,
          event_owner.name AS event_owner_name,
             (CASE
            WHEN email_requests.approved = true
            THEN event_owner.email
            ELSE NULL
          END) AS event_owner_email

        FROM email_requests
        LEFT JOIN users as event_owner ON event_owner.id = email_requests.event_owner_id
        LEFT JOIN events ON events.public_id = email_requests.event_id
        WHERE email_requests.requester_id = $1 AND events.public_id IS NOT NULL
`

	err := getEmailRequests(c, s, query)
	if err != nil {
		return
	}
}

// @Summary Get all email requests received as owner
// @Description Retrieve all email requests for owner from the database
// @Tags messages
// @Produce json
// @Param approvedFilter query string false "Approved filter" Enums(true, false, null) default(null)
// @Success 200 {array} models.EmailRequestResponse "List of email requests"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /messages/email/received-owner-requests [get]
func (s *MessageService) GetAllReceivedOwnerEmailRequests(c *gin.Context) {
	query := `
        SELECT
          email_requests.id,
          email_requests.text,
          email_requests.event_id,
          email_requests.event_owner_id,
          email_requests.requester_id,
          email_requests.approved,
          email_requests.approved_at,
          email_requests.created_at,
          email_requests.updated_at,
          requester.name AS requester_name,
           (CASE
            WHEN email_requests.approved = true
            THEN requester.email
            ELSE NULL
          END) AS requester_email,
          events.name AS event_name,
          events.location AS event_location,
          events.level AS event_level,
          NULL AS event_owner_name,
          NULL AS event_owner_email

        FROM email_requests
        LEFT JOIN users AS requester ON requester.id =  email_requests.requester_id
        LEFT JOIN events ON events.public_id = email_requests.event_id
        WHERE email_requests.event_owner_id = $1 AND events.public_id IS NOT NULL
`

	err := getEmailRequests(c, s, query)
	if err != nil {
		return
	}
}
