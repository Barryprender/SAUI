package food

import (
	"context"
	"errors"

	"saui/statestore"
)

// AddToCart adds one unit of the given item to the session cart.
type AddToCart struct {
	ItemID string
}

func (a AddToCart) Type() string { return CartItemAdded }

func (a AddToCart) Validate(ctx context.Context, s *statestore.Store, sessionID string) error {
	if a.ItemID == "" {
		return errors.New("item_id is required")
	}
	item, ok, err := s.GetFoodItem(ctx, a.ItemID)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("item not found")
	}
	if !item.Available {
		return errors.New(item.Name + " is not available")
	}
	return nil
}

func (a AddToCart) Apply(ctx context.Context, s *statestore.Store, sessionID string) (statestore.Event, error) {
	item, _, err := s.GetFoodItem(ctx, a.ItemID)
	if err != nil {
		return statestore.Event{}, err
	}
	return s.AppendEvent(ctx, sessionID, CartItemAdded, map[string]any{
		"item_id":     item.ID,
		"name":        item.Name,
		"price_cents": item.PriceCents,
		"category":    item.Category,
	})
}

// RemoveFromCart removes an item entirely from the session cart.
type RemoveFromCart struct {
	ItemID string
}

func (a RemoveFromCart) Type() string { return CartItemRemoved }

func (a RemoveFromCart) Validate(ctx context.Context, s *statestore.Store, sessionID string) error {
	if a.ItemID == "" {
		return errors.New("item_id is required")
	}
	events, err := s.EventsBySession(ctx, sessionID)
	if err != nil {
		return err
	}
	if _, ok := cartFromEvents(events)[a.ItemID]; !ok {
		return errors.New("item not in cart")
	}
	return nil
}

func (a RemoveFromCart) Apply(ctx context.Context, s *statestore.Store, sessionID string) (statestore.Event, error) {
	return s.AppendEvent(ctx, sessionID, CartItemRemoved, map[string]any{
		"item_id": a.ItemID,
	})
}

// PlaceOrder validates the cart and commits it as an immutable order event.
// The server computes the total — the client never submits a price.
type PlaceOrder struct{}

func (a PlaceOrder) Type() string { return OrderPlaced }

func (a PlaceOrder) Validate(ctx context.Context, s *statestore.Store, sessionID string) error {
	events, err := s.EventsBySession(ctx, sessionID)
	if err != nil {
		return err
	}
	cart := cartFromEvents(events)
	if len(cart) == 0 {
		return errors.New("cart is empty")
	}
	for _, entry := range cart {
		item, ok, err := s.GetFoodItem(ctx, entry.ItemID)
		if err != nil {
			return err
		}
		if !ok || !item.Available {
			return errors.New(entry.Name + " is no longer available")
		}
	}
	return nil
}

func (a PlaceOrder) Apply(ctx context.Context, s *statestore.Store, sessionID string) (statestore.Event, error) {
	events, err := s.EventsBySession(ctx, sessionID)
	if err != nil {
		return statestore.Event{}, err
	}
	cart := cartFromEvents(events)

	var items []OrderItem
	var total int64
	for _, entry := range cart {
		items = append(items, OrderItem{
			ItemID:     entry.ItemID,
			Name:       entry.Name,
			PriceCents: entry.PriceCents,
			Quantity:   entry.Quantity,
		})
		total += entry.PriceCents * int64(entry.Quantity)
	}
	return s.AppendEvent(ctx, sessionID, OrderPlaced, map[string]any{
		"items":       items,
		"total_cents": total,
	})
}
