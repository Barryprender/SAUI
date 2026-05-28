package banking

import (
	"encoding/json"
	"time"

	"saui/statestore"
)

const (
	AccountFunded = "banking.account.funded"
	TransferSent  = "banking.transfer.sent"
)

type txEntry struct {
	EventID             int64
	RecipientID         string
	RecipientName       string
	RecipientInitials   string
	RecipientColorClass string
	AmountCents         int64
	OccurredAt          time.Time
}

func balanceFromEvents(events []statestore.Event) (cents int64, funded bool) {
	for _, ev := range events {
		switch ev.Type {
		case AccountFunded:
			var p struct {
				AmountCents int64 `json:"amount_cents"`
			}
			if json.Unmarshal(ev.Payload, &p) == nil {
				cents += p.AmountCents
				funded = true
			}
		case TransferSent:
			var p struct {
				AmountCents int64 `json:"amount_cents"`
			}
			if json.Unmarshal(ev.Payload, &p) == nil {
				cents -= p.AmountCents
			}
		}
	}
	return
}

func txFromEvents(events []statestore.Event) []txEntry {
	var out []txEntry
	for _, ev := range events {
		if ev.Type != TransferSent {
			continue
		}
		var p struct {
			RecipientID   string `json:"recipient_id"`
			RecipientName string `json:"recipient_name"`
			AmountCents   int64  `json:"amount_cents"`
		}
		if json.Unmarshal(ev.Payload, &p) != nil {
			continue
		}
		r, _ := RecipientByID(p.RecipientID)
		initials := r.Initials
		if initials == "" {
			initials = nameInitials(p.RecipientName)
		}
		colorClass := r.ColorClass
		if colorClass == "" {
			colorClass = "avatar--blue"
		}
		out = append(out, txEntry{
			EventID:             ev.ID,
			RecipientID:         p.RecipientID,
			RecipientName:       p.RecipientName,
			RecipientInitials:   initials,
			RecipientColorClass: colorClass,
			AmountCents:         p.AmountCents,
			OccurredAt:          ev.OccurredAt,
		})
	}
	return out
}

func nameInitials(name string) string {
	words := make([]byte, 0, 2)
	prev := ' '
	for i := 0; i < len(name) && len(words) < 2; i++ {
		c := name[i]
		if prev == ' ' && c != ' ' {
			words = append(words, c)
		}
		prev = rune(c)
	}
	return string(words)
}
