package handler

import (
	"github.com/ezerw/wheel/db"
	"github.com/ezerw/wheel/middleware"
	"github.com/ezerw/wheel/service"
	"github.com/ezerw/wheel/util"
	"github.com/gin-gonic/gin"
	"time"
)

// Server serves HTTP requests for our wheel api.
type Server struct {
	config        util.Config
	router        *gin.Engine
	peopleService *service.People
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(config util.Config, store db.Store) (*Server, error) {
	server := &Server{
		config:        config,
		peopleService: service.NewPeople(store),
	}

	server.setupRouter()
	return server, nil
}

// Start runs the HTTP server on a specific address.
func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func (s *Server) setupRouter() {
	router := gin.Default()
	router.Use(middleware.Cors())

	auth := router.Group("/api")

	// teams
	//auth.GET("/teams", s.HandleListTeams)
	//auth.GET("/teams/:team-id", s.HandleShowTeam)
	//auth.POST("/teams", s.HandleAddTeam)
	//auth.PUT("/teams/:team-id", s.HandleUpdateTeam)
	//auth.DELETE("/teams/:team-id", s.HandleDeleteTeam)

	// team people
	auth.GET("/teams/:team-id/people", s.HandleListPeople)
	auth.GET("/teams/:team-id/people/:person-id", s.HandleShowPerson)
	auth.POST("/teams/:team-id/people", s.HandleAddPerson)
	auth.PUT("/teams/:team-id/people/:person-id", s.HandleUpdatePerson)
	auth.DELETE("/teams/:team-id/people/:person-id", s.HandleDeletePerson)
	//
	//// team turns
	//auth.GET("/teams/:team-id/turns", h.ListTurns)
	//auth.POST("/teams/:team-id/turns", h.CreateTurn)
	//auth.DELETE("/teams/:team-id/turns/:turn-id", h.DeleteTurn)

	s.router = router
}

// checkTeamExists checks whether the passed teamID correspond to an existent team in the DB.
//func (s *Server) checkTeamExists(teamID string) bool {
//	var team model.Team
//
//	err := h.DB.Get(&team, "SELECT id FROM teams WHERE id=?", teamID)
//	if err != nil {
//		if errors.Is(sql.ErrNoRows, err) {
//			return false
//		}
//		return false
//	}
//	return true
//}

// getNextWorkingDay returns the next working day based in the day passed as parameter.
func (s *Server) getNextWorkingDay(today time.Time, loc time.Location) (*time.Time, error) {
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
