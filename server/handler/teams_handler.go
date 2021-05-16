package handler

//// HandleListTeams handles GET request to /api/teams
//func (s *Server) HandleListTeams(c *gin.Context) {
//	teams := []model.Team{}
//	err := h.DB.Select(&teams, "SELECT id, name FROM teams")
//	if err != nil {
//		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{"data": teams})
//}
//
//// HandleShowTeam handles GET request to /api/teams/:team-id
//func (s *Server) HandleShowTeam(c *gin.Context) {
//	teamID := c.Param("team-id")
//
//	var team model.Team
//	err := h.DB.Get(&team, "SELECT id, name FROM teams WHERE id = ?", teamID)
//	if err != nil {
//		if errors.Is(sql.ErrNoRows, err) {
//			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "team not found"})
//			return
//		}
//		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	// attach team people
//	people := []model.Person{}
//	err = h.DB.Select(&people, `SELECT id, first_name, last_name, email FROM people WHERE team_id = ?`, teamID)
//	if err != nil {
//		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	team.People = people
//
//	c.JSON(http.StatusOK, gin.H{"data": team})
//}
//
//// HandleAddTeam handles POST request to /api/teams
//func (s *Server) HandleAddTeam(c *gin.Context) {
//	var binding model.Team
//
//	err := c.BindJSON(&binding)
//	if err != nil {
//		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	team := model.Team{Name: binding.Name}
//	res, err := h.DB.NamedExec("INSERT INTO teams (name) VALUES (:name)", team)
//	if err != nil {
//		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	id, _ := res.LastInsertId()
//	team.ID = int(id)
//
//	c.JSON(http.StatusCreated, gin.H{"data": team})
//}
//
//// HandleUpdateTeam handles PUT request to /api/teams/:team-id
//func (s *Server) HandleUpdateTeam(c *gin.Context) {}
//
//// HandleDeleteTeam handles DELETE request to /api/teams/:team-id
//func (s *Server) HandleDeleteTeam(c *gin.Context) {
//	teamID := c.Param("team-id")
//
//	teamExists := s.checkTeamExists(teamID)
//	if !teamExists {
//		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "team not found"})
//		return
//	}
//
//	_, err := h.DB.Exec("DELETE FROM teams WHERE id = ?", teamID)
//	if err != nil {
//		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	c.Writer.WriteHeader(http.StatusOK)
//}
