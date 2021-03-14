package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/ezerw/wheel/models"
)

// TeamAddPerson handles POST request to /api/teams/:team-id/people
func (h *Handler) AddPerson(c *gin.Context) {
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
func (h *Handler) DeletePerson(c *gin.Context) {
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
