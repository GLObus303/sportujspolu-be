package references

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/globus303/sportujspolu/models"
	"github.com/globus303/sportujspolu/utils"
)

type ReferencesService struct {
	db *sql.DB
}

func NewReferencesService(db *sql.DB) *ReferencesService {
	return &ReferencesService{db}
}

// @Summary Get all levels
// @Description Retrieves all levels from the database
// @Tags levels
// @Produce json
// @Success 200 {array} models.Level
// @Failure 500 {object} models.ErrorResponse
// @Router /references/levels [get]
func (s *ReferencesService) GetAllLevels(c *gin.Context) {
	var levels []models.Level

	query := "SELECT id, value, label FROM levels"
	rows, err := s.db.Query(query)
	if err != nil {
		log.Println("(GetAllLevels) db.Query", err)
		c.JSON(http.StatusInternalServerError, utils.GetError("Error retrieving levels"))

		return
	}
	defer rows.Close()

	for rows.Next() {
		var level models.Level
		if err := rows.Scan(&level.ID, &level.Value, &level.Label); err != nil {
			log.Println("(GetAllLevels) rows.Scan", err)
			c.JSON(http.StatusInternalServerError, utils.GetError("Error processing levels"))

			return
		}

		levels = append(levels, level)
	}

	if err = rows.Err(); err != nil {
		log.Println("(GetAllLevels) rows.Err", err)
		c.JSON(http.StatusInternalServerError, utils.GetError("Error reading levels"))

		return
	}

	c.JSON(http.StatusOK, levels)
}
