package statestore

import (
	"context"
	"database/sql"
	"strings"
)

type FoodItem struct {
	ID          string
	Name        string
	Description string
	PriceCents  int64
	Category    string
	Available   bool
}

func (s *Store) migrateFoodDemo() error {
	// Create table for new installs.
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS demo_food_items (
			id          TEXT    PRIMARY KEY,
			name        TEXT    NOT NULL,
			description TEXT    NOT NULL DEFAULT '',
			price_cents INTEGER NOT NULL,
			category    TEXT    NOT NULL,
			sort_order  INTEGER NOT NULL DEFAULT 0,
			available   INTEGER NOT NULL DEFAULT 1
		)`)
	if err != nil {
		return err
	}

	// Add sort_order column for existing installs that predate it.
	_, err = s.db.Exec(`ALTER TABLE demo_food_items ADD COLUMN sort_order INTEGER NOT NULL DEFAULT 0`)
	if err != nil && !strings.Contains(err.Error(), "duplicate column name") {
		return err
	}

	// Seed reference data; UPDATE handles sort_order on rows that already exist.
	rows := []struct {
		id, name, desc, cat string
		price               int
		order               int
	}{
		{"bruschetta", "Bruschetta", "Grilled bread, tomato, garlic, basil", "Starters", 850, 1},
		{"burrata", "Burrata", "Fresh burrata, heirloom tomatoes, olive oil", "Starters", 1200, 2},
		{"margherita", "Margherita", "San Marzano tomato, fior di latte, fresh basil", "Mains", 1400, 3},
		{"carbonara", "Pasta Carbonara", "Rigatoni, guanciale, pecorino, egg yolk", "Mains", 1600, 4},
		{"risotto", "Risotto ai Funghi", "Wild mushroom risotto, parmigiano, truffle oil", "Mains", 1800, 5},
		{"tiramisu", "Tiramisù", "Mascarpone, espresso savoiardi, cocoa", "Desserts", 700, 6},
		{"water", "Sparkling Water", "750ml San Pellegrino", "Drinks", 300, 7},
		{"wine", "House Wine", "Glass of house red or white", "Drinks", 900, 8},
	}
	for _, r := range rows {
		_, err := s.db.Exec(`
			INSERT INTO demo_food_items (id, name, description, price_cents, category, sort_order, available)
			VALUES (?, ?, ?, ?, ?, ?, 1)
			ON CONFLICT(id) DO UPDATE SET sort_order = excluded.sort_order`,
			r.id, r.name, r.desc, r.price, r.cat, r.order)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) ListFoodItems(ctx context.Context) ([]FoodItem, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, name, description, price_cents, category, available
		 FROM demo_food_items ORDER BY sort_order`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []FoodItem
	for rows.Next() {
		var item FoodItem
		var available int
		if err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.PriceCents, &item.Category, &available); err != nil {
			return nil, err
		}
		item.Available = available == 1
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *Store) GetFoodItem(ctx context.Context, id string) (FoodItem, bool, error) {
	var item FoodItem
	var available int
	err := s.db.QueryRowContext(ctx,
		`SELECT id, name, description, price_cents, category, available
		 FROM demo_food_items WHERE id = ?`, id,
	).Scan(&item.ID, &item.Name, &item.Description, &item.PriceCents, &item.Category, &available)
	if err == sql.ErrNoRows {
		return FoodItem{}, false, nil
	}
	if err != nil {
		return FoodItem{}, false, err
	}
	item.Available = available == 1
	return item, true, nil
}
