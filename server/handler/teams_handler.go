package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/ezerw/wheel/db"
)

// HandleListTeams handles GET request to /api/teams
func (s *Server) HandleListTeams(c *gin.Context) {
	teams, err := s.teamsService.ListTeams(c.Request.Context())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": teams})
}

// HandleShowTeam handles GET request to /api/teams/:team-id
func (s *Server) HandleShowTeam(c *gin.Context) {
	teamID := c.Param("team-id")
	intTeamID, err := strconv.ParseInt(teamID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	team, err := s.teamsService.GetTeam(c.Request.Context(), intTeamID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Team not found.",
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Add people to the team response
	composite := struct {
		TeamID   int64       `json:"id"`
		TeamName string      `json:"name"`
		People   []db.Person `json:"people"`
	}{}

	composite.TeamName = team.Name
	composite.TeamID = team.ID

	people, _ := s.peopleService.ListPeople(c.Request.Context(), intTeamID)
	composite.People = people

	c.JSON(http.StatusOK, gin.H{"data": composite})
}

// HandleAddTeam handles POST request to /api/teams
func (s *Server) HandleAddTeam(c *gin.Context) {
	binding := struct {
		Name string `json:"name" binding:"required"`
	}{}
	err := c.BindJSON(&binding)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	team, err := s.teamsService.AddTeam(c.Request.Context(), binding.Name)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": team})
}

// HandleUpdateTeam handles PUT request to /api/teams/:team-id
func (s *Server) HandleUpdateTeam(c *gin.Context) {
	teamID := c.Param("team-id")
	intTeamID, err := strconv.ParseInt(teamID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	binding := struct {
		Name string `json:"name,omitempty"`
	}{}
	err = c.BindJSON(&binding)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = s.teamsService.GetTeam(c.Request.Context(), intTeamID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Team not found.",
			})
			return
		}
	}

	updateTeamArgs := db.UpdateTeamParams{
		Name: binding.Name,
		ID:   intTeamID,
	}

	team, err := s.teamsService.UpdateTeam(c.Request.Context(), updateTeamArgs)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": team})
}

// HandleDeleteTeam handles DELETE request to /api/teams/:team-id
func (s *Server) HandleDeleteTeam(c *gin.Context) {
	teamID := c.Param("team-id")

	intTeamID, err := strconv.ParseInt(teamID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = s.teamsService.DeleteTeam(c.Request.Context(), intTeamID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": ""})
}
