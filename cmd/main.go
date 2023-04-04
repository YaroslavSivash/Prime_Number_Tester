package main

import (
	"log"

	"Prime_Number_Tester/internal/config"
	"Prime_Number_Tester/server"
)

func main() {
	conf, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatal(err)
	}
	app, err := server.NewApp(conf)
	if err != nil {
		log.Fatal(err)
	}
	if err = app.Run(conf.AppPort); err != nil {
		log.Fatal(err)
	}
}
