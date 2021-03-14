package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	_ "time/tzdata"

	"github.com/ezerw/wheel/models"
)

// TeamTurnList handles GET request to /api/teams/:team-id/turns
func (h *Handler) TurnList(c *gin.Context) {
	teamID := c.Param("team-id")

	teamExists := h.checkTeamExists(teamID)
	if !teamExists {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "team not found"})
		return
	}

	ql := c.DefaultQuery("limit", "10")
	offset := c.DefaultQuery("offset", "0")
	qo := c.DefaultQuery("order", "desc")
	qdf := c.Query("date_from")
	qdt := c.Query("date_to")

	// validate order
	order := "desc"
	for _, o := range []string{"asc", "desc"} {
		if o == qo {
			order = qo
		}
	}

	dateLayout := "2006-01-02"

	turns := []models.Turn{}
	// params
	var args []interface{}

	turnsSQL := `
		SELECT 
			turns.*,
			people.id "person.id",
			people.first_name "person.first_name",
			people.last_name "person.last_name"
		FROM turns
		JOIN people ON turns.person_id = people.id `

	// Flag that says if date_from is present to format the date_to if needed
	hasDf := false

	// Conditionally apply date range filters
	if qdf != "" {
		df, err := time.Parse(dateLayout, qdf)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "date_from invalid format."})
			return
		}
		hasDf = true
		turnsSQL += "WHERE turns.date >= ? "
		args = append(args, df.Format(dateLayout))
	}

	if qdt != "" {
		dt, err := time.Parse(dateLayout, qdt)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "date_from invalid format."})
			return
		}

		if hasDf {
			turnsSQL += "AND turns.date <= ? "
		} else {
			turnsSQL += "WHERE turns.date <= ? "
		}
		args = append(args, dt.Format(dateLayout))
	}

	// add limit and offset last
	args = append(args, ql, offset)
	turnsSQL += "ORDER BY turns.date %s LIMIT ? OFFSET ?;"

	err := h.DB.Select(&turns, fmt.Sprintf(turnsSQL, order), args...)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": turns})
}

// TeamTurnCreate handles POST request to /api/teams/:team-id/turns
func (h *Handler) TurnCreate(c *gin.Context) {
	teamID := c.Param("team-id")

	teamExists := h.checkTeamExists(teamID)
	if !teamExists {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "team not found"})
		return
	}

	binding := struct {
		PersonID int `json:"person_id" binding:"required"`
	}{}

	err := c.BindJSON(&binding)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	person := models.Person{}
	err = h.DB.Get(&person, `SELECT id, first_name, last_name, team_id FROM people WHERE id = ? AND team_id = ?`, binding.PersonID, teamID)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "person not found in specified team"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	loc, err := time.LoadLocation("Pacific/Auckland")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	today := time.Now().In(loc)
	year, month, day := today.Date()
	created := time.Date(year, month, day, 0, 0, 0, 0, loc)

	date, err := h.getNextWorkingDay(today, *loc)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	turn := models.Turn{
		Date:      *date,
		CreatedAt: created,
		PersonID:  person.ID,
		TeamID:    person.TeamID,
	}

	insertTurnSQL := `
		INSERT INTO turns (person_id, team_id, date, created_at) 
		VALUES (:person_id, :team_id, :date, :created_at)
		ON DUPLICATE KEY UPDATE person_id=:person_id`

	res, err := h.DB.NamedExec(insertTurnSQL, &turn)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	lid, err := res.LastInsertId()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": gin.H{
		"id":     lid,
		"date":   date,
		"person": person,
	}})
}

// TeamTurnCreate handles DELETE request to /api/teams/:team-id/turns/turn-id
func (h *Handler) TurnDelete(c *gin.Context) {
	teamID := c.Param("team-id")
	turnID := c.Param("turn-id")

	teamExists := h.checkTeamExists(teamID)
	if !teamExists {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "team not found"})
		return
	}

	rows, err := h.DB.Exec("DELETE FROM turns WHERE id=?", turnID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ar, err := rows.RowsAffected()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if ar == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "turn not found"})
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}
