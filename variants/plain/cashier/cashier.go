package plain

import (
	"github.com/egonelbre/reactor-demo/variants/plain/item"
	"github.com/egonelbre/reactor-demo/variants/plain/msg"
	"github.com/egonelbre/reactor-demo/variants/plain/reactor"
)

type Reactor struct {
	ID        reactor.ID
	Database  reactor.ID
	GameWorld reactor.ID

	PendingPurchases map[reactor.RequestID]*PendingPurchaseRequest
}

type PendingPurchaseRequest struct {
	Status PendingPurchaseStatus

	UserRequestID reactor.RequestID

	UserID reactor.ID
	ItemID item.ID
}

type PendingPurchaseStatus byte

const (
	PendingPurchaseStatus_Unknown = iota
	PendingPurchaseStatus_WaitingDatabase
	PendingPurchaseStatus_WaitingWorld
)

func (cashier *Reactor) React(ev msg.Message, out *reactor.Outbox) {
	switch ev.Type() {
	case msg.CashierPurchaseItemRequestCode:
		var req msg.CashierPurchaseItemRequest
		userID := req.Parse(ev)

		message := msg.DatabasePurchaseItemRequest{
			RequestID: reactor.NewRequestID(),
			UserID:    userID,
			ItemID:    req.ItemID,
		}
		message.Send(out)

		cashier.PendingPurchases[message.RequestID] = &PendingPurchaseRequest{
			Status:        PendingPurchaseStatus_WaitingDatabase,
			UserRequestID: req.UserRequestID,
			UserID:        userID,
			ItemID:        req.ItemID,
		}

	case msg.DatabasePurchaseItemResponseCode:
		var resp msg.DatabasePurchaseItemResponse
		resp.Parse(ev)

		pending := cashier.PendingPurchases[resp.RequestID]
		assert(pending != nil)
		assert(pending.Status == PendingPurchaseStatus_WaitingDatabase)

		if !resp.Ok {
			delete(cashier.PendingPurchases, resp.RequestID)
			msg.CashierPurchaseItemResponse{
				UserRequestID: resp.RequestID,
				Ok:            false,
			}.Send(out)
			return
		}

		pending.Status = PendingPurchaseStatus_WaitingWorld
		msg.WorldAddItemRequest{
			RequestID: resp.RequestID,
			UserID:    pending.UserID,
			ItemID:    pending.ItemID,
		}.Send(out)

	case msg.WorldAddItemResponseCode:
		var resp msg.WorldAddItemResponse
		resp.Parse(ev)

		pending := cashier.PendingPurchases[resp.RequestID]
		assert(pending != nil)
		assert(pending.Status == PendingPurchaseStatus_WaitingWorld)

		msg.CashierPurchaseItemResponse{
			UserRequestID: pending.UserRequestID,
			Ok:            resp.Ok,
		}.Send(out)

		delete(cashier.PendingPurchases, resp.RequestID)
	}
}

func assert(v bool) {
	if !v {
		panic("assertion failed")
	}
}
