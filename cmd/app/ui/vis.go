package ui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/newssourcecrawler/realtorinstall/cmd/app/client"
)

// CustomTheme implements fyne.Theme if you want a branded look.
type CustomTheme struct{}

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

// PlanTab builds the Plans screen with list and CRUD actions.
func PlanTab(win fyne.Window) fyne.CanvasObject {
	// Domain model for UI
	type Plan struct {
		ID          int64   `json:"id"`
		PropertyID  int64   `json:"property_id"`
		BuyerID     int64   `json:"buyer_id"`
		TotalAmount float64 `json:"total_amount"`
	}

	var items []Plan
	list := widget.NewList(
		func() int { return len(items) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(i widget.ListItemID, obj fyne.CanvasObject) {
			p := items[i]
			obj.(*widget.Label).SetText(
				fmt.Sprintf("#%d – Prop:%d – Buyer:%d – $%.2f",
					p.ID, p.PropertyID, p.BuyerID, p.TotalAmount),
			)
		},
	)

	// Load data from backend
	load := func() {
		resp, err := client.HTTPClient.Get(client.BaseURL + "/plans")
		if err != nil {
			dialog.ShowError(err, win)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			dialog.ShowError(fmt.Errorf("server error: %s", resp.Status), win)
			return
		}
		var data []Plan
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			dialog.ShowError(err, win)
			return
		}
		items = data
		list.Refresh()
	}

	var selectedID int64
	list.OnSelected = func(id widget.ListItemID) {
		selectedID = items[id].ID
	}

	// Form helper to create or edit
	showForm := func(title string, p *Plan, save func(Plan)) {
		propEntry := widget.NewEntry()
		propEntry.SetText(strconv.FormatInt(p.PropertyID, 10))
		buyerEntry := widget.NewEntry()
		buyerEntry.SetText(strconv.FormatInt(p.BuyerID, 10))
		amountEntry := widget.NewEntry()
		amountEntry.SetText(fmt.Sprintf("%.2f", p.TotalAmount))

		dlg := dialog.NewForm(
			title,
			"Save",
			"Cancel",
			[]*widget.FormItem{
				{Text: "Property ID", Widget: propEntry},
				{Text: "Buyer ID", Widget: buyerEntry},
				{Text: "Total Amount", Widget: amountEntry},
			},
			func(bOK bool) {
				if !bOK {
					dlg.Hide()
					return
				}
				pid, _ := strconv.ParseInt(propEntry.Text, 10, 64)
				bid, _ := strconv.ParseInt(buyerEntry.Text, 10, 64)
				amt, _ := strconv.ParseFloat(amountEntry.Text, 64)
				n := Plan{ID: p.ID, PropertyID: pid, BuyerID: bid, TotalAmount: amt}
				save(n)
				dlg.Hide()
			}, win)
		dlg.Show()
	}

	// CRUD buttons
	btnRefresh := widget.NewButton("Refresh", func() { load() })
	btnNew := widget.NewButton("New", func() {
		showForm("New Plan", &Plan{}, func(n Plan) {
			buf, _ := json.Marshal(n)
			resp, err := client.HTTPClient.Post(client.BaseURL+"/plans", "application/json", bytes.NewReader(buf))
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			resp.Body.Close()
			load()
		})
	})
	btnEdit := widget.NewButton("Edit", func() {
		if selectedID == 0 {
			dialog.ShowInformation("Select one", "Please select a plan to edit.", win)
			return
		}
		var current Plan
		for _, pl := range items {
			if pl.ID == selectedID {
				current = pl
				break
			}
		}
		showForm("Edit Plan", &current, func(n Plan) {
			buf, _ := json.Marshal(n)
			req, _ := http.NewRequest(http.MethodPut, client.BaseURL+fmt.Sprintf("/plans/%d", n.ID), bytes.NewReader(buf))
			req.Header.Set("Content-Type", "application/json")
			resp, err := client.HTTPClient.Do(req)
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			resp.Body.Close()
			load()
		})
	})
	btnDel := widget.NewButton("Delete", func() {
		if selectedID == 0 {
			dialog.ShowInformation("Select one", "Please select a plan to delete.", win)
			return
		}
		confirm := dialog.NewConfirm("Confirm Delete", "Are you sure?", func(b bool) {
			if !b {
				return
			}
			req, _ := http.NewRequest(http.MethodDelete, client.BaseURL+fmt.Sprintf("/plans/%d", selectedID), nil)
			resp, err := client.HTTPClient.Do(req)
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			resp.Body.Close()
			load()
		}, win)
		confirm.Show()
	})

	load()
	controls := container.NewHBox(btnRefresh, btnNew, btnEdit, btnDel)
	return container.NewBorder(controls, nil, nil, nil, list)
}

