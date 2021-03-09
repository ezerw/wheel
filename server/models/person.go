package models

type Person struct {
	ID        int    `db:"id" json:"id,omitempty"`
	FirstName string `db:"first_name" json:"first_name" binding:"required"`
	LastName  string `db:"last_name" json:"last_name" binding:"required"`
	TeamID    int    `db:"team_id" json:"-"`
}
