package msg

import (
	"github.com/egonelbre/reactor-demo/variants/plain/item"
	"github.com/egonelbre/reactor-demo/variants/plain/reactor"
)

type Message []uint32

func (m Message) Type() MessageType {
	return MessageType(m[0])
}

func (m Message) Assert(mt MessageType, size int) {
	if m[0] != uint32(mt) {
		panic("expected something else")
	}
	if len(m) < size+2 {
		panic("invalid size")
	}
}

type MessageType uint32

const (
	DatabasePurchaseItemRequestCode MessageType = iota + 1
	DatabasePurchaseItemResponseCode
	CashierPurchaseItemRequestCode
	CashierPurchaseItemResponseCode
	WorldAddItemRequestCode
	WorldAddItemResponseCode
)

func (m MessageType) Len() int {
	switch m {
	case CashierPurchaseItemRequestCode:
		return 2
	case CashierPurchaseItemResponseCode:
		return 2
	case DatabasePurchaseItemRequestCode:
		return 3
	case DatabasePurchaseItemResponseCode:
		return 2
	case WorldAddItemRequestCode:
		return 3
	case WorldAddItemResponseCode:
		return 2
	default:
		return -1
	}
}

type CashierPurchaseItemRequest struct {
	UserRequestID reactor.RequestID
	ItemID        item.ID
}

func (m CashierPurchaseItemRequest) Send(outbox *reactor.Outbox) {
	outbox.Buffer = append(outbox.Buffer,
		uint32(CashierPurchaseItemRequestCode),
		uint32(outbox.Owner),
		uint32(m.UserRequestID),
		uint32(m.ItemID),
	)
}

func (m *CashierPurchaseItemRequest) Parse(message Message) (sender reactor.ID) {
	message.Assert(CashierPurchaseItemRequestCode, CashierPurchaseItemRequestCode.Len())
	m.UserRequestID = reactor.RequestID(message[2])
	m.ItemID = item.ID(message[3])
	return reactor.ID(message[1])
}

type CashierPurchaseItemResponse struct {
	UserRequestID reactor.RequestID
	Ok            bool
}

func (m CashierPurchaseItemResponse) Send(outbox *reactor.Outbox) {
	outbox.Buffer = append(outbox.Buffer,
		uint32(CashierPurchaseItemRequestCode),
		uint32(outbox.Owner),
		uint32(m.UserRequestID),
		boolToUint32(m.Ok),
	)
}

func (m *CashierPurchaseItemResponse) Parse(message Message) (sender reactor.ID) {
	message.Assert(CashierPurchaseItemRequestCode, CashierPurchaseItemRequestCode.Len())
	m.UserRequestID = reactor.RequestID(message[2])
	m.Ok = uint32ToBool(message[3])
	return reactor.ID(message[1])
}

type DatabasePurchaseItemRequest struct {
	RequestID reactor.RequestID
	UserID    reactor.ID
	ItemID    item.ID
}

func (m DatabasePurchaseItemRequest) Send(outbox *reactor.Outbox) {
	outbox.Buffer = append(outbox.Buffer,
		uint32(DatabasePurchaseItemRequestCode),
		uint32(outbox.Owner),
		uint32(m.RequestID),
		uint32(m.UserID),
		uint32(m.ItemID),
	)
}

func (m *DatabasePurchaseItemRequest) Parse(message Message) (sender reactor.ID) {
	message.Assert(DatabasePurchaseItemRequestCode, DatabasePurchaseItemRequestCode.Len())
	m.RequestID = reactor.RequestID(message[2])
	m.UserID = reactor.ID(message[3])
	m.ItemID = item.ID(message[4])
	return reactor.ID(message[1])
}

type DatabasePurchaseItemResponse struct {
	RequestID reactor.RequestID
	Ok        bool
}

func (m DatabasePurchaseItemResponse) Send(outbox *reactor.Outbox) {
	outbox.Buffer = append(outbox.Buffer,
		uint32(DatabasePurchaseItemResponseCode),
		uint32(outbox.Owner),
		uint32(m.RequestID),
		boolToUint32(m.Ok),
	)
}

func (m *DatabasePurchaseItemResponse) Parse(message Message) (sender reactor.ID) {
	message.Assert(DatabasePurchaseItemResponseCode, DatabasePurchaseItemResponseCode.Len())
	m.RequestID = reactor.RequestID(message[2])
	m.Ok = uint32ToBool(message[2])
	return reactor.ID(message[1])
}

type WorldAddItemRequest struct {
	RequestID reactor.RequestID
	UserID    reactor.ID
	ItemID    item.ID
}

func (m WorldAddItemRequest) Send(outbox *reactor.Outbox) {
	outbox.Buffer = append(outbox.Buffer,
		uint32(WorldAddItemRequestCode),
		uint32(outbox.Owner),
		uint32(m.RequestID),
		uint32(m.UserID),
		uint32(m.ItemID),
	)
}

func (m *WorldAddItemRequest) Parse(message Message) (sender reactor.ID) {
	message.Assert(WorldAddItemRequestCode, WorldAddItemRequestCode.Len())
	m.RequestID = reactor.RequestID(message[2])
	m.UserID = reactor.ID(message[3])
	m.ItemID = item.ID(message[4])
	return reactor.ID(message[1])
}

type WorldAddItemResponse struct {
	RequestID reactor.RequestID
	Ok        bool
}

func (m WorldAddItemResponse) Send(outbox *reactor.Outbox) {
	outbox.Buffer = append(outbox.Buffer,
		uint32(WorldAddItemResponseCode),
		uint32(outbox.Owner),
		uint32(m.RequestID),
		boolToUint32(m.Ok),
	)
}

func (m *WorldAddItemResponse) Parse(message Message) (sender reactor.ID) {
	message.Assert(WorldAddItemResponseCode, WorldAddItemResponseCode.Len())
	m.RequestID = reactor.RequestID(message[2])
	m.Ok = uint32ToBool(message[2])
	return reactor.ID(message[1])
}

func uint32ToBool(v uint32) bool { return v == 1 }
func boolToUint32(v bool) uint32 {
	if v {
		return 1
	}
	return 0
}
