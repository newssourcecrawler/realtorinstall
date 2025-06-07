package ui

import "fyne.io/fyne/v2"

// CustomTheme implements fyne.Theme if you want a branded look.
type CustomTheme struct{}

// … implement Theme methods …

func PropertyTab() fyne.CanvasObject    { /* build your form/table */ }
func BuyerTab() fyne.CanvasObject       { /* … */ }
func PlanTab() fyne.CanvasObject        { /* … */ }
func InstallmentTab() fyne.CanvasObject { /* … */ }
func PaymentTab() fyne.CanvasObject     { /* … */ }
func CommissionTab() fyne.CanvasObject  { /* … */ }
func SalesTab() fyne.CanvasObject       { /* … */ }
func LettingsTab() fyne.CanvasObject    { /* … */ }
func ReportTab() fyne.CanvasObject      { /* … */ }
