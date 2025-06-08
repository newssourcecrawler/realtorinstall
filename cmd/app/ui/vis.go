package ui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/newssourcecrawler/realtorinstall/cmd/app/client"
)

// CustomTheme implements fyne.Theme if you want a branded look.
type CustomTheme struct{}

// … implement Theme methods …

func PropertyTab() fyne.CanvasObject {
	addressEntry := widget.NewEntry()
	addressEntry.SetPlaceHolder("Address")
	cityEntry := widget.NewEntry()
	cityEntry.SetPlaceHolder("City")
	zipEntry := widget.NewEntry()
	zipEntry.SetPlaceHolder("ZIP Code")
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

	var (
		// For Properties tab
		properties []Property
		propList   *widget.List

		// For Buyers tab
		//buyers    []Buyer
		//buyerList *widget.List

		// Base URL of your API
		//apiURL = "https://localhost:8443"
	)

	addPropBtn := widget.NewButton("Add Property", func() {
		p := map[string]string{
			"address": addressEntry.Text,
			"city":    cityEntry.Text,
			"zip":     zipEntry.Text,
		}
		body, _ := json.Marshal(p)
		resp, err := client.HTTPClient.Post(BaseURL+"/properties", "application/json", bytes.NewReader(body))
		if err != nil {
			showError(win, fmt.Sprintf("Failed to POST property: %v", err))
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != client.HTTPClient.StatusOK {
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
		req, _ := client.HTTPClient.NewRequest("DELETE", BaseURL+fmt.Sprintf("/properties/%d", selectedPropID), nil)
		resp, err := client.HTTPClient.Do(req)
		if err != nil {
			showError(win, fmt.Sprintf("DELETE failed: %v", err))
			return
		}
		resp.Body.Close()
		if resp.StatusCode != client.HTTPClient.StatusOK {
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

func BuyerTab() fyne.CanvasObject {
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Full Name")
	emailEntry := widget.NewEntry()
	emailEntry.SetPlaceHolder("Email")

	// Buyer mirrors internal/models.Buyer + API JSON. Adjust fields to match your internal model.
	type Buyer struct {
		ID    int64  `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	var (
		// For Properties tab
		//properties []Property
		//propList   *widget.List

		// For Buyers tab
		buyers    []Buyer
		buyerList *widget.List

		// Base URL of your API
		//apiURL = "https://localhost:8443"
	)

	addBuyerBtn := widget.NewButton("Add Buyer", func() {
		b := map[string]string{
			"name":  nameEntry.Text,
			"email": emailEntry.Text,
		}
		body, _ := json.Marshal(b)
		resp, err := client.HTTPClient.Post(BaseURL+"/buyers", "application/json", bytes.NewReader(body))
		if err != nil {
			showError(win, fmt.Sprintf("Failed to POST buyer: %v", err))
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != client.HTTPClient.StatusOK {
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
		req, _ := client.HTTPClient.NewRequest("DELETE", BaseURL+fmt.Sprintf("/buyers/%d", selectedBuyerID), nil)
		resp, err := client.HTTPClient.Do(req)
		if err != nil {
			showError(win, fmt.Sprintf("DELETE failed: %v", err))
			return
		}
		resp.Body.Close()
		if resp.StatusCode != client.HTTPClient.StatusOK {
			showError(win, fmt.Sprintf("API error: %s", resp.Status))
			return
		}
		loadBuyers(win)
	})
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

func PlanTab() fyne.CanvasObject        { /* … */ }
func InstallmentTab() fyne.CanvasObject { /* … */ }
func PaymentTab() fyne.CanvasObject     { /* … */ }
func CommissionTab() fyne.CanvasObject  { /* … */ }
func SalesTab() fyne.CanvasObject       { /* … */ }
func LettingsTab() fyne.CanvasObject    { /* … */ }
func ReportTab() fyne.CanvasObject      { /* … */ }

// SettingsTab lets the user view/edit the server URL and test connectivity.
func SettingsTab(win fyne.Window) fyne.CanvasObject {
	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("Server Base URL")
	urlEntry.SetText(client.BaseURL)

	status := widget.NewLabel("")

	testBtn := widget.NewButton("Test Connection", func() {
		// A lightweight health check; adjust path if needed
		resp, err := client.HTTPClient.Get(client.BaseURL + "/health")
		if err != nil {
			status.SetText("❌ " + err.Error())
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			status.SetText("✅ OK")
		} else {
			status.SetText(fmt.Sprintf("⚠ %s", resp.Status))
		}
	})

	saveBtn := widget.NewButton("Save", func() {
		client.SetBaseURL(urlEntry.Text)
		status.SetText("Saved: " + client.BaseURL)
	})

	return container.NewVBox(
		widget.NewLabelWithStyle("Connection Settings", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		urlEntry,
		container.NewHBox(testBtn, saveBtn),
		status,
	)
}
