package food

import (
	"context"
	"encoding/json"
	"time"

	"saui/statestore"
)

const (
	MenuProjectionName  = "food.menu"
	CartProjectionName  = "food.cart"
	OrderProjectionName = "food.order"
)

// MenuProjection is the read model for the restaurant menu, grouped by category.
type MenuProjection struct {
	Categories []MenuCategory
}

type MenuCategory struct {
	Name  string
	Items []statestore.FoodItem
}

func (p MenuProjection) Name() string { return MenuProjectionName }

// CartProjection is the server-computed cart for a session.
type CartItem struct {
	ItemID     string
	Name       string
	PriceCents int64
	Quantity   int
}

type CartProjection struct {
	Items      []CartItem
	TotalCents int64
	ItemCount  int
}

func (p CartProjection) Name() string { return CartProjectionName }

// OrderProjection is the most recent order placed by a session.
type OrderProjection struct {
	HasOrder   bool
	Items      []OrderItem
	TotalCents int64
	PlacedAt   time.Time
}

func (p OrderProjection) Name() string { return OrderProjectionName }

func init() {
	statestore.RegisterProjection(MenuProjectionName, buildMenu)
	statestore.RegisterProjection(CartProjectionName, buildCart)
	statestore.RegisterProjection(OrderProjectionName, buildOrder)
}

func buildMenu(ctx context.Context, s *statestore.Store, _ string) (statestore.Projection, error) {
	items, err := s.ListFoodItems(ctx)
	if err != nil {
		return nil, err
	}

	catMap := map[string]*MenuCategory{}
	var catOrder []string
	for _, item := range items {
		if _, ok := catMap[item.Category]; !ok {
			catMap[item.Category] = &MenuCategory{Name: item.Category}
			catOrder = append(catOrder, item.Category)
		}
		catMap[item.Category].Items = append(catMap[item.Category].Items, item)
	}

	proj := MenuProjection{}
	for _, name := range catOrder {
		proj.Categories = append(proj.Categories, *catMap[name])
	}
	return proj, nil
}

func buildCart(ctx context.Context, s *statestore.Store, sessionID string) (statestore.Projection, error) {
	events, err := s.EventsBySession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	proj := CartProjection{}
	for _, e := range cartFromEvents(events) {
		proj.Items = append(proj.Items, CartItem{
			ItemID:     e.ItemID,
			Name:       e.Name,
			PriceCents: e.PriceCents,
			Quantity:   e.Quantity,
		})
		proj.TotalCents += e.PriceCents * int64(e.Quantity)
		proj.ItemCount += e.Quantity
	}
	return proj, nil
}

func buildOrder(ctx context.Context, s *statestore.Store, sessionID string) (statestore.Projection, error) {
	events, err := s.EventsBySession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	proj := OrderProjection{}
	for _, ev := range events {
		if ev.Type != OrderPlaced {
			continue
		}
		var p struct {
			Items      []OrderItem `json:"items"`
			TotalCents int64       `json:"total_cents"`
		}
		if err := json.Unmarshal(ev.Payload, &p); err != nil {
			continue
		}
		proj.HasOrder = true
		proj.Items = p.Items
		proj.TotalCents = p.TotalCents
		proj.PlacedAt = ev.OccurredAt
	}
	return proj, nil
}