func InstallmentTab(win fyne.Window) fyne.CanvasObject {
	// Domain model for UI
	type Installment struct {
		ID         int64   `json:"id"`
		PlanID     int64   `json:"plan_id"`
		DueDate    string  `json:"due_date"`
		AmountDue  float64 `json:"amount_due"`
		AmountPaid float64 `json:"amount_paid"`
	}

	var items []Installment
	list := widget.NewList(
		func() int { return len(items) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(i widget.ListItemID, obj fyne.CanvasObject) {
			inst := items[i]
			obj.(*widget.Label).SetText(
				fmt.Sprintf("#%d – Plan:%d – %s – Due:$%.2f Paid:$%.2f",
					inst.ID, inst.PlanID, inst.DueDate, inst.AmountDue, inst.AmountPaid),
			)
		},
	)

	// Load data from backend
	load := func() {
		resp, err := client.HTTPClient.Get(client.BaseURL + "/installments")
		if err != nil {
			dialog.ShowError(err, win)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			dialog.ShowError(fmt.Errorf("server error: %s", resp.Status), win)
			return
		}
		var data []Installment
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			dialog.ShowError(err, win)
			return
		}
		items = data
		list.Refresh()
	}

	var selectedID int64
	list.OnSelected = func(id widget.ListItemID) {
		selectedID = items[id].ID
	}

	// Form helper to create or edit
	showForm := func(title string, inst *Installment, save func(Installment)) {
		planEntry := widget.NewEntry()
		planEntry.SetText(strconv.FormatInt(inst.PlanID, 10))
		dueEntry := widget.NewEntry()
		dueEntry.SetText(inst.DueDate)
		paidEntry := widget.NewEntry()
		paidEntry.SetText(fmt.Sprintf("%.2f", inst.AmountPaid))

		dlg := dialog.NewForm(
			title,
			"Save",
			"Cancel",
			[]*widget.FormItem{
				{Text: "Plan ID", Widget: planEntry},
				{Text: "Due Date (YYYY-MM-DD)", Widget: dueEntry},
				{Text: "Amount Paid", Widget: paidEntry},
			},
			func(bOK bool) {
				if !bOK {
					dlg.Hide()
					return
				}
				pid, _ := strconv.ParseInt(planEntry.Text, 10, 64)
				paid, _ := strconv.ParseFloat(paidEntry.Text, 64)
				n := Installment{
					ID:         inst.ID,
					PlanID:     pid,
					DueDate:    dueEntry.Text,
					AmountDue:  inst.AmountDue,
					AmountPaid: paid,
				}
				save(n)
				dlg.Hide()
			}, win)
		dlg.Show()
	}

	// CRUD buttons
	btnRefresh := widget.NewButton("Refresh", func() { load() })
	btnNew := widget.NewButton("New", func() {
		showForm("New Installment", &Installment{}, func(n Installment) {
			buf, _ := json.Marshal(n)
			resp, err := client.HTTPClient.Post(client.BaseURL+"/installments", "application/json", bytes.NewReader(buf))
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			resp.Body.Close()
			load()
		})
	})
	btnEdit := widget.NewButton("Edit", func() {
		if selectedID == 0 {
			dialog.ShowInformation("Select one", "Please select an installment to edit.", win)
			return
		}
		var current Installment
		for _, inst := range items {
			if inst.ID == selectedID {
				current = inst
				break
			}
		}
		showForm("Edit Installment", &current, func(n Installment) {
			buf, _ := json.Marshal(n)
			req, _ := http.NewRequest(http.MethodPut, client.BaseURL+fmt.Sprintf("/installments/%d", n.ID), bytes.NewReader(buf))
			req.Header.Set("Content-Type", "application/json")
			resp, err := client.HTTPClient.Do(req)
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			resp.Body.Close()
			load()
		})
	})
	btnDel := widget.NewButton("Delete", func() {
		if selectedID == 0 {
			dialog.ShowInformation("Select one", "Please select an installment to delete.", win)
			return
		}
		confirm := dialog.NewConfirm("Confirm Delete", "Are you sure?", func(b bool) {
			if !b {
				return
			}
			req, _ := http.NewRequest(http.MethodDelete, client.BaseURL+fmt.Sprintf("/installments/%d", selectedID), nil)
			resp, err := client.HTTPClient.Do(req)
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			resp.Body.Close()
			load()
		}, win)
		confirm.Show()
	})

	load()
	controls := container.NewHBox(btnRefresh, btnNew, btnEdit, btnDel)
	return container.NewBorder(controls, nil, nil, nil, list)
}

