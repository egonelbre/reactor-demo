package reactor

type ID uint32

type RequestID uint32

type Outbox struct {
	Owner  ID
	Buffer []uint32
}

func NewRequestID() RequestID { return 0 }
