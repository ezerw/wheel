package handler

import (
	_ "time/tzdata" // required

)

//var dateLayout = "2006-01-02"

//// HandleListTurns handles GET request to /api/teams/:team-id/turns
//func (s *Server) HandleListTurns(c *gin.Context) {
//	teamID := c.Param("team-id")
//
//	teamExists := h.checkTeamExists(teamID)
//	if !teamExists {
//		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "team not found"})
//		return
//	}
//
//	ql := c.DefaultQuery("limit", "10")
//	offset := c.DefaultQuery("offset", "0")
//	qo := c.DefaultQuery("order", "desc")
//	qdf := c.Query("date_from")
//	qdt := c.Query("date_to")
//
//	// validate order
//	order := "desc"
//	for _, o := range []string{"asc", "desc"} {
//		if o == qo {
//			order = qo
//		}
//	}
//
//	turns := []model.Turn{}
//	// params
//	var args []interface{}
//
//	turnsSQL := `
//		SELECT
//			turns.id,
//			turns.date,
//			turns.created_at,
//			people.id "person.id",
//			people.first_name "person.first_name",
//			people.last_name "person.last_name"
//		FROM turns
//		JOIN people ON turns.person_id = people.id `
//
//	// Flag that says if date_from is present to format the date_to if needed
//	hasDf := false
//
//	// Conditionally apply date range filters
//	if qdf != "" {
//		df, err := time.Parse(dateLayout, qdf)
//		if err != nil {
//			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "date_from invalid format."})
//			return
//		}
//		hasDf = true
//		turnsSQL += "WHERE turns.date >= ? "
//		args = append(args, df.Format(dateLayout))
//	}
//
//	if qdt != "" {
//		dt, err := time.Parse(dateLayout, qdt)
//		if err != nil {
//			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "date_to invalid format."})
//			return
//		}
//
//		if hasDf {
//			turnsSQL += "AND turns.date <= ? "
//		} else {
//			turnsSQL += "WHERE turns.date <= ? "
//		}
//		args = append(args, dt.Format(dateLayout))
//	}
//
//	// add limit and offset last
//	args = append(args, ql, offset)
//	turnsSQL += "ORDER BY turns.date %s LIMIT ? OFFSET ?;"
//
//	err := h.DB.Select(&turns, fmt.Sprintf(turnsSQL, order), args...)
//	if err != nil {
//		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{"data": turns})
//}
//
//// HandleShowTurn handles GET request to /api/teams/:team-id/turns/:turn-id
//func (s *Server) HandleShowTurn(c *gin.Context){}
//
//// HandleCreateTurn handles POST request to /api/teams/:team-id/turns
//func (s *Server) HandleCreateTurn(c *gin.Context) {
//	teamID := c.Param("team-id")
//
//	teamExists := h.checkTeamExists(teamID)
//	if !teamExists {
//		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "team not found"})
//		return
//	}
//
//	binding := struct {
//		PersonID int `json:"person_id" binding:"required"`
//	}{}
//
//	err := c.BindJSON(&binding)
//	if err != nil {
//		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	person := model.Person{}
//	err = h.DB.Get(&person, `SELECT id, first_name, last_name, team_id FROM people WHERE id = ? AND team_id = ?`, binding.PersonID, teamID)
//	if err != nil {
//		if errors.Is(sql.ErrNoRows, err) {
//			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "person not found in specified team"})
//			return
//		}
//		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	loc, err := time.LoadLocation("Pacific/Auckland")
//	if err != nil {
//		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	today := time.Now().In(loc)
//	date, err := h.getNextWorkingDay(today, *loc)
//	if err != nil {
//		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	turn := model.Turn{
//		Date:      *date,
//		CreatedAt: today,
//		PersonID:  person.ID,
//		TeamID:    person.TeamID,
//	}
//
//	insertTurnSQL := `
//		REPLACE INTO turns (person_id, team_id, date, created_at)
//		VALUES (:person_id, :team_id, :date, :created_at);`
//
//	res, err := h.DB.NamedExec(insertTurnSQL, &turn)
//	if err != nil {
//		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//	lid, err := res.LastInsertId()
//	if err != nil {
//		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	turnSQL := `
//		SELECT
//			turns.id,
//			turns.date,
//			turns.created_at,
//			people.id "person.id",
//			people.first_name "person.first_name",
//			people.last_name "person.last_name"
//		FROM turns
//		JOIN people ON turns.person_id = people.id
//		WHERE turns.id = ?`
//
//	t := model.Turn{}
//	err = h.DB.Get(&t, turnSQL, lid)
//	if err != nil {
//		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	c.JSON(http.StatusCreated, gin.H{"data": t})
//}
//
//// HandleUpdateTurn handles PUT request to /api/teams/:team-id/turns/:turn-id
//func (s *Server) HandleUpdateTurn(c *gin.Context) {}
//
//// HandleDeleteTurn handles DELETE request to /api/teams/:team-id/turns/:turn-id
//func (s *Server) HandleDeleteTurn(c *gin.Context) {
//	teamID := c.Param("team-id")
//	turnID := c.Param("turn-id")
//
//	teamExists := h.checkTeamExists(teamID)
//	if !teamExists {
//		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "team not found"})
//		return
//	}
//
//	rows, err := h.DB.Exec("DELETE FROM turns WHERE id=?", turnID)
//	if err != nil {
//		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	ar, err := rows.RowsAffected()
//	if err != nil {
//		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	if ar == 0 {
//		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "turn not found"})
//		return
//	}
//
//	c.Writer.WriteHeader(http.StatusOK)
//}
