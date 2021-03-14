package db

import (
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func Init() (*sqlx.DB, error) {
	var (
		db *sqlx.DB
		err error
	)

	for i := 0; i < 10; i++ {
		db, err = sqlx.Connect("mysql", os.Getenv("MYSQL_DSN"))
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	// create teams table
	createTeamsTableSQL := `
	CREATE TABLE IF NOT EXISTS teams (
	    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	    name VARCHAR(100) NOT NULL,
	    UNIQUE (name)
	);`
	db.MustExec(createTeamsTableSQL)

	// create people table
	createPeopleTableSQL := `
	CREATE TABLE IF NOT EXISTS people (
		id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		first_name VARCHAR(100) NOT NULL,
		last_name VARCHAR(100) NOT NULL,
		team_id INT NOT NULL,
		FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE
	);`
	db.MustExec(createPeopleTableSQL)

	// create turns table
	createTurnsTableSQL := `
	CREATE TABLE IF NOT EXISTS turns (
		id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		person_id INT NOT NULL,
		team_id INT NOT NULL,
		date DATE NOT NULL,
		created_at DATE NOT NULL,
		UNIQUE(person_id, date),
		FOREIGN KEY (person_id) REFERENCES people(id) ON DELETE CASCADE
	);`
	db.MustExec(createTurnsTableSQL)

	return db, nil
}