// PaymentTab builds the Payments screen with list and CRUD actions.
func PaymentTab(win fyne.Window) fyne.CanvasObject {
	// Domain model for UI
	type Payment struct {
		ID            int64   `json:"id"`
		InstallmentID int64   `json:"installment_id"`
		PaymentDate   string  `json:"payment_date"`
		Amount        float64 `json:"amount"`
	}

	var items []Payment
	list := widget.NewList(
		func() int { return len(items) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(i widget.ListItemID, obj fyne.CanvasObject) {
			p := items[i]
			obj.(*widget.Label).SetText(
				fmt.Sprintf("#%d – Inst:%d – %s – $%.2f", p.ID, p.InstallmentID, p.PaymentDate, p.Amount),
			)
		},
	)

	// Load data from backend
	load := func() {
		resp, err := client.HTTPClient.Get(client.BaseURL + "/payments")
		if err != nil {
			dialog.ShowError(err, win)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			dialog.ShowError(fmt.Errorf("server error: %s", resp.Status), win)
			return
		}
		var data []Payment
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			dialog.ShowError(err, win)
			return
		}
		items = data
		list.Refresh()
	}

	var selectedID int64
	list.OnSelected = func(id widget.ListItemID) {
		selectedID = items[id].ID
	}

	// Form helper to create or edit
	showForm := func(title string, p *Payment, save func(Payment)) {
		instEntry := widget.NewEntry()
		instEntry.SetText(strconv.FormatInt(p.InstallmentID, 10))
		dateEntry := widget.NewEntry()
		dateEntry.SetText(p.PaymentDate)
		amountEntry := widget.NewEntry()
		amountEntry.SetText(fmt.Sprintf("%.2f", p.Amount))

		dlg := dialog.NewForm(
			title,
			"Save",
			"Cancel",
			[]*widget.FormItem{
				{Text: "Installment ID", Widget: instEntry},
				{Text: "Payment Date (YYYY-MM-DD)", Widget: dateEntry},
				{Text: "Amount", Widget: amountEntry},
			},
			func(bOK bool) {
				if !bOK {
					dlg.Hide()
					return
				}
				iid, _ := strconv.ParseInt(instEntry.Text, 10, 64)
				amt, _ := strconv.ParseFloat(amountEntry.Text, 64)
				n := Payment{ID: p.ID, InstallmentID: iid, PaymentDate: dateEntry.Text, Amount: amt}
				save(n)
				dlg.Hide()
			}, win)
		dlg.Show()
	}

	// CRUD buttons
	btnRefresh := widget.NewButton("Refresh", func() { load() })

	btnNew := widget.NewButton("New", func() {
		showForm("New Payment", &Payment{}, func(n Payment) {
			buf, _ := json.Marshal(n)
			resp, err := client.HTTPClient.Post(client.BaseURL+"/payments", "application/json", bytes.NewReader(buf))
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			resp.Body.Close()
			load()
		})
	})

	btnEdit := widget.NewButton("Edit", func() {
		if selectedID == 0 {
			dialog.ShowInformation("Select one", "Please select a payment to edit.", win)
			return
		}
		var current Payment
		for _, pay := range items {
			if pay.ID == selectedID {
				current = pay
				break
			}
		}
		showForm("Edit Payment", &current, func(n Payment) {
			buf, _ := json.Marshal(n)
			req, _ := http.NewRequest(http.MethodPut, client.BaseURL+fmt.Sprintf("/payments/%d", n.ID), bytes.NewReader(buf))
			req.Header.Set("Content-Type", "application/json")
			resp, err := client.HTTPClient.Do(req)
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			resp.Body.Close()
			load()
		})
	})

	btnDel := widget.NewButton("Delete", func() {
		if selectedID == 0 {
			dialog.ShowInformation("Select one", "Please select a payment to delete.", win)
			return
		}
		confirm := dialog.NewConfirm("Confirm Delete", "Are you sure?", func(b bool) {
			if !b {
				return
			}
			req, _ := http.NewRequest(http.MethodDelete, client.BaseURL+fmt.Sprintf("/payments/%d", selectedID), nil)
			resp, err := client.HTTPClient.Do(req)
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			resp.Body.Close()
			load()
		}, win)
		confirm.Show()
	})

	load()
	controls := container.NewHBox(btnRefresh, btnNew, btnEdit, btnDel)
	return container.NewBorder(controls, nil, nil, nil, list)
}

