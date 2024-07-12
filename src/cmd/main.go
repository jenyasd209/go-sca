package main

import (
	"flag"

	"go-sca/src/controller"
	"go-sca/src/database"
	"go-sca/src/model"
	"go-sca/src/repos"
	"go-sca/src/server"
	"go-sca/src/server/handlers"
	"gorm.io/gorm"
)

func main() {
	var dbname string
	flag.StringVar(&dbname, "db", "default.db", "Write name of SQLite db")
	flag.Parse()

	db, err := database.NewDatabase(dbname, &gorm.Config{})
	if err != nil {
		panic(err)
	}

	validator := controller.NewBreedValidator()
	err = validator.Init()
	if err != nil {
		panic(err)
	}

	catsRepo := repos.NewGenericRepo[model.SpyCat](db)
	spyCatHandler := handlers.NewCatHandler(controller.NewSpyCatController(catsRepo, validator))

	targetRepo := repos.NewGenericRepo[model.Target](db)
	targetHandler := handlers.NewTargetHandler(controller.NewTargetController(targetRepo))

	missionsRepo := repos.NewGenericRepo[model.Mission](db)
	missionsHandler := handlers.NewMissionHandler(controller.NewMissionController(missionsRepo))

	s := server.NewServer(":8080", spyCatHandler, missionsHandler, targetHandler)
	err = s.Listen()
	if err != nil {
		panic(err)
	}
}
