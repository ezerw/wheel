package handler

import (
	"context"
	"database/sql"
	"github.com/ezerw/wheel/db"
	"github.com/ezerw/wheel/middleware"
	"github.com/ezerw/wheel/service"
	"github.com/ezerw/wheel/util"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// Server serves HTTP requests for our wheel api.
type Server struct {
	config        util.Config
	router        *gin.Engine
	peopleService *service.People
	teamsService  *service.Teams
	turnsService  *service.Turns
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(config util.Config, store db.Store) (*Server, error) {
	server := &Server{
		config:        config,
		peopleService: service.NewPeople(store),
		teamsService:  service.NewTeams(store),
		turnsService:  service.NewTurns(store),
	}

	server.setupRouter()
	return server, nil
}

// Start runs the HTTP server on a specific address.
func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

// setupRouter defines a router and add routes to it.
func (s *Server) setupRouter() {
	router := gin.Default()
	router.Use(middleware.Cors())

	// TODO: authenticate requests
	auth := router.Group("/api")

	// teams
	auth.GET("/teams", s.HandleListTeams)
	auth.GET("/teams/:team-id", s.HandleShowTeam)
	auth.POST("/teams", s.HandleAddTeam)
	auth.PUT("/teams/:team-id", s.HandleUpdateTeam)
	auth.DELETE("/teams/:team-id", s.HandleDeleteTeam)

	// team people
	auth.GET("/teams/:team-id/people", s.HandleListPeople)
	auth.GET("/teams/:team-id/people/:person-id", s.HandleShowPerson)
	auth.POST("/teams/:team-id/people", s.HandleAddPerson)
	auth.PUT("/teams/:team-id/people/:person-id", s.HandleUpdatePerson)
	auth.DELETE("/teams/:team-id/people/:person-id", s.HandleDeletePerson)

	// team turns
	auth.GET("/teams/:team-id/turns", s.HandleListTurns)
	auth.POST("/teams/:team-id/turns", s.HandleUpsertTurn)

	s.router = router
}

// teamExists checks if the specified teamID exists in the DB.
func (s *Server) teamExists(ctx context.Context, teamID int64) (bool, error) {
	_, err := s.teamsService.GetTeam(ctx, teamID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, nil
	}
	return true, nil
}
