// cmd/app/main.go
package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/newssourcecrawler/realtorinstall/cmd/app/client"
	"github.com/newssourcecrawler/realtorinstall/cmd/app/ui"
)

// temp vars to track which item is selected
var selectedPropID int64
var selectedBuyerID int64

func main() {

	// 1) Load & trust your server’s cert
	certPool := x509.NewCertPool()
	crtPath := "certs/server.crt"
	pem, err := os.ReadFile(crtPath)
	if err != nil {
		log.Fatalf("Failed to read CA cert %s: %v", crtPath, err)
	}
	if !certPool.AppendCertsFromPEM(pem) {
		log.Fatalf("Failed to append CA cert")
	}

	// 2) Setup global HTTP client
	client.SetupHTTPClient(&tls.Config{
		RootCAs:    certPool,
		MinVersion: tls.VersionTLS12,
	})

	// 3) Start Fyne
	myApp := app.NewWithID("com.newssourcecrawler.realtorinstall")
	myWin := myApp.NewWindow("Realtor Sales, Lettings and Installment Suite")
	myApp.Settings().SetTheme(&ui.CustomTheme{})
	showLogin(myApp)

	// Graceful shutdown on Ctrl+C
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	go myWin.ShowAndRun()
	<-quit
	myWin.Close()
	time.Sleep(200 * time.Millisecond)
}

func prev() {
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

}

// showLogin displays a simple login window and, on success, calls showMain.
func showLogin(a fyne.App) {
	w := a.NewWindow("Login")
	w.Resize(fyne.NewSize(300, 200))

	username := widget.NewEntry()
	username.SetPlaceHolder("Username")
	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("Password")
	status := widget.NewLabel("")

	loginBtn := widget.NewButton("Login", func() {
		token, err := client.Login(username.Text, password.Text)
		if err != nil {
			status.SetText("❌ " + err.Error())
			return
		}
		client.SetAuthToken(token)
		w.Close()
		showMain(a)
	})

	w.SetContent(container.NewVBox(
		widget.NewLabelWithStyle("Realtor Sales, Lettings and Installment Suite", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		username,
		password,
		loginBtn,
		status,
	))
	w.ShowAndRun()
}

// showMain builds your tabbed main UI after login.
func showMain(a fyne.App) {
	w := a.NewWindow("Realtor Sales, Lettings and Installment Suite")
	w.Resize(fyne.NewSize(1024, 768))

	tabs := container.NewAppTabs(
		container.NewTabItem("Settings", ui.SettingsTab(w)),
		container.NewTabItem("Properties", ui.PropertyTab()),
		container.NewTabItem("Buyers", ui.BuyerTab()),
		container.NewTabItem("Plans", ui.PlanTab()),
		container.NewTabItem("Installments", ui.InstallmentTab()),
		container.NewTabItem("Payments", ui.PaymentTab()),
		container.NewTabItem("Commissions", ui.CommissionTab()),
		container.NewTabItem("Sales", ui.SalesTab()),
		container.NewTabItem("Lettings", ui.LettingsTab()),
		container.NewTabItem("Reports", ui.ReportTab()),
	)
	tabs.SetTabLocation(container.TabLocationTop)

	w.SetContent(tabs)
	w.ShowAndRun()
}

// loadProperties calls GET /properties and updates propList.
func loadProperties(win fyne.Window) {
	resp, err := http.Get(BaseURL + "/properties")
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
	resp, err := http.Get(BaseURL + "/buyers")
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
