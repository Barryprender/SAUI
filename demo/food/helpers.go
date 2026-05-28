package food

import "fmt"

func formatPrice(cents int64) string {
	return fmt.Sprintf("€%.2f", float64(cents)/100)
}

// CartItemQty returns the current quantity of itemID in the cart, 0 if absent.
func CartItemQty(cart CartProjection, itemID string) int {
	for _, item := range cart.Items {
		if item.ItemID == itemID {
			return item.Quantity
		}
	}
	return 0
}

// CartItemName returns the name of itemID in the cart, "" if absent.
func CartItemName(cart CartProjection, itemID string) string {
	for _, item := range cart.Items {
		if item.ItemID == itemID {
			return item.Name
		}
	}
	return ""
}
