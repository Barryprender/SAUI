package distributed

import (
	"context"
	"time"

	"saui/statestore"
)

const PipelineProjectionName = "distributed.pipeline"

type StageStatus string

const (
	StagePending StageStatus = "pending"
	StageActive  StageStatus = "active"
	StageDone    StageStatus = "done"
	StageFailed  StageStatus = "failed"
)

type OrderStatus string

const (
	StatusNew                 OrderStatus = "new"
	StatusAwaitingInventory   OrderStatus = "awaiting_inventory"
	StatusAwaitingPayment     OrderStatus = "awaiting_payment"
	StatusPaymentDeclined     OrderStatus = "payment_declined"
	StatusAwaitingFulfillment OrderStatus = "awaiting_fulfillment"
	StatusFulfilled           OrderStatus = "fulfilled"
	StatusCancelled           OrderStatus = "cancelled"
)

type PipelineNode struct {
	Index         int
	Label         string
	ServiceName   string
	ColorClass    string
	Status        StageStatus
	EventID       int64
	JustCompleted bool
}

type EventLogEntry struct {
	EventID     int64
	Type        string
	ServiceName string
	ColorClass  string
	Label       string
	OccurredAt  time.Time
}

type PipelineProjection struct {
	Nodes           []PipelineNode
	EventLog        []EventLogEntry
	Status          OrderStatus
	OrderNumber     int
	LatestEventType string
}

func (p PipelineProjection) Name() string { return PipelineProjectionName }

func init() {
	statestore.RegisterProjection(PipelineProjectionName, func(ctx context.Context, s *statestore.Store, sessionID string) (statestore.Projection, error) {
		events, err := s.EventsBySession(ctx, sessionID)
		if err != nil {
			return nil, err
		}
		return buildPipelineFromEvents(events), nil
	})
}

func buildPipelineFromEvents(allEvents []statestore.Event) PipelineProjection {
	// Find the start of the current cycle (after last terminal event).
	cycleStart := 0
	orderCount := 0
	for i, ev := range allEvents {
		if ev.Type == OrderCreated {
			orderCount++
		}
		if ev.Type == OrderFulfilled || ev.Type == OrderCancelled {
			cycleStart = i + 1
		}
	}

	events := allEvents[cycleStart:]
	nodes := freshNodes()
	status := StatusNew
	var eventLog []EventLogEntry
	var latestType string

	for _, ev := range events {
		meta, ok := eventMetaFor(ev.Type)
		if !ok {
			continue
		}
		eventLog = append(eventLog, EventLogEntry{
			EventID:     ev.ID,
			Type:        ev.Type,
			ServiceName: meta.ServiceName,
			ColorClass:  meta.ColorClass,
			Label:       meta.Label,
			OccurredAt:  ev.OccurredAt,
		})
		latestType = ev.Type

		switch ev.Type {
		case OrderCreated:
			status = StatusAwaitingInventory
			nodes[0].Status = StageDone
			nodes[0].EventID = ev.ID
		case InventoryReserved:
			status = StatusAwaitingPayment
			nodes[1].Status = StageDone
			nodes[1].EventID = ev.ID
		case PaymentCaptured:
			status = StatusAwaitingFulfillment
			nodes[2].Status = StageDone
			nodes[2].EventID = ev.ID
		case PaymentDeclined:
			status = StatusPaymentDeclined
			nodes[2].Status = StageFailed
			nodes[2].EventID = ev.ID
		case OrderFulfilled:
			status = StatusFulfilled
			nodes[3].Status = StageDone
			nodes[3].EventID = ev.ID
		case OrderCancelled:
			status = StatusCancelled
		}
	}

	// Mark the active node for the current waiting state.
	switch status {
	case StatusAwaitingInventory:
		nodes[1].Status = StageActive
	case StatusAwaitingPayment:
		nodes[2].Status = StageActive
	case StatusAwaitingFulfillment:
		nodes[3].Status = StageActive
	}

	// Reverse event log: newest first.
	for i, j := 0, len(eventLog)-1; i < j; i, j = i+1, j-1 {
		eventLog[i], eventLog[j] = eventLog[j], eventLog[i]
	}

	// Flag the node that just transitioned for the pulse animation.
	if idx := nodeIndexForEvent(latestType); idx >= 0 {
		nodes[idx].JustCompleted = true
	}

	return PipelineProjection{
		Nodes:           nodes,
		EventLog:        eventLog,
		Status:          status,
		OrderNumber:     orderCount,
		LatestEventType: latestType,
	}
}

// currentStatusFromStore is used by actions to validate pipeline state.
func currentStatusFromStore(ctx context.Context, s *statestore.Store, sessionID string) (OrderStatus, error) {
	events, err := s.EventsBySession(ctx, sessionID)
	if err != nil {
		return StatusNew, err
	}
	return buildPipelineFromEvents(events).Status, nil
}
