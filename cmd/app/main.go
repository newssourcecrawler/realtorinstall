// cmd/app/main.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// temp vars to track which item is selected
var selectedPropID int64
var selectedBuyerID int64

// Property mirrors internal/models.Property + API JSON.
type Property struct {
	ID           int64  `json:"id"`
	Address      string `json:"address"`
	City         string `json:"city"`
	ZIP          string `json:"zip"`
	ListingDate  string `json:"listing_date"`
	CreatedAt    string `json:"created_at"`
	LastModified string `json:"last_modified"`
}

// Buyer mirrors internal/models.Buyer + API JSON. Adjust fields to match your internal model.
type Buyer struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var (
	// For Properties tab
	properties []Property
	propList   *widget.List

	// For Buyers tab
	buyers    []Buyer
	buyerList *widget.List

	// Base URL of your API
	apiURL = "http://localhost:8080"
)

func main() {
	myApp := app.New()
	myWin := myApp.NewWindow("Realtor Installment Assistant")

	// 1) Build the Properties tab
	propForm := buildPropertyForm(myWin)

	// 2) Build the Buyers tab
	buyerForm := buildBuyerForm(myWin)

	// 3) Combine into tabs
	tabs := container.NewAppTabs(
		container.NewTabItem("Properties", propForm),
		container.NewTabItem("Buyers", buyerForm),
	)
	tabs.SetTabLocation(container.TabLocationTop)

	myWin.SetContent(tabs)
	myWin.Resize(fyne.NewSize(800, 500))

	// Initial load for both tabs
	loadProperties(myWin)
	loadBuyers(myWin)

	// Graceful shutdown on Ctrl+C
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	go myWin.ShowAndRun()
	<-quit
	myWin.Close()
	time.Sleep(200 * time.Millisecond)
}

// buildPropertyForm constructs the UI for the "Properties" tab.
func buildPropertyForm(win fyne.Window) fyne.CanvasObject {
	// Form fields
	addressEntry := widget.NewEntry()
	addressEntry.SetPlaceHolder("Address")
	cityEntry := widget.NewEntry()
	cityEntry.SetPlaceHolder("City")
	zipEntry := widget.NewEntry()
	zipEntry.SetPlaceHolder("ZIP Code")

	addPropBtn := widget.NewButton("Add Property", func() {
		p := map[string]string{
			"address": addressEntry.Text,
			"city":    cityEntry.Text,
			"zip":     zipEntry.Text,
		}
		body, _ := json.Marshal(p)
		resp, err := http.Post(apiURL+"/properties", "application/json", bytes.NewReader(body))
		if err != nil {
			showError(win, fmt.Sprintf("Failed to POST property: %v", err))
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			b, _ := io.ReadAll(resp.Body)
			showError(win, fmt.Sprintf("API error: %s", string(b)))
			return
		}
		// Clear inputs & reload
		addressEntry.SetText("")
		cityEntry.SetText("")
		zipEntry.SetText("")
		loadProperties(win)
	})

	// List widget
	propList = widget.NewList(
		func() int { return len(properties) + 1 },
		func() fyne.CanvasObject { return widget.NewLabel("cell") },
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
			if id == 0 {
				label.SetText("ID | Address | City | ZIP | Listed At")
				label.TextStyle = fyne.TextStyle{Bold: true}
				return
			}
			p := properties[id-1]
			selectedPropID = p.ID
			label.TextStyle = fyne.TextStyle{}
			label.SetText(fmt.Sprintf(
				"%d | %s | %s | %s | %s",
				p.ID, p.Address, p.City, p.ZIP, p.ListingDate,
			))
		},
	)

	editPropBtn := widget.NewButton("Edit Selected", func() {
		if selectedPropID == 0 {
			showError(win, "No property selected")
			return
		}
		// Example: open a popup to edit fields (implementation omitted)
		// call http.NewRequest("PUT", apiURL+fmt.Sprintf("/properties/%d", selectedPropID), <body>)
		// then loadProperties(win)
	})

	deletePropBtn := widget.NewButton("Delete Selected", func() {
		if selectedPropID == 0 {
			showError(win, "No property selected")
			return
		}
		req, _ := http.NewRequest("DELETE", apiURL+fmt.Sprintf("/properties/%d", selectedPropID), nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			showError(win, fmt.Sprintf("DELETE failed: %v", err))
			return
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			showError(win, fmt.Sprintf("API error: %s", resp.Status))
			return
		}
		loadProperties(win)
	})
	refreshPropBtn := widget.NewButton("Refresh Properties", func() {
		loadProperties(win)
	})

	return container.NewVBox(
		widget.NewLabelWithStyle("Properties", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		addressEntry, cityEntry, zipEntry,
		addPropBtn,
		widget.NewSeparator(),
		container.NewHBox(editPropBtn, deletePropBtn, refreshPropBtn),
		propList,
	)
}