// CommissionTab builds the Commissions screen with list and CRUD actions.
func CommissionTab(win fyne.Window) fyne.CanvasObject {
	// Domain model for UI
	type Commission struct {
		ID               int64   `json:"id"`
		BeneficiaryID    int64   `json:"beneficiary_id"`
		TransactionType  string  `json:"transaction_type"`
		TransactionID    int64   `json:"transaction_id"`
		RateOrAmount     float64 `json:"rate_or_amount"`
		CalculatedAmount float64 `json:"calculated_amount"`
		Memo             string  `json:"memo"`
	}

	var items []Commission
	list := widget.NewList(
		func() int { return len(items) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(i widget.ListItemID, obj fyne.CanvasObject) {
			c := items[i]
			obj.(*widget.Label).SetText(
				fmt.Sprintf("#%d – B:%d – %s %d – Rate: %.2f → Amt: %.2f",
					c.ID, c.BeneficiaryID, c.TransactionType, c.TransactionID,
					c.RateOrAmount, c.CalculatedAmount),
			)
		},
	)

	// Load data from backend
	load := func() {
		resp, err := client.HTTPClient.Get(client.BaseURL + "/commissions")
		if err != nil {
			dialog.ShowError(err, win)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			dialog.ShowError(fmt.Errorf("server error: %s", resp.Status), win)
			return
		}
		var data []Commission
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			dialog.ShowError(err, win)
			return
		}
		items = data
		list.Refresh()
	}

	var selectedID int64
	list.OnSelected = func(id widget.ListItemID) {
		selectedID = items[id].ID
	}

	// Form helper to create or edit
	showForm := func(title string, c *Commission, save func(Commission)) {
		benEntry := widget.NewEntry()
		benEntry.SetText(strconv.FormatInt(c.BeneficiaryID, 10))
		tranType := widget.NewSelect([]string{"sale", "letting", "introduction"}, nil)
		tranType.SetSelected(c.TransactionType)
		tranID := widget.NewEntry()
		tranID.SetText(strconv.FormatInt(c.TransactionID, 10))
		rateEntry := widget.NewEntry()
		rateEntry.SetText(fmt.Sprintf("%.2f", c.RateOrAmount))
		calcEntry := widget.NewEntry()
		calcEntry.SetText(fmt.Sprintf("%.2f", c.CalculatedAmount))
		memoEntry := widget.NewEntry()
		memoEntry.SetText(c.Memo)

		dlg := dialog.NewForm(
			title,
			"Save",
			"Cancel",
			[]*widget.FormItem{
				{Text: "Beneficiary ID", Widget: benEntry},
				{Text: "Transaction Type", Widget: tranType},
				{Text: "Transaction ID", Widget: tranID},
				{Text: "RateOrAmount", Widget: rateEntry},
				{Text: "CalculatedAmount", Widget: calcEntry},
				{Text: "Memo", Widget: memoEntry},
			},
			func(bOK bool) {
				if !bOK {
					dlg.Hide()
					return
				}
				bid, _ := strconv.ParseInt(benEntry.Text, 10, 64)
				tid, _ := strconv.ParseInt(tranID.Text, 10, 64)
				rate, _ := strconv.ParseFloat(rateEntry.Text, 64)
				calc, _ := strconv.ParseFloat(calcEntry.Text, 64)
				n := Commission{ID: c.ID, BeneficiaryID: bid,
					TransactionType:  tranType.Selected,
					TransactionID:    tid,
					RateOrAmount:     rate,
					CalculatedAmount: calc,
					Memo:             memoEntry.Text,
				}
				save(n)
				dlg.Hide()
			}, win)
		dlg.Show()
	}

	// CRUD buttons
	btnRefresh := widget.NewButton("Refresh", func() { load() })
	btnNew := widget.NewButton("New", func() {
		showForm("New Commission", &Commission{}, func(n Commission) {
			buf, _ := json.Marshal(n)
			resp, err := client.HTTPClient.Post(client.BaseURL+"/commissions", "application/json", bytes.NewReader(buf))
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			resp.Body.Close()
			load()
		})
	})
	btnEdit := widget.NewButton("Edit", func() {
		if selectedID == 0 {
			dialog.ShowInformation("Select one", "Please select a commission to edit.", win)
			return
		}
		var current Commission
		for _, c := range items {
			if c.ID == selectedID {
				current = c
				break
			}
		}
		showForm("Edit Commission", &current, func(n Commission) {
			buf, _ := json.Marshal(n)
			req, _ := http.NewRequest(http.MethodPut, client.BaseURL+fmt.Sprintf("/commissions/%d", n.ID), bytes.NewReader(buf))
			req.Header.Set("Content-Type", "application/json")
			resp, err := client.HTTPClient.Do(req)
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			resp.Body.Close()
			load()
		})
	})
	btnDel := widget.NewButton("Delete", func() {
		if selectedID == 0 {
			dialog.ShowInformation("Select one", "Please select a commission to delete.", win)
			return
		}
		confirm := dialog.NewConfirm("Confirm Delete", "Are you sure?", func(b bool) {
			if !b {
				return
			}
			req, _ := http.NewRequest(http.MethodDelete, client.BaseURL+fmt.Sprintf("/commissions/%d", selectedID), nil)
			resp, err := client.HTTPClient.Do(req)
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			resp.Body.Close()
			load()
		}, win)
		confirm.Show()
	})

	load()
	controls := container.NewHBox(btnRefresh, btnNew, btnEdit, btnDel)
	return container.NewBorder(controls, nil, nil, nil, list)
}

