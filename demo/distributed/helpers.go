package distributed

// EventMeta maps an event type to its display properties.
type EventMeta struct {
	ServiceName string
	ColorClass  string
	Label       string
}

var eventMetaMap = map[string]EventMeta{
	OrderCreated:      {"Order Service", "log--order", "Order created"},
	InventoryReserved: {"Inventory Service", "log--inventory", "Inventory reserved"},
	PaymentCaptured:   {"Payment Service", "log--payment", "Payment captured"},
	PaymentDeclined:   {"Payment Service", "log--payment", "Payment declined"},
	OrderFulfilled:    {"Order Service", "log--order", "Order fulfilled"},
	OrderCancelled:    {"Order Service", "log--order", "Order cancelled"},
}

func eventMetaFor(eventType string) (EventMeta, bool) {
	m, ok := eventMetaMap[eventType]
	return m, ok
}

func nodeIndexForEvent(eventType string) int {
	switch eventType {
	case OrderCreated:
		return 0
	case InventoryReserved:
		return 1
	case PaymentCaptured, PaymentDeclined:
		return 2
	case OrderFulfilled, OrderCancelled:
		return 3
	}
	return -1
}

func freshNodes() []PipelineNode {
	return []PipelineNode{
		{Index: 0, Label: "Order", ServiceName: "Order Service", ColorClass: "node--order", Status: StagePending},
		{Index: 1, Label: "Inventory", ServiceName: "Inventory Service", ColorClass: "node--inventory", Status: StagePending},
		{Index: 2, Label: "Payment", ServiceName: "Payment Service", ColorClass: "node--payment", Status: StagePending},
		{Index: 3, Label: "Fulfillment", ServiceName: "Order Service", ColorClass: "node--order", Status: StagePending},
	}
}
