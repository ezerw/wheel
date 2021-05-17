package service

import (
	"context"
	"github.com/ezerw/wheel/db"
	"time"
)

// Turns is the service in charge of interact with the turns table in the database.
type Turns struct {
	store db.Store
}

// NewTurns creates a new TeamsService instance.
func NewTurns(store db.Store) *Turns {
	return &Turns{store: store}
}

// ListTurns gets turns from the DB based on passed params.
func (s *Turns) ListTurns(
	ctx context.Context,
	teamID int64,
	dateFrom time.Time,
	dateTo time.Time,
	limit int64,
	offset int64,
) ([]db.Turn, error) {
	if !dateFrom.IsZero() && !dateTo.IsZero() {
		args := db.ListTurnsWithBothDatesParams{
			TeamID: teamID,
			Date:   dateFrom,
			Date_2: dateTo,
			Limit:  int32(limit),
			Offset: int32(offset),
		}
		return s.store.ListTurnsWithBothDates(ctx, args)
	}

	if !dateFrom.IsZero() {
		args := db.ListTurnsWithDateFromParams{
			TeamID: teamID,
			Date:   dateFrom,
			Limit:  int32(limit),
			Offset: int32(offset),
		}
		return s.store.ListTurnsWithDateFrom(ctx, args)
	}

	if !dateTo.IsZero() {
		args := db.ListTurnsWithDateToParams{
			TeamID: teamID,
			Date:   dateTo,
			Limit:  int32(limit),
			Offset: int32(offset),
		}
		return s.store.ListTurnsWithDateTo(ctx, args)
	}

	args := db.ListTurnsParams{
		TeamID: teamID,
		Limit:  int32(limit),
		Offset: int32(offset),
	}
	return s.store.ListTurns(ctx, args)
}

// GetTurn gets one turn from the DB using id and teamID as params.
func (s *Turns) GetTurn(ctx context.Context, args db.GetTurnParams) (*db.Turn, error) {
	turn, err := s.store.GetTurn(ctx, args)
	if err != nil {
		return nil, err
	}

	return &turn, nil
}

// GetTurnByDate gets one turn from the DB using date and teamID as params.
func (s *Turns) GetTurnByDate(ctx context.Context, args db.GetTurnByDateParams) (*db.Turn, error) {
	turn, err := s.store.GetTurnByDate(ctx, args)
	if err != nil {
		return nil, err
	}

	return &turn, nil
}

// AddTurn adds a turn to the DB for the specified team.
func (s *Turns) AddTurn(ctx context.Context, teamID int64, personID int64, date time.Time) (*db.Turn, error) {
	args := db.CreateTurnParams{
		PersonID: personID,
		TeamID:   teamID,
		Date:     date,
	}

	result, err := s.store.CreateTurn(ctx, args)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	getTurnArgs := db.GetTurnParams{
		ID:     id,
		TeamID: teamID,
	}

	turn, err := s.GetTurn(ctx, getTurnArgs)
	if err != nil {
		return nil, err
	}

	return turn, nil
}

// UpdateTurn updates a turn in the DB.
func (s *Turns) UpdateTurn(ctx context.Context, args db.UpdateTurnParams) (*db.Turn, error) {
	_, err := s.store.UpdateTurn(ctx, args)
	if err != nil {
		return nil, err
	}

	return &db.Turn{
		ID:       args.ID,
		PersonID: args.PersonID,
		TeamID:   args.TeamID,
		Date:     args.Date,
	}, nil
}
