// cmd/app/main.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// Property matches the JSON shape from the API (and internal/models.Property).
type Property struct {
	ID           int64  `json:"id"`
	Address      string `json:"address"`
	City         string `json:"city"`
	ZIP          string `json:"zip"`
	ListingDate  string `json:"listing_date"`
	CreatedAt    string `json:"created_at"`
	LastModified string `json:"last_modified"`
}

var (
	// properties holds the last‐fetched slice of Property.
	properties []Property

	// propList is the Fyne List widget; we keep it global so loadProperties can call propList.Refresh().
	propList *widget.List

	// apiURL is the base URL for the API. Change if your server is remote.
	apiURL = "http://localhost:8080"
)

func main() {
	// 1. Create a Fyne application
	myApp := app.New()
	myWin := myApp.NewWindow("Realtor Installment Assistant")

	// 2. Build the "Add Property" form fields and button
	addressEntry := widget.NewEntry()
	addressEntry.SetPlaceHolder("Address")

	cityEntry := widget.NewEntry()
	cityEntry.SetPlaceHolder("City")

	zipEntry := widget.NewEntry()
	zipEntry.SetPlaceHolder("ZIP Code")

	addBtn := widget.NewButton("Add Property", func() {
		// Gather input from each field
		p := map[string]string{
			"address": addressEntry.Text,
			"city":    cityEntry.Text,
			"zip":     zipEntry.Text,
		}

		// POST to /properties
		bodyBytes, _ := json.Marshal(p)
		resp, err := http.Post(apiURL+"/properties", "application/json", bytes.NewReader(bodyBytes))
		if err != nil {
			dialog.ShowError(fmt.Errorf("failed to POST: %w", err), myWin)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			b, _ := io.ReadAll(resp.Body)
			dialog.ShowError(fmt.Errorf("API error: %s", string(b)), myWin)
			return
		}

		// Clear form on success
		addressEntry.SetText("")
		cityEntry.SetText("")
		zipEntry.SetText("")

		// Refresh the list
		loadProperties(myWin)
	})

	// 3. Create the List widget (initially empty)
	propList = widget.NewList(
		// Length: number of rows = len(properties)+1 (header)
		func() int {
			return len(properties) + 1
		},
		// Create template: a single label per row
		func() fyne.CanvasObject {
			return widget.NewLabel("cell")
		},
		// Update each row by ID
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
			if id == 0 {
				// Header row
				label.SetText("ID | Address | City | ZIP | Listed At")
				label.TextStyle = fyne.TextStyle{Bold: true}
				return
			}
			// Data rows
			p := properties[id-1]
			label.TextStyle = fyne.TextStyle{}
			label.SetText(fmt.Sprintf(
				"%d | %s | %s | %s | %s",
				p.ID, p.Address, p.City, p.ZIP, p.ListingDate,
			))
		},
	)

	// 4. Create a "Refresh Properties" button
	refreshBtn := widget.NewButton("Refresh Properties", func() {
		loadProperties(myWin)
	})

	// 5. Arrange all widgets in a vertical box
	form := container.NewVBox(
		widget.NewLabelWithStyle("Add New Property", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		addressEntry,
		cityEntry,
		zipEntry,
		addBtn,
		widget.NewSeparator(),
		refreshBtn,
		propList,
	)

	myWin.SetContent(form)
	myWin.Resize(fyne.NewSize(600, 400))

	// 6. Initial load of properties
	loadProperties(myWin)

	myWin.ShowAndRun()
}

// loadProperties fetches properties from GET /properties and refreshes propList.
func loadProperties(win fyne.Window) {
	resp, err := http.Get(apiURL + "/properties")
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to GET: %w", err), win)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		dialog.ShowError(fmt.Errorf("API error: %s", string(b)), win)
		return
	}

	var ps []Property
	if err := json.NewDecoder(resp.Body).Decode(&ps); err != nil {
		dialog.ShowError(fmt.Errorf("decode error: %w", err), win)
		return
	}

	properties = ps
	propList.Refresh()
}
