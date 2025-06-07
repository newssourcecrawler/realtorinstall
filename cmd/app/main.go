package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/your-org/realtorinstall/client"
	"github.com/your-org/realtorinstall/ui"
)

func main() {
	// 1) Load & trust your server’s cert
	certPool := x509.NewCertPool()
	crtPath := "certs/server.crt"
	pem, err := ioutil.ReadFile(crtPath)
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
	myApp.Settings().SetTheme(&ui.CustomTheme{}) // optional

	showLogin(myApp)
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
		widget.NewLabelWithStyle("Realtor Installment", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		username,
		password,
		loginBtn,
		status,
	))
	w.ShowAndRun()
}

// showMain builds your tabbed main UI after login.
func showMain(a fyne.App) {
	w := a.NewWindow("Realtor Installment Assistant")
	w.Resize(fyne.NewSize(1024, 768))

	tabs := container.NewAppTabs(
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
