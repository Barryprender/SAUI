package distributed

import (
	"context"
	"errors"

	"saui/statestore"
)

type PlaceOrder struct{}

func (a PlaceOrder) Type() string { return OrderCreated }

func (a PlaceOrder) Validate(ctx context.Context, s *statestore.Store, sessionID string) error {
	status, err := currentStatusFromStore(ctx, s, sessionID)
	if err != nil {
		return err
	}
	if status != StatusNew && status != StatusFulfilled && status != StatusCancelled {
		return errors.New("order already in progress")
	}
	return nil
}

func (a PlaceOrder) Apply(ctx context.Context, s *statestore.Store, sessionID string) (statestore.Event, error) {
	return s.AppendEvent(ctx, sessionID, OrderCreated, map[string]any{
		"service": "order",
	})
}

type ReserveInventory struct{}

func (a ReserveInventory) Type() string { return InventoryReserved }

func (a ReserveInventory) Validate(ctx context.Context, s *statestore.Store, sessionID string) error {
	status, err := currentStatusFromStore(ctx, s, sessionID)
	if err != nil {
		return err
	}
	if status != StatusAwaitingInventory {
		return errors.New("not awaiting inventory")
	}
	return nil
}

func (a ReserveInventory) Apply(ctx context.Context, s *statestore.Store, sessionID string) (statestore.Event, error) {
	return s.AppendEvent(ctx, sessionID, InventoryReserved, map[string]any{
		"service": "inventory",
	})
}

type CapturePayment struct{}

func (a CapturePayment) Type() string { return PaymentCaptured }

func (a CapturePayment) Validate(ctx context.Context, s *statestore.Store, sessionID string) error {
	status, err := currentStatusFromStore(ctx, s, sessionID)
	if err != nil {
		return err
	}
	if status != StatusAwaitingPayment && status != StatusPaymentDeclined {
		return errors.New("not awaiting payment")
	}
	return nil
}

func (a CapturePayment) Apply(ctx context.Context, s *statestore.Store, sessionID string) (statestore.Event, error) {
	return s.AppendEvent(ctx, sessionID, PaymentCaptured, map[string]any{
		"service": "payment",
	})
}

type DeclinePayment struct{}

func (a DeclinePayment) Type() string { return PaymentDeclined }

func (a DeclinePayment) Validate(ctx context.Context, s *statestore.Store, sessionID string) error {
	status, err := currentStatusFromStore(ctx, s, sessionID)
	if err != nil {
		return err
	}
	if status != StatusAwaitingPayment {
		return errors.New("not awaiting payment")
	}
	return nil
}

func (a DeclinePayment) Apply(ctx context.Context, s *statestore.Store, sessionID string) (statestore.Event, error) {
	return s.AppendEvent(ctx, sessionID, PaymentDeclined, map[string]any{
		"service": "payment",
	})
}

type FulfillOrder struct{}

func (a FulfillOrder) Type() string { return OrderFulfilled }

func (a FulfillOrder) Validate(ctx context.Context, s *statestore.Store, sessionID string) error {
	status, err := currentStatusFromStore(ctx, s, sessionID)
	if err != nil {
		return err
	}
	if status != StatusAwaitingFulfillment {
		return errors.New("not awaiting fulfillment")
	}
	return nil
}

func (a FulfillOrder) Apply(ctx context.Context, s *statestore.Store, sessionID string) (statestore.Event, error) {
	return s.AppendEvent(ctx, sessionID, OrderFulfilled, map[string]any{
		"service": "order",
	})
}

type CancelOrder struct{}

func (a CancelOrder) Type() string { return OrderCancelled }

func (a CancelOrder) Validate(ctx context.Context, s *statestore.Store, sessionID string) error {
	status, err := currentStatusFromStore(ctx, s, sessionID)
	if err != nil {
		return err
	}
	if status == StatusNew || status == StatusFulfilled || status == StatusCancelled {
		return errors.New("nothing to cancel")
	}
	return nil
}

func (a CancelOrder) Apply(ctx context.Context, s *statestore.Store, sessionID string) (statestore.Event, error) {
	return s.AppendEvent(ctx, sessionID, OrderCancelled, map[string]any{
		"service": "order",
	})
}
