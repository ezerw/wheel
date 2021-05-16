package service

import (
	"context"
	"github.com/ezerw/wheel/db"
)

// People is the service in charge of interact with the people table in the database.
type People struct {
	store db.Store
}

// NewPeople creates a new PeopleService instance.
func NewPeople(store db.Store) *People {
	return &People{store: store}
}

// ListPeople gets people of a team from the DB.
func (s *People) ListPeople(ctx context.Context, teamID int64) ([]db.Person, error) {
	return s.store.ListPeople(ctx, teamID)
}

// GetPerson gets one person of the team from the DB.
func (s *People) GetPerson(ctx context.Context, args db.GetPersonParams) (*db.Person, error) {
	person, err := s.store.GetPerson(ctx, args)
	if err != nil {
		return nil, err
	}

	return &person, nil
}

// AddPerson add one person to the team in the DB.
func (s *People) AddPerson(ctx context.Context, args db.CreatePersonParams) (*db.Person, error) {
	result, err := s.store.CreatePerson(ctx, args)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &db.Person{
		ID:        id,
		FirstName: args.FirstName,
		LastName:  args.LastName,
		Email:     args.Email,
		TeamID:    args.TeamID,
	}, nil
}

// UpdatePerson updates a person from the team in the DB.
func (s *People) UpdatePerson(ctx context.Context, args db.UpdatePersonParams) (*db.Person, error) {
	_, err := s.store.UpdatePerson(ctx, args)
	if err != nil {
		return nil, err
	}

	return &db.Person{
		ID:        args.ID,
		FirstName: args.FirstName,
		LastName:  args.LastName,
		Email:     args.Email,
		TeamID:    args.TeamID,
	}, nil
}

// DeletePerson deletes a person from the team from the DB.
func (s *People) DeletePerson(ctx context.Context, args db.DeletePersonParams) error {
	return s.store.DeletePerson(ctx, args)
}
