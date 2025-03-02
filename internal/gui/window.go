package gui

import (
	"GoRestApi/internal/api"
	"GoRestApi/internal/client"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"io"
)

type Window struct {
	window      fyne.Window
	methodEntry *widget.Select
	urlEntry    *widget.Entry
	bodyEntry   *widget.Entry
	response    *widget.Entry
}

func NewWindow(a fyne.App) *Window {
	w := &Window{
		window: a.NewWindow("API Request Tool"),
	}
	w.setupUI()
	return w
}

func (w *Window) Show() {
	w.window.ShowAndRun()
}

func (w *Window) setupUI() {
	methodLabel := widget.NewLabel("HTTP Method:")
	w.methodEntry = widget.NewSelect([]string{"GET", "POST", "PUT", "DELETE", "PATCH"}, nil)
	w.methodEntry.SetSelected("GET")

	urlLabel := widget.NewLabel("URL:")
	w.urlEntry = widget.NewEntry()
	w.urlEntry.SetPlaceHolder("https://api.example.com/endpoint")

	bodyLabel := widget.NewLabel("Request Body:")
	w.bodyEntry = widget.NewMultiLineEntry()
	w.bodyEntry.SetPlaceHolder("{\n  \"key\": \"value\"\n}")

	sendButton := widget.NewButton("Send Request", func() {
		w.HandleSendReq()
	})

	responseLabel := widget.NewLabel("Response:")
	w.response = widget.NewMultiLineEntry()
	w.response.SetPlaceHolder("Response will appear here...")
	w.response.Disable() // Make it read-only

	form := container.NewVBox(
		methodLabel,
		w.methodEntry,
		urlLabel,
		w.urlEntry,
		bodyLabel,
		w.bodyEntry,
		sendButton,
		responseLabel,
		w.response,
	)

	scrollContainer := container.NewScroll(form)

	w.window.SetContent(scrollContainer)
	w.window.Resize(fyne.NewSize(600, 500))
	w.window.CenterOnScreen()
}

func (w *Window) HandleSendReq() {
	r := api.Req{
		w.methodEntry.Selected,
		w.urlEntry.Text,
		nil,
		w.bodyEntry.Text,
	}

	res, err := client.SendReq(r)
	if err != nil {

	}

	read, err := io.ReadAll(res.Body)

	fmt.Println(string(read))
	w.response.SetText(r.Body)

}
