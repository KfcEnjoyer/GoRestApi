package main

import (
	"GoRestApi/internal/gui"
	"GoRestApi/internal/utills"
	"log"
)

func main() {
	log.Println("Initializing API Request Tool...")

	utills.EnsureDirectoriesAndFiles()
	app := gui.New()

	app.Run()

	log.Println("Application closed")
}
