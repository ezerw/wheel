package service

import (
	"context"
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

func (s *Teams) ListTeams(ctx context.Context) ([]db.Team, error) {
}

func (s *Teams) GetTeam(ctx context.Context, teamID int64) (*db.Team, error) {
}

func (s *Teams) AddTeam(ctx context.Context, teamName string) (*db.Team, error) {
}

func (s *Teams) UpdateTeam(ctx context.Context, args db.UpdateTeamParams) (*db.Team, error) {
}

func (s *Teams) DeleteTeam(ctx context.Context, teamID int64) error {
}