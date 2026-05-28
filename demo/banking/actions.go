package banking

import (
	"context"
	"errors"

	"saui/statestore"
)

// FundAccount seeds the initial €2,500 balance once per session.
type FundAccount struct{}

func (a FundAccount) Type() string { return AccountFunded }

func (a FundAccount) Validate(ctx context.Context, s *statestore.Store, sessionID string) error {
	n, err := s.CountEventsBySessionAndType(ctx, sessionID, AccountFunded)
	if err != nil {
		return err
	}
	if n > 0 {
		return errors.New("already funded")
	}
	return nil
}

func (a FundAccount) Apply(ctx context.Context, s *statestore.Store, sessionID string) (statestore.Event, error) {
	return s.AppendEvent(ctx, sessionID, AccountFunded, map[string]any{
		"amount_cents": int64(250000),
	})
}

// SendTransfer debits the account and records a transfer to a preset recipient.
type SendTransfer struct {
	RecipientID string
	AmountCents int64
}

func (a SendTransfer) Type() string { return TransferSent }

func (a SendTransfer) Validate(ctx context.Context, s *statestore.Store, sessionID string) error {
	if a.RecipientID == "" {
		return errors.New("select a recipient")
	}
	if a.AmountCents <= 0 {
		return errors.New("amount must be positive")
	}
	if _, ok := RecipientByID(a.RecipientID); !ok {
		return errors.New("unknown recipient")
	}
	events, err := s.EventsBySession(ctx, sessionID)
	if err != nil {
		return err
	}
	balance, _ := balanceFromEvents(events)
	if a.AmountCents > balance {
		return errors.New("insufficient funds")
	}
	return nil
}

func (a SendTransfer) Apply(ctx context.Context, s *statestore.Store, sessionID string) (statestore.Event, error) {
	r, _ := RecipientByID(a.RecipientID)
	return s.AppendEvent(ctx, sessionID, TransferSent, map[string]any{
		"recipient_id":   a.RecipientID,
		"recipient_name": r.Name,
		"amount_cents":   a.AmountCents,
	})
}
