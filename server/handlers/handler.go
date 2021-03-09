package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	_ "time/tzdata"

	"github.com/ezerw/wheel/models"
)

type Handler struct {
	DB *sqlx.DB
}

// TeamList handles GET request to /api/teams
func (h *Handler) TeamList(c *gin.Context) {
	teams := []models.Team{}
	err := h.DB.Select(&teams, "SELECT id, name FROM teams")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": teams})
}

// TeamShow handles GET request to /api/teams/:team-id
func (h *Handler) TeamShow(c *gin.Context) {
	teamID := c.Param("team-id")

	var team models.Team
	err := h.DB.Get(&team, "SELECT id, name FROM teams WHERE id = ?", teamID)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "team not found"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// attach team people
	people := []models.Person{}
	err = h.DB.Select(&people, `SELECT id, first_name, last_name FROM people WHERE team_id = ?`, teamID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	team.People = people

	c.JSON(http.StatusOK, gin.H{"data": team})
}

// TeamCreate handles POST request to /api/teams
func (h *Handler) TeamCreate(c *gin.Context) {
	var binding models.Team

	err := c.BindJSON(&binding)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	team := models.Team{Name: binding.Name}
	res, err := h.DB.NamedExec("INSERT INTO teams (name) VALUES (:name)", team)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, err := res.LastInsertId()
	team.ID = int(id)

	c.JSON(http.StatusCreated, gin.H{"data": team})
}

// TeamDelete handles DELETE request to /api/teams/:team-id
func (h *Handler) TeamDelete(c *gin.Context) {
	teamID := c.Param("team-id")

	teamExists := h.checkTeamExists(teamID)
	if !teamExists {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "team not found"})
		return
	}

	_, err := h.DB.Exec("DELETE FROM teams WHERE id = ?", teamID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}

// TeamAddPerson handles POST request to /api/teams/:team-id/people
func (h *Handler) TeamAddPerson(c *gin.Context) {
	teamID := c.Param("team-id")

	teamExists := h.checkTeamExists(teamID)
	if !teamExists {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "team not found"})
		return
	}

	intTeamID, err := strconv.ParseInt(teamID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var binding models.Person
	err = c.BindJSON(&binding)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	person := models.Person{
		FirstName: binding.FirstName,
		LastName:  binding.LastName,
		TeamID:    int(intTeamID),
	}
	res, err := h.DB.NamedExec(
		"INSERT INTO people (first_name, last_name, team_id) VALUES (:first_name, :last_name, :team_id)",
		person,
	)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newID, err := res.LastInsertId()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	person.ID = int(newID)

	c.JSON(http.StatusCreated, gin.H{"data": person})

}

// TeamDeletePerson handles DELETE request to /api/teams/:team-id/people/:person-id
func (h *Handler) TeamDeletePerson(c *gin.Context) {
	teamID := c.Param("team-id")
	personID := c.Param("person-id")

	teamExists := h.checkTeamExists(teamID)
	if !teamExists {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "team not found"})
		return
	}

	rows, err := h.DB.Exec("DELETE FROM people WHERE id = ?", personID)
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
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "person not found"})
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}

// TeamTurnList handles GET request to /api/teams/:team-id/turns
func (h *Handler) TeamTurnList(c *gin.Context) {
	teamID := c.Param("team-id")

	teamExists := h.checkTeamExists(teamID)
	if !teamExists {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "team not found"})
		return
	}

	ql := c.DefaultQuery("limit", "10")
	offset := c.DefaultQuery("offset", "0")
	qo := c.DefaultQuery("order", "desc")

	// validate limit
	limit := "10"
	for _, l := range []string{"50", "100"} {
		if l == ql {
			limit = l
		}
	}

	// validate order
	order := "desc"
	for _, o := range []string{"asc", "desc"} {
		if o == qo {
			order = qo
		}
	}

	turns := []models.Turn{}
	turnsSQL := `
		SELECT 
			turns.*,
			people.id "person.id",
			people.first_name "person.first_name",
			people.last_name "person.last_name"
		FROM turns
		JOIN people ON turns.person_id = people.id
		ORDER BY turns.date %s
		LIMIT ? 
		OFFSET ?;
	`
	query := fmt.Sprintf(turnsSQL, order)

	err := h.DB.Select(&turns, query, limit, offset)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": turns})
}

// TeamTurnCreate handles POST request to /api/teams/:team-id/turns
func (h *Handler) TeamTurnCreate(c *gin.Context) {
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

	turn := models.Turn{Date: *date, CreatedAt: created, PersonID: person.ID, TeamID: person.TeamID}

	replaceTurnSQL := `
		REPLACE INTO turns (person_id, team_id, date, created_at) 
		VALUES (:person_id, :team_id, :date, :created_at);`

	res, err := h.DB.NamedExec(replaceTurnSQL, &turn)
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
func (h *Handler) TeamTurnDelete(c *gin.Context) {
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

// checkTeamExists checks whether the passed teamID correspond to an existent team in the DB.
func (h *Handler) checkTeamExists(teamID string) bool {
	var team models.Team

	err := h.DB.Get(&team, "SELECT id FROM teams WHERE id=?", teamID)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return false
		}
		return false
	}
	return true
}

// getNextWorkingDay returns the next working day based in the day passed as parameter.
func (h *Handler) getNextWorkingDay(today time.Time, loc time.Location) (*time.Time, error) {
	var next time.Time

	year, month, day := today.Date()
	tomorrow := time.Date(year, month, day+1, 0, 0, 0, 0, &loc)

	switch tomorrow.Weekday() {
	case time.Saturday:
		next = tomorrow.Add(48 * time.Hour)
	case time.Sunday:
		next = tomorrow.Add(24 * time.Hour)
	default:
		next = tomorrow
	}

	return &next, nil
}
