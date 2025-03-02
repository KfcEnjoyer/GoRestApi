package main

import (
	"GoRestApi/internal/gui"
	"log"
)

func main() {
	log.Println("Initializing API Request Tool...")

	app := gui.New()

	app.Run()

	log.Println("Application closed")
}
