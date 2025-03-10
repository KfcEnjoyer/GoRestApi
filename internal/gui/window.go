package gui

import (
	"GoRestApi/internal/api"
	"GoRestApi/internal/storage"
	"GoRestApi/internal/utills"
	"bytes"
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"io"
	"log"
	"net/http"
	"time"
)

type Window struct {
	window               fyne.Window
	reqName              *widget.Entry
	methodEntry          *widget.Select
	urlEntry             *widget.Entry
	bodyEntry            *widget.Entry
	statusCodeText       *canvas.Text
	statusContainer      *fyne.Container
	responseText         *widget.Label
	requestNamesList     *widget.List
	requestDetailsList   *widget.List
	requestsData         map[string][]api.Req
	newReqName           *widget.Entry
	detailsPanel         *fyne.Container
	currentRequestName   string
	currentRequests      []api.Req
	selectedRequestIndex int
}

func NewWindow(a fyne.App) *Window {
	w := &Window{
		window:               a.NewWindow("API Request Tool"),
		selectedRequestIndex: -1,
	}

	iconPath := "internal/gui/gorest.png"

	icon, err := fyne.LoadResourceFromPath(iconPath)
	if err == nil {
		w.window.SetIcon(icon)
	}

	w.setupUI()
	return w
}

func (w *Window) Show() {
	w.window.ShowAndRun()
}

func (w *Window) loadSavedRequests() {
	var err error
	w.requestsData, err = storage.LoadRequests()
	if err != nil {
		log.Println("Error loading requests:", err)
	}
}

