package service

import (
	"context"
	"github.com/ezerw/wheel/db"
)

// Turns is the service in charge of interact with the turns table in the database.
type Turns struct {
	store db.Store
}

// NewTurns creates a new TeamsService instance.
func NewTurns(store db.Store) *Turns {
	return &Turns{store: store}
}

func (s *Turns) ListTurns(ctx context.Context, teamID int64) ([]db.Turn, error) {
}

func (s *Turns) GetTurn(ctx context.Context, turnID int64) (*db.Turn, error) {
}

func (s *Turns) AddTurn(ctx context.Context, args db.CreateTurnParams) (*db.Turn, error) {
}

func (s *Turns) UpdateTurn(ctx context.Context, args db.UpdateTurnParams) (*db.Turn, error) {
}

func (s *Turns) DeleteTurn(ctx context.Context, turnID int64) error {
}