// SalesTab builds the Sales screen with list and CRUD actions.
func SalesTab(win fyne.Window) fyne.CanvasObject {
	// Domain model for UI
	type Sale struct {
		ID         int64   `json:"id"`
		PropertyID int64   `json:"property_id"`
		SaleDate   string  `json:"sale_date"`
		SalePrice  float64 `json:"sale_price"`
	}

	var items []Sale
	list := widget.NewList(
		func() int { return len(items) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(i widget.ListItemID, obj fyne.CanvasObject) {
			s := items[i]
			obj.(*widget.Label).SetText(
				fmt.Sprintf("#%d – Prop:%d – %s – $%.2f", s.ID, s.PropertyID, s.SaleDate, s.SalePrice),
			)
		},
	)

	// Load data from backend
	load := func() {
		resp, err := client.HTTPClient.Get(client.BaseURL + "/sales")
		if err != nil {
			dialog.ShowError(err, win)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			dialog.ShowError(fmt.Errorf("server error: %s", resp.Status), win)
			return
		}
		var data []Sale
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			dialog.ShowError(err, win)
			return
		}
		items = data
		list.Refresh()
	}

	var selectedID int64
	list.OnSelected = func(id widget.ListItemID) {
		selectedID = items[id].ID
	}

	// Form helper to create or edit
	showForm := func(title string, s *Sale, save func(Sale)) {
		propEntry := widget.NewEntry()
		propEntry.SetText(strconv.FormatInt(s.PropertyID, 10))
		saleDate := widget.NewEntry()
		saleDate.SetText(s.SaleDate)
		salePrice := widget.NewEntry()
		salePrice.SetText(fmt.Sprintf("%.2f", s.SalePrice))

		dlg := dialog.NewForm(
			title,
			"Save",
			"Cancel",
			[]*widget.FormItem{
				{Text: "Property ID", Widget: propEntry},
				{Text: "Sale Date (YYYY-MM-DD)", Widget: saleDate},
				{Text: "Sale Price", Widget: salePrice},
			},
			func(bOK bool) {
				if !bOK {
					dlg.Hide()
					return
				}
				pid, _ := strconv.ParseInt(propEntry.Text, 10, 64)
				price, _ := strconv.ParseFloat(salePrice.Text, 64)
				n := Sale{ID: s.ID, PropertyID: pid, SaleDate: saleDate.Text, SalePrice: price}
				save(n)
				dlg.Hide()
			}, win)
		dlg.Show()
	}

	// CRUD buttons
	btnRefresh := widget.NewButton("Refresh", func() { load() })

	btnNew := widget.NewButton("New", func() {
		showForm("New Sale", &Sale{}, func(n Sale) {
			buf, _ := json.Marshal(n)
			resp, err := client.HTTPClient.Post(client.BaseURL+"/sales", "application/json", bytes.NewReader(buf))
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			resp.Body.Close()
			load()
		})
	})

	btnEdit := widget.NewButton("Edit", func() {
		if selectedID == 0 {
			dialog.ShowInformation("Select one", "Please select a sale to edit.", win)
			return
		}
		var current Sale
		for _, s := range items {
			if s.ID == selectedID {
				current = s
				break
			}
		}
		showForm("Edit Sale", &current, func(n Sale) {
			buf, _ := json.Marshal(n)
			req, _ := http.NewRequest(http.MethodPut, client.BaseURL+fmt.Sprintf("/sales/%d", n.ID), bytes.NewReader(buf))
			req.Header.Set("Content-Type", "application/json")
			resp, err := client.HTTPClient.Do(req)
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			resp.Body.Close()
			load()
		})
	})

	btnDel := widget.NewButton("Delete", func() {
		if selectedID == 0 {
			dialog.ShowInformation("Select one", "Please select a sale to delete.", win)
			return
		}
		confirm := dialog.NewConfirm("Confirm Delete", "Are you sure?", func(b bool) {
			if !b {
				return
			}
			req, _ := http.NewRequest(http.MethodDelete, client.BaseURL+fmt.Sprintf("/sales/%d", selectedID), nil)
			resp, err := client.HTTPClient.Do(req)
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			resp.Body.Close()
			load()
		}, win)
		confirm.Show()
	})

	load()
	controls := container.NewHBox(btnRefresh, btnNew, btnEdit, btnDel)
	return container.NewBorder(controls, nil, nil, nil, list)
}

