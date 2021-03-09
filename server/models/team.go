package models

type Team struct {
	ID     int      `db:"id" json:"id,omitempty"`
	Name   string   `db:"name" json:"name" binding:"required"`
	People []Person `db:"people" json:"people"`
}
