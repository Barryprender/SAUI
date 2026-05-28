package banking

import "fmt"

// Recipient is a preset transfer target.
type Recipient struct {
	ID         string
	Name       string
	Initials   string
	Bank       string
	ColorClass string
}

var allRecipients = []Recipient{
	{ID: "sarah-chen", Name: "Sarah Chen", Initials: "SC", Bank: "Revolut", ColorClass: "avatar--violet"},
	{ID: "marcus-webb", Name: "Marcus Webb", Initials: "MW", Bank: "Monzo", ColorClass: "avatar--blue"},
	{ID: "elena-sousa", Name: "Elena Sousa", Initials: "ES", Bank: "N26", ColorClass: "avatar--teal"},
	{ID: "james-obrien", Name: "James O'Brien", Initials: "JO", Bank: "Revolut", ColorClass: "avatar--rose"},
}

func Recipients() []Recipient { return allRecipients }

func RecipientByID(id string) (Recipient, bool) {
	for _, r := range allRecipients {
		if r.ID == id {
			return r, true
		}
	}
	return Recipient{}, false
}

func FormatBalance(cents int64) string {
	negative := cents < 0
	if negative {
		cents = -cents
	}
	euros := cents / 100
	frac := cents % 100
	var s string
	if euros >= 1000 {
		s = fmt.Sprintf("%d,%03d.%02d", euros/1000, euros%1000, frac)
	} else {
		s = fmt.Sprintf("%d.%02d", euros, frac)
	}
	if negative {
		return "-€" + s
	}
	return "€" + s
}
