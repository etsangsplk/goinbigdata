package main

import (
	"fmt"
	"html/template"
	"os"
	"time"
)

type Account struct {
	FirstName string
	LastName  string
}

type Purchase struct {
	Date          time.Time
	Description   string
	AmountInCents int
}

type Statement struct {
	FromDate  time.Time
	ToDate    time.Time
	Account   Account
	Purchases []Purchase
}
type PurchaseItem struct {
	Name        string
	Description string
}

type PurchaseItems struct {
	Items []PurchaseItem
}

func main() {
	fmap := template.FuncMap{
		"formatAsDollars": formatAsDollars,
		"formatAsDate":    formatAsDate,
		"urgentNote":      urgentNote,
	}
	t := template.Must(template.New("email.tmpl").Funcs(fmap).ParseFiles("email.tmpl"))
	err := t.Execute(os.Stdout, createMockStatement())
	if err != nil {
		panic(err)
	}

	gmap := template.FuncMap{
		"Name":        getName,
		"Description": getDescription,
	}
	m := template.Must(template.New("map.tmpl").Funcs(gmap).ParseFiles("map.tmpl"))
	err = m.Execute(os.Stdout, createPurchaseItems())
	if err != nil {
		panic(err)
	}
}

func formatAsDollars(valueInCents int) (string, error) {
	dollars := valueInCents / 100
	cents := valueInCents % 100
	return fmt.Sprintf("$%d.%2d", dollars, cents), nil
}

func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d/%d/%d", day, month, year)
}

func urgentNote(acc Account) string {
	return fmt.Sprintf("You have earned 100 VIP points that can be used for purchases")
}

func createPurchaseItems() PurchaseItems {
	return PurchaseItems{
		Items: []PurchaseItem{
			{
				Name:        "item1",
				Description: "description1",
			},
			{
				Name:        "item2",
				Description: "description2",
			},
		},
	}
}

func getName(p PurchaseItem) string {
	return p.Name
}

func getDescription(p PurchaseItem) string {
	return p.Description
}

func createMockStatement() Statement {
	return Statement{
		FromDate: time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC),
		ToDate:   time.Date(2016, 2, 1, 0, 0, 0, 0, time.UTC),
		Account: Account{
			FirstName: "John",
			LastName:  "Dow",
		},
		Purchases: []Purchase{
			{
				Date:          time.Date(2016, 1, 3, 0, 0, 0, 0, time.UTC),
				Description:   "Shovel",
				AmountInCents: 2326,
			},
			{
				Date:          time.Date(2016, 1, 8, 0, 0, 0, 0, time.UTC),
				Description:   "Staple remover",
				AmountInCents: 5432,
			},
		},
	}
}
