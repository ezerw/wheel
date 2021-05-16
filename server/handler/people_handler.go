package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"

	"github.com/ezerw/wheel/db"
)

// HandleListPeople handles GET request to /api/teams/:team-id/people
func (s *Server) HandleListPeople(c *gin.Context) {
	teamID := c.Param("team-id")

	intTeamID, err := strconv.ParseInt(teamID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	people, err := s.peopleService.ListPeople(c.Request.Context(), intTeamID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": people})
}

// HandleShowPerson handles GET request to /api/teams/:team-id/people/:person-id
func (s *Server) HandleShowPerson(c *gin.Context) {
	teamID := c.Param("team-id")
	personID := c.Param("person-id")

	intTeamID, err := strconv.ParseInt(teamID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	intPersonID, err := strconv.ParseInt(personID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	person, err := s.peopleService.GetPerson(c.Request.Context(), db.GetPersonParams{
		ID:     intPersonID,
		TeamID: intTeamID,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": person})
}

// HandleAddPerson handles POST request to /api/teams/:team-id/people
func (s *Server) HandleAddPerson(c *gin.Context) {
	teamID := c.Param("team-id")

	intTeamID, err := strconv.ParseInt(teamID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	binding := struct {
		FirstName string `json:"first_name" binding:"required"`
		LastName  string `json:"last_name" binding:"required"`
		Email     string `json:"email" binding:"required"`
	}{}
	err = c.BindJSON(&binding)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	args := db.CreatePersonParams{
		FirstName: binding.FirstName,
		LastName:  binding.LastName,
		Email:     binding.Email,
		TeamID:    intTeamID,
	}

	person, err := s.peopleService.AddPerson(c.Request.Context(), args)
	if err != nil {
		errN, _ := err.(*mysql.MySQLError)
		if errN.Number == 1062 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "The specified email already exists in the database.",
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": person})
}

// HandleUpdatePerson handles PUT request to /api/teams/:team-id/people/:person-id
func (s *Server) HandleUpdatePerson(c *gin.Context) {
	teamID := c.Param("team-id")
	personID := c.Param("person-id")

	intTeamID, err := strconv.ParseInt(teamID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	intPersonID, err := strconv.ParseInt(personID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	binding := struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		TeamID    int64  `json:"team_id"`
	}{}
	err = c.BindJSON(&binding)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the person exists in the specified team
	person, err := s.peopleService.GetPerson(c.Request.Context(), db.GetPersonParams{
		ID:     intPersonID,
		TeamID: intTeamID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Person not found in the specified team.",
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	args := db.UpdatePersonParams{
		FirstName: binding.FirstName,
		LastName:  binding.LastName,
		Email:     binding.Email,
		TeamID:    binding.TeamID,
		ID:        intPersonID,
	}

	person, err = s.peopleService.UpdatePerson(c.Request.Context(), args)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": person})
}

// HandleDeletePerson handles DELETE request to /api/teams/:team-id/people/:person-id
func (s *Server) HandleDeletePerson(c *gin.Context) {
	teamID := c.Param("team-id")
	personID := c.Param("person-id")

	intTeamID, err := strconv.ParseInt(teamID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	intPersonID, err := strconv.ParseInt(personID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	args := db.DeletePersonParams{
		ID:     intPersonID,
		TeamID: intTeamID,
	}

	err = s.peopleService.DeletePerson(c.Request.Context(), args)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": ""})
}
