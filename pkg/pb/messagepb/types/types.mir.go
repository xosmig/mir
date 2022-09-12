package messagepbtypes

import (
	mirreflect "github.com/filecoin-project/mir/codegen/mirreflect"
	mscpb "github.com/filecoin-project/mir/pkg/pb/availabilitypb/mscpb"
	types "github.com/filecoin-project/mir/pkg/pb/bcbpb/types"
	isspb "github.com/filecoin-project/mir/pkg/pb/isspb"
	messagepb "github.com/filecoin-project/mir/pkg/pb/messagepb"
	reflectutil "github.com/filecoin-project/mir/pkg/util/reflectutil"
)

type Message struct {
	DestModule string
	Type       Message_Type
}

type Message_Type interface {
	mirreflect.GeneratedType
	isMessage_Type()
	Pb() messagepb.Message_Type
}

type Message_TypeWrapper[T any] interface {
	Message_Type
	Unwrap() *T
}

func Message_TypeFromPb(pb messagepb.Message_Type) Message_Type {
	switch pb := pb.(type) {
	case *messagepb.Message_Iss:
		return &Message_Iss{Iss: pb.Iss}
	case *messagepb.Message_Bcb:
		return &Message_Bcb{Bcb: types.MessageFromPb(pb.Bcb)}
	case *messagepb.Message_MultisigCollector:
		return &Message_MultisigCollector{MultisigCollector: pb.MultisigCollector}
	}
	return nil
}

type Message_Iss struct {
	Iss *isspb.ISSMessage
}

func (*Message_Iss) isMessage_Type() {}

func (w *Message_Iss) Unwrap() *isspb.ISSMessage {
	return w.Iss
}

func (w *Message_Iss) Pb() messagepb.Message_Type {
	return &messagepb.Message_Iss{Iss: w.Iss}
}

func (*Message_Iss) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*messagepb.Message_Iss]()}
}

type Message_Bcb struct {
	Bcb *types.Message
}

func (*Message_Bcb) isMessage_Type() {}

func (w *Message_Bcb) Unwrap() *types.Message {
	return w.Bcb
}

func (w *Message_Bcb) Pb() messagepb.Message_Type {
	return &messagepb.Message_Bcb{Bcb: (w.Bcb).Pb()}
}

func (*Message_Bcb) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*messagepb.Message_Bcb]()}
}

type Message_MultisigCollector struct {
	MultisigCollector *mscpb.Message
}

func (*Message_MultisigCollector) isMessage_Type() {}

func (w *Message_MultisigCollector) Unwrap() *mscpb.Message {
	return w.MultisigCollector
}

func (w *Message_MultisigCollector) Pb() messagepb.Message_Type {
	return &messagepb.Message_MultisigCollector{MultisigCollector: w.MultisigCollector}
}

func (*Message_MultisigCollector) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*messagepb.Message_MultisigCollector]()}
}

func MessageFromPb(pb *messagepb.Message) *Message {
	return &Message{
		DestModule: pb.DestModule,
		Type:       Message_TypeFromPb(pb.Type),
	}
}

func (m *Message) Pb() *messagepb.Message {
	return &messagepb.Message{
		DestModule: m.DestModule,
		Type:       (m.Type).Pb(),
	}
}

func (*Message) MirReflect() mirreflect.Type {
	return mirreflect.TypeImpl{PbType_: reflectutil.TypeOf[*messagepb.Message]()}
}