// LettingsTab builds the Lettings screen with list and CRUD actions.
func LettingsTab(win fyne.Window) fyne.CanvasObject {
	// Domain model for UI
	type Letting struct {
		ID         int64   `json:"id"`
		PropertyID int64   `json:"property_id"`
		StartDate  string  `json:"start_date"`
		EndDate    string  `json:"end_date"`
		RentAmount float64 `json:"rent_amount"`
	}

	var items []Letting
	list := widget.NewList(
		func() int { return len(items) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(i widget.ListItemID, obj fyne.CanvasObject) {
			lt := items[i]
			obj.(*widget.Label).SetText(
				fmt.Sprintf("#%d – Prop:%d – %s → %s – $%.2f",
					lt.ID, lt.PropertyID, lt.StartDate, lt.EndDate, lt.RentAmount),
			)
		},
	)

	// Load data from backend
	load := func() {
		resp, err := client.HTTPClient.Get(client.BaseURL + "/lettings")
		if err != nil {
			dialog.ShowError(err, win)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			dialog.ShowError(fmt.Errorf("server error: %s", resp.Status), win)
			return
		}
		var data []Letting
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			dialog.ShowError(err, win)
			return
		}
		items = data
		list.Refresh()
	}

	var selectedID int64
	list.OnSelected = func(id widget.ListItemID) {
		selectedID = items[id].ID
	}

	// Form helper to create or edit
	showForm := func(title string, lt *Letting, save func(Letting)) {
		propEntry := widget.NewEntry()
		propEntry.SetText(strconv.FormatInt(lt.PropertyID, 10))
		startEntry := widget.NewEntry()
		startEntry.SetText(lt.StartDate)
		endEntry := widget.NewEntry()
		endEntry.SetText(lt.EndDate)
		rentEntry := widget.NewEntry()
		rentEntry.SetText(fmt.Sprintf("%.2f", lt.RentAmount))

		formItems := []*widget.FormItem{
			{Text: "Property ID", Widget: propEntry},
			{Text: "Start Date (YYYY-MM-DD)", Widget: startEntry},
			{Text: "End Date (YYYY-MM-DD or blank)", Widget: endEntry},
			{Text: "Rent Amount", Widget: rentEntry},
		}

		dlg := dialog.NewForm(
			title,
			"Save",
			"Cancel",
			formItems,
			func(bOK bool) {
				if !bOK {
					dlg.Hide()
					return
				}
				// parse entries
				pid, _ := strconv.ParseInt(propEntry.Text, 10, 64)
				rentAmt, _ := strconv.ParseFloat(rentEntry.Text, 64)
				n := Letting{ID: lt.ID, PropertyID: pid,
					StartDate:  startEntry.Text,
					EndDate:    endEntry.Text,
					RentAmount: rentAmt,
				}
				save(n)
				dlg.Hide()
			}, win)
		dlg.Show()
	}

	// CRUD buttons
	btnRefresh := widget.NewButton("Refresh", func() {
		load()
	})

	btnNew := widget.NewButton("New", func() {
		showForm("New Letting", &Letting{}, func(n Letting) {
			// POST JSON
			buf, _ := json.Marshal(n)
			resp, err := client.HTTPClient.Post(
				client.BaseURL+"/lettings",
				"application/json",
				bytes.NewReader(buf),
			)
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			resp.Body.Close()
			load()
		})
	})

	btnEdit := widget.NewButton("Edit", func() {
		if selectedID == 0 {
			dialog.ShowInformation("Select one", "Please select a letting to edit.", win)
			return
		}
		// find selected
		var current Letting
		for _, lt := range items {
			if lt.ID == selectedID {
				current = lt
				break
			}
		}
		showForm("Edit Letting", &current, func(n Letting) {
			// PUT JSON
			buf, _ := json.Marshal(n)
			req, _ := http.NewRequest(http.MethodPut,
				client.BaseURL+fmt.Sprintf("/lettings/%d", n.ID),
				bytes.NewReader(buf),
			)
			req.Header.Set("Content-Type", "application/json")
			resp, err := client.HTTPClient.Do(req)
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			resp.Body.Close()
			load()
		})
	})

	btnDel := widget.NewButton("Delete", func() {
		if selectedID == 0 {
			dialog.ShowInformation("Select one", "Please select a letting to delete.", win)
			return
		}
		confirm := dialog.NewConfirm("Confirm Delete",
			"Are you sure you want to delete this letting?",
			func(b bool) {
				if !b {
					return
				}
				req, _ := http.NewRequest(http.MethodDelete,
					client.BaseURL+fmt.Sprintf("/lettings/%d", selectedID), nil)
				resp, err := client.HTTPClient.Do(req)
				if err != nil {
					dialog.ShowError(err, win)
					return
				}
				resp.Body.Close()
				load()
			}, win)
		confirm.Show()
	})

	load()

	// Layout: buttons up top, list below
	controls := container.NewHBox(btnRefresh, btnNew, btnEdit, btnDel)
	return container.NewBorder(controls, nil, nil, nil, list)
}

