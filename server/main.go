package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	"github.com/ezerw/wheel/db"
	"github.com/ezerw/wheel/handler"
	"github.com/ezerw/wheel/middleware"
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

	h := handler.Handler{DB: conn}

	r := gin.Default()
	r.Use(middleware.Cors())
	r.Use(middleware.Authenticated()) // require a session for api endpoints

	// teams
	r.GET("/teams", h.TeamList)
	r.GET("/teams/:team-id", h.TeamShow)
	r.POST("/teams", h.TeamCreate)
	r.DELETE("/teams/:team-id", h.TeamDelete)

	// team people
	r.POST("/teams/:team-id/people", h.AddPerson)
	r.DELETE("/teams/:team-id/people/:person-id", h.DeletePerson)

	// team turns
	r.GET("/teams/:team-id/turns", h.TurnList)
	r.POST("/teams/:team-id/turns", h.TurnCreate)
	r.DELETE("/teams/:team-id/turns/:turn-id", h.TurnDelete)

	_ = r.Run()
}
