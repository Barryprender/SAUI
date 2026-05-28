package banking

import (
	"context"
	"time"

	"saui/statestore"
)

const AccountProjectionName = "banking.account"

// TxRecord is a single transfer in the account history.
type TxRecord struct {
	EventID             int64
	RecipientID         string
	RecipientName       string
	RecipientInitials   string
	RecipientColorClass string
	AmountCents         int64
	OccurredAt          time.Time
}

// AccountProjection is the read model for the banking demo session.
type AccountProjection struct {
	IsFunded     bool
	BalanceCents int64
	Transfers    []TxRecord // newest first
}

func (p AccountProjection) Name() string { return AccountProjectionName }

func init() {
	statestore.RegisterProjection(AccountProjectionName, buildAccount)
}

func buildAccount(ctx context.Context, s *statestore.Store, sessionID string) (statestore.Projection, error) {
	events, err := s.EventsBySession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	balance, funded := balanceFromEvents(events)
	entries := txFromEvents(events)

	proj := AccountProjection{
		IsFunded:     funded,
		BalanceCents: balance,
	}
	// Reverse to newest-first
	for i := len(entries) - 1; i >= 0; i-- {
		e := entries[i]
		proj.Transfers = append(proj.Transfers, TxRecord{
			EventID:             e.EventID,
			RecipientID:         e.RecipientID,
			RecipientName:       e.RecipientName,
			RecipientInitials:   e.RecipientInitials,
			RecipientColorClass: e.RecipientColorClass,
			AmountCents:         e.AmountCents,
			OccurredAt:          e.OccurredAt,
		})
	}
	return proj, nil
}
