package handlers

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	_ "time/tzdata"

	"github.com/ezerw/wheel/models"
)

type Handler struct {
	DB *sqlx.DB
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
