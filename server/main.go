package main

import (
	"github.com/ezerw/wheel/db"
	"github.com/ezerw/wheel/handlers"
	"github.com/ezerw/wheel/middlewares"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	conn, err := db.Init()
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	handler := handlers.Handler{DB: conn}

	r := gin.Default()
	r.Use(middlewares.Cors())
	r.Use(middlewares.Authenticated()) // require a session for api endpoints

	// teams
	r.GET("/teams", handler.TeamList)
	r.GET("/teams/:team-id", handler.TeamShow)
	r.POST("/teams", handler.TeamCreate)
	r.DELETE("/teams/:team-id", handler.TeamDelete)

	// team people
	r.POST("/teams/:team-id/people", handler.AddPerson)
	r.DELETE("/teams/:team-id/people/:person-id", handler.DeletePerson)

	// team turns
	r.GET("/teams/:team-id/turns", handler.TurnList)
	r.POST("/teams/:team-id/turns", handler.TurnCreate)
	r.DELETE("/teams/:team-id/turns/:turn-id", handler.TurnDelete)

	_ = r.Run()
}
