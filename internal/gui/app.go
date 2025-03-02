package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

type App struct {
	fyneApp fyne.App
	Window  *Window
}

func New() *App {
	a := &App{
		fyneApp: app.New(),
	}
	a.Window = NewWindow(a.fyneApp)
	return a
}

func (a *App) Run() {
	a.Window.Show()
}
