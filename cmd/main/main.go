package main

import (
	"GoRestApi/internal/gui"
	"GoRestApi/internal/setup"
	"log"
)

func main() {
	log.Println("Initializing API Request Tool...")

	setup.EnsureDirectoriesAndFiles()
	app := gui.New()

	app.Run()

	log.Println("Application closed")
}