// buildBuyerForm constructs the UI for the "Buyers" tab.
func buildBuyerForm(win fyne.Window) fyne.CanvasObject {
	// Form fields
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Full Name")
	emailEntry := widget.NewEntry()
	emailEntry.SetPlaceHolder("Email")

	addBuyerBtn := widget.NewButton("Add Buyer", func() {
		b := map[string]string{
			"name":  nameEntry.Text,
			"email": emailEntry.Text,
		}
		body, _ := json.Marshal(b)
		resp, err := http.Post(apiURL+"/buyers", "application/json", bytes.NewReader(body))
		if err != nil {
			showError(win, fmt.Sprintf("Failed to POST buyer: %v", err))
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			bs, _ := io.ReadAll(resp.Body)
			showError(win, fmt.Sprintf("API error: %s", string(bs)))
			return
		}
		nameEntry.SetText("")
		emailEntry.SetText("")
		loadBuyers(win)
	})

	// List widget
	buyerList = widget.NewList(
		func() int { return len(buyers) + 1 },
		func() fyne.CanvasObject { return widget.NewLabel("cell") },
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
			if id == 0 {
				label.SetText("ID | Name | Email")
				label.TextStyle = fyne.TextStyle{Bold: true}
				return
			}
			b := buyers[id-1]
			selectedBuyerID = b.ID
			//label.TextStyle = fyne.TextStyle{}
			label.SetText(fmt.Sprintf("%d | %s | %s", b.ID, b.Name, b.Email))
		},
	)

	editBuyerBtn := widget.NewButton("Edit Selected", func() {
		if selectedBuyerID == 0 {
			showError(win, "No buyer selected")
			return
		}
		// Example: open a popup to edit buyer (implementation omitted)
		// call http.NewRequest("PUT", apiURL+fmt.Sprintf("/buyers/%d", selectedBuyerID), <body>)
		// then loadBuyers(win)
	})
	deleteBuyerBtn := widget.NewButton("Delete Selected", func() {
		if selectedBuyerID == 0 {
			showError(win, "No buyer selected")
			return
		}
		req, _ := http.NewRequest("DELETE", apiURL+fmt.Sprintf("/buyers/%d", selectedBuyerID), nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			showError(win, fmt.Sprintf("DELETE failed: %v", err))
			return
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			showError(win, fmt.Sprintf("API error: %s", resp.Status))
			return
		}
		loadBuyers(win)
	refreshBuyerBtn := widget.NewButton("Refresh Buyers", func() {
		loadBuyers(win)
	})

	return container.NewVBox(
		widget.NewLabelWithStyle("Buyers", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		nameEntry, emailEntry,
		addBuyerBtn,
		widget.NewSeparator(),
		container.NewHBox(editBuyerBtn, deleteBuyerBtn, refreshBuyerBtn),
		buyerList,
	)
}

// loadProperties calls GET /properties and updates propList.
func loadProperties(win fyne.Window) {
	resp, err := http.Get(apiURL + "/properties")
	if err != nil {
		showError(win, fmt.Sprintf("Failed to GET properties: %v", err))
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		showError(win, fmt.Sprintf("API error: %s", string(b)))
		return
	}
	var ps []Property
	if err := json.NewDecoder(resp.Body).Decode(&ps); err != nil {
		showError(win, fmt.Sprintf("Decode error: %v", err))
		return
	}
	properties = ps
	propList.Refresh()
}

// loadBuyers calls GET /buyers and updates buyerList.
func loadBuyers(win fyne.Window) {
	resp, err := http.Get(apiURL + "/buyers")
	if err != nil {
		showError(win, fmt.Sprintf("Failed to GET buyers: %v", err))
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		showError(win, fmt.Sprintf("API error: %s", string(b)))
		return
	}
	var bs []Buyer
	if err := json.NewDecoder(resp.Body).Decode(&bs); err != nil {
		showError(win, fmt.Sprintf("Decode error: %v", err))
		return
	}
	buyers = bs
	buyerList.Refresh()
}

// showError pops up a modal dialog with an error message.
func showError(w fyne.Window, msg string) {
	var popup *widget.PopUp
	popup = widget.NewModalPopUp(
		container.NewVBox(
			widget.NewLabelWithStyle("Error", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			widget.NewLabel(msg),
			widget.NewButton("OK", func() {
				popup.Hide()
			}),
		),
		w.Canvas(),
	)
	popup.Show()
}
