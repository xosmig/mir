package messagepb

import (
	mscpb "github.com/filecoin-project/mir/pkg/pb/availabilitypb/mscpb"
	bcbpb "github.com/filecoin-project/mir/pkg/pb/bcbpb"
	isspb "github.com/filecoin-project/mir/pkg/pb/isspb"
)

type Message_Type = isMessage_Type

type Message_TypeWrapper[T any] interface {
	Message_Type
	Unwrap() *T
}

func (w *Message_Iss) Unwrap() *isspb.ISSMessage {
	return w.Iss
}

func (w *Message_Bcb) Unwrap() *bcbpb.Message {
	return w.Bcb
}

func (w *Message_MultisigCollector) Unwrap() *mscpb.Message {
	return w.MultisigCollector
}
