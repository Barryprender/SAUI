package food

import (
	"encoding/json"

	"saui/statestore"
)

const (
	CartItemAdded   = "food.cart.item_added"
	CartItemRemoved = "food.cart.item_removed"
	OrderPlaced     = "food.order.placed"
)

// OrderItem is the line-item record stored in an OrderPlaced event payload.
type OrderItem struct {
	ItemID     string `json:"item_id"`
	Name       string `json:"name"`
	PriceCents int64  `json:"price_cents"`
	Quantity   int    `json:"quantity"`
}

type cartEntry struct {
	ItemID     string
	Name       string
	PriceCents int64
	Quantity   int
}

// cartFromEvents replays food cart events into a map keyed by item ID.
func cartFromEvents(events []statestore.Event) map[string]*cartEntry {
	cart := map[string]*cartEntry{}
	for _, ev := range events {
		switch ev.Type {
		case CartItemAdded:
			var p struct {
				ItemID     string `json:"item_id"`
				Name       string `json:"name"`
				PriceCents int64  `json:"price_cents"`
			}
			if err := json.Unmarshal(ev.Payload, &p); err != nil {
				continue
			}
			if e, ok := cart[p.ItemID]; ok {
				e.Quantity++
			} else {
				cart[p.ItemID] = &cartEntry{
					ItemID:     p.ItemID,
					Name:       p.Name,
					PriceCents: p.PriceCents,
					Quantity:   1,
				}
			}
		case CartItemRemoved:
			var p struct {
				ItemID string `json:"item_id"`
			}
			if err := json.Unmarshal(ev.Payload, &p); err != nil {
				continue
			}
			delete(cart, p.ItemID)
		case OrderPlaced:
			cart = map[string]*cartEntry{}
		}
	}
	return cart
}