func (w *Window) setupUI() {
	w.loadSavedRequests()

	requestNames := make([]string, 0, len(w.requestsData))
	for name := range w.requestsData {
		requestNames = append(requestNames, name)
	}

	w.requestNamesList = widget.NewList(
		func() int {
			return len(requestNames)
		},
		func() fyne.CanvasObject {
			label := widget.NewLabel("Request Name")
			label.TextStyle = fyne.TextStyle{Bold: true}

			label.Importance = widget.HighImportance
			return container.NewHBox(label)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			name := requestNames[id]
			count := len(w.requestsData[name])
			label := obj.(*fyne.Container).Objects[0].(*widget.Label)
			label.SetText(fmt.Sprintf("%s (%d)", name, count))
		},
	)

	w.requestDetailsList = widget.NewList(
		func() int {
			return len(w.currentRequests)
		},
		func() fyne.CanvasObject {
			label := widget.NewLabel("Request Details")
			label.TextStyle = fyne.TextStyle{Bold: true}
			label.Importance = widget.HighImportance
			return container.NewHBox(label)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			req := w.currentRequests[id]
			label := obj.(*fyne.Container).Objects[0].(*widget.Label)
			label.SetText(fmt.Sprintf("%s: %s", req.Method, req.URL))
		},
	)

	w.requestNamesList.OnSelected = func(id widget.ListItemID) {
		w.currentRequestName = requestNames[id]
		w.currentRequests = w.requestsData[w.currentRequestName]
		w.requestDetailsList.Refresh()

		w.selectedRequestIndex = -1
	}

	w.requestDetailsList.OnSelected = func(id widget.ListItemID) {
		if id < len(w.currentRequests) {
			w.selectedRequestIndex = int(id)
			req := w.currentRequests[id]
			w.reqName.SetText(w.currentRequestName)
			w.methodEntry.SetSelected(req.Method)
			w.urlEntry.SetText(req.URL)
			w.bodyEntry.SetText(req.Body)
		}
	}

	w.newReqName = widget.NewEntry()
	w.newReqName.SetPlaceHolder("New request name")

	addReqButton := widget.NewButton("Create New Request", func() {
		if w.newReqName.Text != "" {
			w.reqName.SetText(w.newReqName.Text)
			w.methodEntry.SetSelected("GET")
			w.urlEntry.SetText("")
			w.bodyEntry.SetText("")
			w.newReqName.SetText("")

			if _, exists := w.requestsData[w.reqName.Text]; !exists {
				requestNames = append(requestNames, w.reqName.Text)
				w.requestNamesList.Refresh()
			}
		}
	})
	addReqButton.Importance = widget.HighImportance

	listsSeparator := widget.NewSeparator()

	collectionsLabel := widget.NewLabelWithStyle("Request Collections:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	collectionsLabel.Importance = widget.HighImportance
	requestsLabel := widget.NewLabelWithStyle("Requests in Collection:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	requestsLabel.Importance = widget.HighImportance
	createNewLabel := widget.NewLabelWithStyle("Create New:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	createNewLabel.Importance = widget.HighImportance

	namesListContainer := container.NewGridWrap(fyne.NewSize(250, 250), w.requestNamesList)
	detailsListContainer := container.NewGridWrap(fyne.NewSize(250, 250), w.requestDetailsList)

	newReqNameContainer := container.NewGridWrap(fyne.NewSize(250, 40), w.newReqName)

	sidebarContent := container.NewVBox(
		collectionsLabel,
		namesListContainer,
		listsSeparator,
		requestsLabel,
		detailsListContainer,
		createNewLabel,
		newReqNameContainer,
		addReqButton,
	)

	nameLabel := widget.NewLabelWithStyle("Request Name:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	nameLabel.Importance = widget.HighImportance
	w.reqName = widget.NewEntry()
	w.reqName.SetPlaceHolder("example/post")

	methodLabel := widget.NewLabelWithStyle("HTTP Method:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	methodLabel.Importance = widget.HighImportance
	w.methodEntry = widget.NewSelect([]string{"GET", "POST", "PUT", "DELETE", "PATCH"}, nil)
	w.methodEntry.SetSelected("GET")

	urlLabel := widget.NewLabelWithStyle("URL:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	urlLabel.Importance = widget.HighImportance
	w.urlEntry = widget.NewEntry()
	w.urlEntry.SetPlaceHolder("https://api.example.com/endpoint")

	bodyLabel := widget.NewLabelWithStyle("Request Body:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	bodyLabel.Importance = widget.HighImportance
	w.bodyEntry = widget.NewMultiLineEntry()
	w.bodyEntry.SetPlaceHolder("{\n  \"key\": \"value\"\n}")

	saveReqButton := widget.NewButton("Save Request", func() {
		w.HandleSaveReq()
	})
	saveReqButton.Importance = widget.HighImportance

	deleteReqButton := widget.NewButton("Delete Selected Request", func() {
		w.HandleDeleteReq()
	})
	deleteReqButton.Importance = widget.HighImportance

	sendButton := widget.NewButton("Send Request", func() {
		w.HandleSendReq()
	})
	sendButton.Importance = widget.HighImportance

	responseLabel := widget.NewLabelWithStyle("Response:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	responseLabel.Importance = widget.HighImportance
	w.responseText = widget.NewLabel("Response will appear here...")
	w.responseText.Wrapping = fyne.TextWrapWord
	w.responseText.Alignment = fyne.TextAlignLeading
	w.responseText.TextStyle = fyne.TextStyle{Monospace: true}

	responseTextContainer := container.NewGridWrap(fyne.NewSize(400, 750), w.responseText)
	responseTextContainer.Resize(fyne.NewSize(600, -1))

	responseScrollContainer := container.NewVScroll(responseTextContainer)
	responseScrollContainer.SetMinSize(fyne.NewSize(400, 600))

	statusLabel := widget.NewLabelWithStyle("Status Code:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	statusLabel.Importance = widget.HighImportance

	w.statusCodeText = canvas.NewText("", color.White)
	w.statusCodeText.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}
	w.statusCodeText.Alignment = fyne.TextAlignLeading

	statusRect := canvas.NewRectangle(color.Transparent)
	statusRect.SetMinSize(fyne.NewSize(600, 30))

	w.statusContainer = container.NewGridWrap(fyne.NewSize(250, 20),
		statusRect,
		w.statusCodeText,
	)

	w.statusContainer.Resize(fyne.NewSize(600, 30))

	responseContainer := container.NewVBox(
		responseLabel,
		w.statusContainer,
		responseScrollContainer,
	)

	reqNameContainer := container.NewGridWrap(fyne.NewSize(600, 40), w.reqName)
	methodContainer := container.NewGridWrap(fyne.NewSize(600, 40), w.methodEntry)
	urlContainer := container.NewGridWrap(fyne.NewSize(600, 40), w.urlEntry)
	bodyContainer := container.NewGridWrap(fyne.NewSize(600, 200), w.bodyEntry)
	buttonsContainer := container.NewHBox(saveReqButton, deleteReqButton, sendButton)

	w.detailsPanel = container.NewVBox(
		nameLabel,
		reqNameContainer,
		methodLabel,
		methodContainer,
		urlLabel,
		urlContainer,
		bodyLabel,
		bodyContainer,
		buttonsContainer,
		responseContainer,
	)

	detailsScroll := container.NewScroll(w.detailsPanel)

	split := container.NewHSplit(
		container.NewVBox(sidebarContent),
		detailsScroll,
	)
	split.Offset = 0.25

	w.window.SetContent(split)
	w.window.Resize(fyne.NewSize(1920, 1080))
	w.window.CenterOnScreen()
}

func (w *Window) refreshRequestLists() {
	w.loadSavedRequests()

	requestNames := make([]string, 0, len(w.requestsData))
	for name := range w.requestsData {
		requestNames = append(requestNames, name)
	}

	if w.currentRequestName != "" {
		if reqs, exists := w.requestsData[w.currentRequestName]; exists {
			w.currentRequests = reqs
		} else {
			w.currentRequests = nil
		}
	}

	w.requestNamesList.Refresh()
	w.requestDetailsList.Refresh()
}

func (w *Window) HandleSendReq() {
	r := api.Req{
		Method:  w.methodEntry.Selected,
		URL:     w.urlEntry.Text,
		Headers: nil,
		Body:    w.bodyEntry.Text,
	}

	res, err := api.SendReq(r)
	if err != nil {
		w.responseText.Text = err.Error()
		w.setStatusCode(500, "Internal Error")
		w.responseText.Refresh()
		utills.CreateReqLog(utills.RequestLogger{
			Method:     r.Method,
			Url:        r.URL,
			Message:    "Failed to create request: " + err.Error(),
			StatusCode: 500, // Or another appropriate code
			TimeStamp:  time.Now().Format("2006-01-02 15:04:05"),
		})
		return
	}

	if res != nil {
		w.setStatusCode(res.StatusCode, http.StatusText(res.StatusCode))

		read, err := io.ReadAll(res.Body)
		if err != nil {
			w.responseText.SetText(err.Error())
			utills.CreateReqLog(utills.RequestLogger{
				Method:     r.Method,
				Url:        r.URL,
				Message:    "Failed to create request: " + err.Error(),
				StatusCode: 500, // Or another appropriate code
				TimeStamp:  time.Now().Format("2006-01-02 15:04:05"),
			})
			return
		}

		var prettyJSON bytes.Buffer
		err = json.Indent(&prettyJSON, read, "", "  ")
		if err != nil {
			w.responseText.Text = string(read)
			utills.CreateReqLog(utills.RequestLogger{
				Method:     r.Method,
				Url:        r.URL,
				Message:    "Failed to create request: " + err.Error(),
				StatusCode: 500, // Or another appropriate code
				TimeStamp:  time.Now().Format("2006-01-02 15:04:05"),
			})
		} else {
			w.responseText.Text = prettyJSON.String()
		}

		w.responseText.Refresh()
		return
	}

	w.responseText.Text = "No response"
	w.responseText.Refresh()
}

func (w *Window) HandleSaveReq() {
	name := w.reqName.Text
	if name == "" {
		dialog.ShowInformation("Error", "Request name cannot be empty", w.window)
		return
	}

	r := api.Req{
		Method:  w.methodEntry.Selected,
		URL:     w.urlEntry.Text,
		Headers: nil,
		Body:    w.bodyEntry.Text,
	}

	err := storage.SaveRequest(name, r)
	if err != nil {
		dialog.ShowError(err, w.window)
		utills.CreateLog(utills.ErrorLogger{
			Error:     err.Error(),
			TimeStamp: time.Now().Format("2006-01-02 15:04:05"),
		})
		return
	}

	dialog.ShowInformation("Success", "Request saved successfully", w.window)

	w.currentRequestName = name

	w.refreshRequestLists()
}

func (w *Window) HandleDeleteReq() {
	if w.currentRequestName == "" || w.selectedRequestIndex < 0 || w.selectedRequestIndex >= len(w.currentRequests) {
		dialog.ShowInformation("Error", "Please select a request to delete", w.window)
		return
	}

	dialog.ShowConfirm("Confirm Deletion",
		"Are you sure you want to delete this request?",
		func(confirmed bool) {
			if confirmed {
				err := storage.DeleteRequest(w.currentRequestName, w.selectedRequestIndex)
				if err != nil {
					dialog.ShowError(err, w.window)
					return
				}

				w.selectedRequestIndex = -1
				w.refreshRequestLists()
				dialog.ShowInformation("Success", "Request deleted successfully", w.window)
			}
		},
		w.window)
}

func (w *Window) setStatusCode(status int, text string) {
	var c color.Color
	var textColor color.Color
	switch {
	case status >= 200 && status < 300:
		c = color.NRGBA{R: 0, G: 100, B: 0, A: 200}         // Dark green
		textColor = color.NRGBA{R: 0, G: 255, B: 0, A: 255} // Bright green
	case status >= 300 && status < 400:
		c = color.NRGBA{R: 100, G: 100, B: 0, A: 200}         // Dark yellow
		textColor = color.NRGBA{R: 255, G: 255, B: 0, A: 255} // Bright yellow
	case status >= 400 && status < 600:
		c = color.NRGBA{R: 100, G: 0, B: 0, A: 200}         // Dark red
		textColor = color.NRGBA{R: 255, G: 0, B: 0, A: 255} // Bright red
	default:
		c = color.Gray{0x55} // Medium gray
		textColor = color.White
	}

	statusRect := w.statusContainer.Objects[0].(*canvas.Rectangle)
	statusRect.FillColor = c
	statusRect.Refresh()

	w.statusCodeText.Text = fmt.Sprintf("%d %s", status, text)
	w.statusCodeText.Color = textColor
	w.statusCodeText.Refresh()
}
