package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/ezerw/wheel/models"
)

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
