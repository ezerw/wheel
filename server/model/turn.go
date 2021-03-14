package model

import "time"

type Turn struct {
	ID        int       `db:"id" json:"id,omitempty"`
	Date      time.Time `db:"date" json:"date"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	PersonID  int       `db:"person_id" json:"-"`
	TeamID    int       `db:"team_id" json:"-"`
	Person    Person    `db:"person" json:"person"`
}
