package service

import (
	"context"
	
	"github.com/pkg/errors"

	"github.com/ezerw/wheel/db"
)

// Teams is the service in charge of interact with the teams table in the database.
type Teams struct {
	store db.Store
}

// NewTeams creates a new TeamsService instance.
func NewTeams(store db.Store) *Teams {
	return &Teams{store: store}
}

// ListTeams gets a list of teams from the DB.
func (s *Teams) ListTeams(ctx context.Context) ([]db.Team, error) {
	return s.store.ListTeams(ctx)
}

// GetTeam gets a team from the DB.
func (s *Teams) GetTeam(ctx context.Context, teamID int64) (*db.Team, error) {
	team, err := s.store.GetTeam(ctx, teamID)
	if err != nil {
		return nil, err
	}

	return &team, nil
}

// AddTeam adds a team to the DB.
func (s *Teams) AddTeam(ctx context.Context, teamName string) (*db.Team, error) {
	result, err := s.store.CreateTeam(ctx, teamName)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &db.Team{
		ID:   id,
		Name: teamName,
	}, nil
}

// UpdateTeam updates a team name in the DB.
func (s *Teams) UpdateTeam(ctx context.Context, args db.UpdateTeamParams) (*db.Team, error) {
	result, err := s.store.UpdateTeam(ctx, args)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, errors.New("0 rows affected")
	}

	return &db.Team{
		ID:   args.ID,
		Name: args.Name,
	}, nil
}

// DeleteTeam deletes a team from the DB
func (s *Teams) DeleteTeam(ctx context.Context, teamID int64) error {
	return s.store.DeleteTeam(ctx, teamID)
}
