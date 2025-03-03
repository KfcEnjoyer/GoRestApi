package gui

import (
	"GoRestApi/internal/api"
	"GoRestApi/internal/client"
	"GoRestApi/internal/storage"
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
)

type Window struct {
	window               fyne.Window
	reqName              *widget.Entry
	methodEntry          *widget.Select
	urlEntry             *widget.Entry
	bodyEntry            *widget.Entry
	statusCodeEntry      *widget.Entry
	statusContainer      *fyne.Container
	response             *widget.Entry
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

	// Create lists of request names for the sidebar
	requestNames := make([]string, 0, len(w.requestsData))
	for name := range w.requestsData {
		requestNames = append(requestNames, name)
	}

	// Setup request names list with larger text
	w.requestNamesList = widget.NewList(
		func() int {
			return len(requestNames)
		},
		func() fyne.CanvasObject {
			label := widget.NewLabel("Request Name")
			label.TextStyle = fyne.TextStyle{Bold: true}
			// Make text bigger
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

	// Setup request details list (initially empty) with larger text
	w.requestDetailsList = widget.NewList(
		func() int {
			return len(w.currentRequests)
		},
		func() fyne.CanvasObject {
			label := widget.NewLabel("Request Details")
			label.TextStyle = fyne.TextStyle{Bold: true}
			// Make text bigger
			label.Importance = widget.HighImportance
			return container.NewHBox(label)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			req := w.currentRequests[id]
			label := obj.(*fyne.Container).Objects[0].(*widget.Label)
			label.SetText(fmt.Sprintf("%s: %s", req.Method, req.URL))
		},
	)

	// Setup request name selection handler
	w.requestNamesList.OnSelected = func(id widget.ListItemID) {
		w.currentRequestName = requestNames[id]
		w.currentRequests = w.requestsData[w.currentRequestName]
		w.requestDetailsList.Refresh()

		// Reset selection in the details list
		w.selectedRequestIndex = -1
	}

	// Setup request detail selection handler
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

	// Create field for new request name
	w.newReqName = widget.NewEntry()
	w.newReqName.SetPlaceHolder("New request name")

	// Create button to add new request
	addReqButton := widget.NewButton("Create New Request", func() {
		if w.newReqName.Text != "" {
			// Set the form with the new name
			w.reqName.SetText(w.newReqName.Text)
			w.methodEntry.SetSelected("GET")
			w.urlEntry.SetText("")
			w.bodyEntry.SetText("")
			w.newReqName.SetText("")

			// If this is a new name, refresh the list
			if _, exists := w.requestsData[w.reqName.Text]; !exists {
				requestNames = append(requestNames, w.reqName.Text)
				w.requestNamesList.Refresh()
			}
		}
	})
	// Make button bigger
	addReqButton.Importance = widget.HighImportance

	// Create sidebar with two lists
	listsSeparator := widget.NewSeparator()

	// Create labels with larger text
	collectionsLabel := widget.NewLabelWithStyle("Request Collections:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	collectionsLabel.Importance = widget.HighImportance
	requestsLabel := widget.NewLabelWithStyle("Requests in Collection:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	requestsLabel.Importance = widget.HighImportance
	createNewLabel := widget.NewLabelWithStyle("Create New:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	createNewLabel.Importance = widget.HighImportance

	namesListContainer := container.NewGridWrap(fyne.NewSize(250, 250), w.requestNamesList)
	detailsListContainer := container.NewGridWrap(fyne.NewSize(250, 250), w.requestDetailsList)

	// For the new request name entry (fixed height)
	newReqNameContainer := container.NewGridWrap(fyne.NewSize(250, 40), w.newReqName)

	// Create sidebar
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

	// Setup details panel
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
	w.response = widget.NewMultiLineEntry()
	w.response.SetPlaceHolder("Response will appear here...")
	w.response.Disable()

	statusLabel := widget.NewLabelWithStyle("Status Code:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	statusLabel.Importance = widget.HighImportance

	w.statusCodeEntry = widget.NewEntry()

	w.statusContainer = container.NewStack(
		canvas.NewRectangle(color.Transparent),
		w.statusCodeEntry,
	)
	// Create fixed size containers for input fields
	// Replace layout.NewFixedGridLayout with container.NewGridWrap

	// Update other containers to use appropriate layouts
	reqNameContainer := container.NewGridWrap(fyne.NewSize(600, 40), w.reqName)
	methodContainer := container.NewGridWrap(fyne.NewSize(600, 40), w.methodEntry)
	urlContainer := container.NewGridWrap(fyne.NewSize(600, 40), w.urlEntry)
	bodyContainer := container.NewGridWrap(fyne.NewSize(600, 200), w.bodyEntry)
	responseContainer := container.NewGridWrap(fyne.NewSize(600, 300), w.response)

	buttonsContainer := container.NewHBox(saveReqButton, deleteReqButton, sendButton)

	// Create details panel
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
		responseLabel,
		w.statusContainer,
		responseContainer,
	)

	detailsScroll := container.NewScroll(w.detailsPanel)

	// Create split layout with sidebar and details panel
	split := container.NewHSplit(
		container.NewVBox(sidebarContent),
		detailsScroll,
	)
	split.Offset = 0.25 // 25% of space for sidebar, 75% for details

	w.window.SetContent(split)
	w.window.Resize(fyne.NewSize(1920, 1080))
	w.window.CenterOnScreen()
}

func (w *Window) refreshRequestLists() {
	// Reload data
	w.loadSavedRequests()

	// Rebuild request names
	requestNames := make([]string, 0, len(w.requestsData))
	for name := range w.requestsData {
		requestNames = append(requestNames, name)
	}

	// Update current requests if needed
	if w.currentRequestName != "" {
		if reqs, exists := w.requestsData[w.currentRequestName]; exists {
			w.currentRequests = reqs
		} else {
			w.currentRequests = nil
		}
	}

	// Refresh both lists
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

	res, err := client.SendReq(r)
	if err != nil {
		w.response.SetText(err.Error())
		log.Println(err)
		return
	}

	if res != nil {
		w.setStatusCode(res.StatusCode, http.StatusText(res.StatusCode))
		read, err := io.ReadAll(res.Body)
		if err != nil {
			w.response.SetText(err.Error())
			log.Println(err)
			return
		}

		w.response.SetText(string(read))
		return
	}

	w.response.SetText("No response")
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
		return
	}

	dialog.ShowInformation("Success", "Request saved successfully", w.window)

	// Save the current name
	w.currentRequestName = name

	// Refresh lists
	w.refreshRequestLists()
}

func (w *Window) HandleDeleteReq() {
	// Check if we have a selection
	if w.currentRequestName == "" || w.selectedRequestIndex < 0 || w.selectedRequestIndex >= len(w.currentRequests) {
		dialog.ShowInformation("Error", "Please select a request to delete", w.window)
		return
	}

	// Confirm deletion
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
	switch {
	case status >= 200 && status < 300:
		c = color.RGBA{R: 0, G: 200, B: 0, A: 255} // Green
	case status >= 300 && status < 400:
		c = color.NRGBA{R: 255, G: 255, B: 0, A: 255} // Yellow
	case status >= 400 && status < 600:
		c = color.NRGBA{R: 255, G: 0, B: 0, A: 255} // Red
	default:
		c = color.Gray{0xCC}
	}

	// Update the container background
	rect := canvas.NewRectangle(c)
	rect.Resize(w.statusContainer.MinSize())
	w.statusContainer.Objects[0] = rect
	w.statusCodeEntry.SetText(fmt.Sprintf("%d %s", status, text))
	w.statusContainer.Refresh()
}