// ReportTab builds the Reporting screen with a dropdown and output area.
func ReportTab(win fyne.Window) fyne.CanvasObject {
	// Map of user-friendly names to API endpoints
	reports := map[string]string{
		"Commissions by Beneficiary": "/reports/commissions/beneficiary",
		"Outstanding Installments":   "/reports/installments/outstanding",
		"Monthly Sales Volume":       "/reports/sales/monthly",
		"Active Rent Roll":           "/reports/lettings/rentroll",
		"Top Property Payments":      "/reports/properties/top-payments",
	}

	// Dropdown for selecting report
	reportSelect := widget.NewSelect(getReportKeys(reports), nil)

	// Button to run the selected report
	runBtn := widget.NewButton("Run Report", nil)

	// Multi-line entry to display JSON results
	output := widget.NewMultiLineEntry()
	output.SetPlaceHolder("Report output will appear here...")
	output.SetReadOnly(true)

	runBtn.OnTapped = func() {
		choice := reportSelect.Selected
		if choice == "" {
			output.SetText("Please select a report to run.")
			return
		}
		endpoint := client.BaseURL + reports[choice]
		resp, err := client.HTTPClient.Get(endpoint)
		if err != nil {
			output.SetText(fmt.Sprintf("Request error: %v", err))
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			output.SetText(fmt.Sprintf("API error: %s", resp.Status))
			return
		}
		var data interface{}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			output.SetText(fmt.Sprintf("Decode error: %v", err))
			return
		}
		pretty, _ := json.MarshalIndent(data, "", "  ")
		output.SetText(string(pretty))
	}

	// Layout: Top controls and output area
	controls := container.NewHBox(
		widget.NewLabelWithStyle("Select Report:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		reportSelect,
		runBtn,
	)

	return container.NewBorder(
		controls, // north
		nil,      // south
		nil,      // west
		nil,      // east
		output,   // center
	)
}

// getReportKeys returns sorted keys of the reports map
func getReportKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

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
