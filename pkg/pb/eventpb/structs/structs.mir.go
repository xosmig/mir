package eventpbstructs

import (
	model "github.com/filecoin-project/mir/codegen/proto-converter/model"
	structs "github.com/filecoin-project/mir/pkg/pb/availabilitypb/structs"
	bcbpb "github.com/filecoin-project/mir/pkg/pb/bcbpb"
	eventpb "github.com/filecoin-project/mir/pkg/pb/eventpb"
	isspb "github.com/filecoin-project/mir/pkg/pb/isspb"
	mempoolpb "github.com/filecoin-project/mir/pkg/pb/mempoolpb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

type Event struct {
	Type       Event_Type
	Next       []*Event
	DestModule string
}

type Event_Type interface {
	isEvent_Type()
	Pb() eventpb.Event_Type
}

func Event_TypeFromPb(pb eventpb.Event_Type) Event_Type {
	switch pb := pb.(type) {
	case *eventpb.Event_Init:
		return &Event_Init{Init: pb.Init}
	case *eventpb.Event_Tick:
		return &Event_Tick{Tick: pb.Tick}
	case *eventpb.Event_WalAppend:
		return &Event_WalAppend{WalAppend: pb.WalAppend}
	case *eventpb.Event_WalEntry:
		return &Event_WalEntry{WalEntry: pb.WalEntry}
	case *eventpb.Event_WalTruncate:
		return &Event_WalTruncate{WalTruncate: pb.WalTruncate}
	case *eventpb.Event_NewRequests:
		return &Event_NewRequests{NewRequests: pb.NewRequests}
	case *eventpb.Event_HashRequest:
		return &Event_HashRequest{HashRequest: pb.HashRequest}
	case *eventpb.Event_HashResult:
		return &Event_HashResult{HashResult: pb.HashResult}
	case *eventpb.Event_SignRequest:
		return &Event_SignRequest{SignRequest: pb.SignRequest}
	case *eventpb.Event_SignResult:
		return &Event_SignResult{SignResult: pb.SignResult}
	case *eventpb.Event_VerifyNodeSigs:
		return &Event_VerifyNodeSigs{VerifyNodeSigs: pb.VerifyNodeSigs}
	case *eventpb.Event_NodeSigsVerified:
		return &Event_NodeSigsVerified{NodeSigsVerified: pb.NodeSigsVerified}
	case *eventpb.Event_RequestReady:
		return &Event_RequestReady{RequestReady: pb.RequestReady}
	case *eventpb.Event_SendMessage:
		return &Event_SendMessage{SendMessage: pb.SendMessage}
	case *eventpb.Event_MessageReceived:
		return &Event_MessageReceived{MessageReceived: pb.MessageReceived}
	case *eventpb.Event_Deliver:
		return &Event_Deliver{Deliver: pb.Deliver}
	case *eventpb.Event_Iss:
		return &Event_Iss{Iss: pb.Iss}
	case *eventpb.Event_VerifyRequestSig:
		return &Event_VerifyRequestSig{VerifyRequestSig: pb.VerifyRequestSig}
	case *eventpb.Event_RequestSigVerified:
		return &Event_RequestSigVerified{RequestSigVerified: pb.RequestSigVerified}
	case *eventpb.Event_StoreVerifiedRequest:
		return &Event_StoreVerifiedRequest{StoreVerifiedRequest: pb.StoreVerifiedRequest}
	case *eventpb.Event_AppSnapshotRequest:
		return &Event_AppSnapshotRequest{AppSnapshotRequest: pb.AppSnapshotRequest}
	case *eventpb.Event_AppSnapshot:
		return &Event_AppSnapshot{AppSnapshot: pb.AppSnapshot}
	case *eventpb.Event_AppRestoreState:
		return &Event_AppRestoreState{AppRestoreState: pb.AppRestoreState}
	case *eventpb.Event_TimerDelay:
		return &Event_TimerDelay{TimerDelay: pb.TimerDelay}
	case *eventpb.Event_TimerRepeat:
		return &Event_TimerRepeat{TimerRepeat: pb.TimerRepeat}
	case *eventpb.Event_TimerGarbageCollect:
		return &Event_TimerGarbageCollect{TimerGarbageCollect: pb.TimerGarbageCollect}
	case *eventpb.Event_Bcb:
		return &Event_Bcb{Bcb: pb.Bcb}
	case *eventpb.Event_Mempool:
		return &Event_Mempool{Mempool: pb.Mempool}
	case *eventpb.Event_Availability:
		return &Event_Availability{Availability: structs.EventFromPb(pb.Availability)}
	case *eventpb.Event_NewEpoch:
		return &Event_NewEpoch{NewEpoch: pb.NewEpoch}
	case *eventpb.Event_NewConfig:
		return &Event_NewConfig{NewConfig: pb.NewConfig}
	case *eventpb.Event_TestingString:
		return &Event_TestingString{TestingString: pb.TestingString}
	case *eventpb.Event_TestingUint:
		return &Event_TestingUint{TestingUint: pb.TestingUint}
	}
	return nil
}

type Event_Init struct {
	Init *eventpb.Init
}

func (*Event_Init) isEvent_Type() {}

func (w *Event_Init) Pb() eventpb.Event_Type {
	return &eventpb.Event_Init{Init: w.Init}
}

type Event_Tick struct {
	Tick *eventpb.Tick
}

func (*Event_Tick) isEvent_Type() {}

func (w *Event_Tick) Pb() eventpb.Event_Type {
	return &eventpb.Event_Tick{Tick: w.Tick}
}

type Event_WalAppend struct {
	WalAppend *eventpb.WALAppend
}

func (*Event_WalAppend) isEvent_Type() {}

func (w *Event_WalAppend) Pb() eventpb.Event_Type {
	return &eventpb.Event_WalAppend{WalAppend: w.WalAppend}
}

type Event_WalEntry struct {
	WalEntry *eventpb.WALEntry
}

func (*Event_WalEntry) isEvent_Type() {}

func (w *Event_WalEntry) Pb() eventpb.Event_Type {
	return &eventpb.Event_WalEntry{WalEntry: w.WalEntry}
}

type Event_WalTruncate struct {
	WalTruncate *eventpb.WALTruncate
}

func (*Event_WalTruncate) isEvent_Type() {}

func (w *Event_WalTruncate) Pb() eventpb.Event_Type {
	return &eventpb.Event_WalTruncate{WalTruncate: w.WalTruncate}
}

type Event_NewRequests struct {
	NewRequests *eventpb.NewRequests
}

func (*Event_NewRequests) isEvent_Type() {}

func (w *Event_NewRequests) Pb() eventpb.Event_Type {
	return &eventpb.Event_NewRequests{NewRequests: w.NewRequests}
}

type Event_HashRequest struct {
	HashRequest *eventpb.HashRequest
}

func (*Event_HashRequest) isEvent_Type() {}

func (w *Event_HashRequest) Pb() eventpb.Event_Type {
	return &eventpb.Event_HashRequest{HashRequest: w.HashRequest}
}

type Event_HashResult struct {
	HashResult *eventpb.HashResult
}

func (*Event_HashResult) isEvent_Type() {}

func (w *Event_HashResult) Pb() eventpb.Event_Type {
	return &eventpb.Event_HashResult{HashResult: w.HashResult}
}

type Event_SignRequest struct {
	SignRequest *eventpb.SignRequest
}

func (*Event_SignRequest) isEvent_Type() {}

func (w *Event_SignRequest) Pb() eventpb.Event_Type {
	return &eventpb.Event_SignRequest{SignRequest: w.SignRequest}
}

type Event_SignResult struct {
	SignResult *eventpb.SignResult
}

func (*Event_SignResult) isEvent_Type() {}

func (w *Event_SignResult) Pb() eventpb.Event_Type {
	return &eventpb.Event_SignResult{SignResult: w.SignResult}
}

type Event_VerifyNodeSigs struct {
	VerifyNodeSigs *eventpb.VerifyNodeSigs
}

func (*Event_VerifyNodeSigs) isEvent_Type() {}

func (w *Event_VerifyNodeSigs) Pb() eventpb.Event_Type {
	return &eventpb.Event_VerifyNodeSigs{VerifyNodeSigs: w.VerifyNodeSigs}
}

type Event_NodeSigsVerified struct {
	NodeSigsVerified *eventpb.NodeSigsVerified
}

func (*Event_NodeSigsVerified) isEvent_Type() {}

func (w *Event_NodeSigsVerified) Pb() eventpb.Event_Type {
	return &eventpb.Event_NodeSigsVerified{NodeSigsVerified: w.NodeSigsVerified}
}

type Event_RequestReady struct {
	RequestReady *eventpb.RequestReady
}

func (*Event_RequestReady) isEvent_Type() {}

func (w *Event_RequestReady) Pb() eventpb.Event_Type {
	return &eventpb.Event_RequestReady{RequestReady: w.RequestReady}
}

type Event_SendMessage struct {
	SendMessage *eventpb.SendMessage
}

func (*Event_SendMessage) isEvent_Type() {}

func (w *Event_SendMessage) Pb() eventpb.Event_Type {
	return &eventpb.Event_SendMessage{SendMessage: w.SendMessage}
}

type Event_MessageReceived struct {
	MessageReceived *eventpb.MessageReceived
}

func (*Event_MessageReceived) isEvent_Type() {}

func (w *Event_MessageReceived) Pb() eventpb.Event_Type {
	return &eventpb.Event_MessageReceived{MessageReceived: w.MessageReceived}
}

type Event_Deliver struct {
	Deliver *eventpb.Deliver
}

func (*Event_Deliver) isEvent_Type() {}

func (w *Event_Deliver) Pb() eventpb.Event_Type {
	return &eventpb.Event_Deliver{Deliver: w.Deliver}
}

type Event_Iss struct {
	Iss *isspb.ISSEvent
}

func (*Event_Iss) isEvent_Type() {}

func (w *Event_Iss) Pb() eventpb.Event_Type {
	return &eventpb.Event_Iss{Iss: w.Iss}
}

type Event_VerifyRequestSig struct {
	VerifyRequestSig *eventpb.VerifyRequestSig
}

func (*Event_VerifyRequestSig) isEvent_Type() {}

func (w *Event_VerifyRequestSig) Pb() eventpb.Event_Type {
	return &eventpb.Event_VerifyRequestSig{VerifyRequestSig: w.VerifyRequestSig}
}

type Event_RequestSigVerified struct {
	RequestSigVerified *eventpb.RequestSigVerified
}

func (*Event_RequestSigVerified) isEvent_Type() {}

func (w *Event_RequestSigVerified) Pb() eventpb.Event_Type {
	return &eventpb.Event_RequestSigVerified{RequestSigVerified: w.RequestSigVerified}
}

type Event_StoreVerifiedRequest struct {
	StoreVerifiedRequest *eventpb.StoreVerifiedRequest
}

func (*Event_StoreVerifiedRequest) isEvent_Type() {}

func (w *Event_StoreVerifiedRequest) Pb() eventpb.Event_Type {
	return &eventpb.Event_StoreVerifiedRequest{StoreVerifiedRequest: w.StoreVerifiedRequest}
}

type Event_AppSnapshotRequest struct {
	AppSnapshotRequest *eventpb.AppSnapshotRequest
}

func (*Event_AppSnapshotRequest) isEvent_Type() {}

func (w *Event_AppSnapshotRequest) Pb() eventpb.Event_Type {
	return &eventpb.Event_AppSnapshotRequest{AppSnapshotRequest: w.AppSnapshotRequest}
}

type Event_AppSnapshot struct {
	AppSnapshot *eventpb.AppSnapshot
}

func (*Event_AppSnapshot) isEvent_Type() {}

func (w *Event_AppSnapshot) Pb() eventpb.Event_Type {
	return &eventpb.Event_AppSnapshot{AppSnapshot: w.AppSnapshot}
}

type Event_AppRestoreState struct {
	AppRestoreState *eventpb.AppRestoreState
}

func (*Event_AppRestoreState) isEvent_Type() {}

func (w *Event_AppRestoreState) Pb() eventpb.Event_Type {
	return &eventpb.Event_AppRestoreState{AppRestoreState: w.AppRestoreState}
}

type Event_TimerDelay struct {
	TimerDelay *eventpb.TimerDelay
}

func (*Event_TimerDelay) isEvent_Type() {}

func (w *Event_TimerDelay) Pb() eventpb.Event_Type {
	return &eventpb.Event_TimerDelay{TimerDelay: w.TimerDelay}
}

type Event_TimerRepeat struct {
	TimerRepeat *eventpb.TimerRepeat
}

func (*Event_TimerRepeat) isEvent_Type() {}

func (w *Event_TimerRepeat) Pb() eventpb.Event_Type {
	return &eventpb.Event_TimerRepeat{TimerRepeat: w.TimerRepeat}
}

type Event_TimerGarbageCollect struct {
	TimerGarbageCollect *eventpb.TimerGarbageCollect
}

func (*Event_TimerGarbageCollect) isEvent_Type() {}

func (w *Event_TimerGarbageCollect) Pb() eventpb.Event_Type {
	return &eventpb.Event_TimerGarbageCollect{TimerGarbageCollect: w.TimerGarbageCollect}
}

type Event_Bcb struct {
	Bcb *bcbpb.Event
}

func (*Event_Bcb) isEvent_Type() {}

func (w *Event_Bcb) Pb() eventpb.Event_Type {
	return &eventpb.Event_Bcb{Bcb: w.Bcb}
}

type Event_Mempool struct {
	Mempool *mempoolpb.Event
}

func (*Event_Mempool) isEvent_Type() {}

func (w *Event_Mempool) Pb() eventpb.Event_Type {
	return &eventpb.Event_Mempool{Mempool: w.Mempool}
}

type Event_Availability struct {
	Availability *structs.Event
}

func (*Event_Availability) isEvent_Type() {}

func (w *Event_Availability) Pb() eventpb.Event_Type {
	return &eventpb.Event_Availability{Availability: (w.Availability).Pb()}
}

type Event_NewEpoch struct {
	NewEpoch *eventpb.NewEpoch
}

func (*Event_NewEpoch) isEvent_Type() {}

func (w *Event_NewEpoch) Pb() eventpb.Event_Type {
	return &eventpb.Event_NewEpoch{NewEpoch: w.NewEpoch}
}

type Event_NewConfig struct {
	NewConfig *eventpb.NewConfig
}

func (*Event_NewConfig) isEvent_Type() {}

func (w *Event_NewConfig) Pb() eventpb.Event_Type {
	return &eventpb.Event_NewConfig{NewConfig: w.NewConfig}
}

type Event_TestingString struct {
	TestingString *wrapperspb.StringValue
}

func (*Event_TestingString) isEvent_Type() {}

func (w *Event_TestingString) Pb() eventpb.Event_Type {
	return &eventpb.Event_TestingString{TestingString: w.TestingString}
}

type Event_TestingUint struct {
	TestingUint *wrapperspb.UInt64Value
}

func (*Event_TestingUint) isEvent_Type() {}

func (w *Event_TestingUint) Pb() eventpb.Event_Type {
	return &eventpb.Event_TestingUint{TestingUint: w.TestingUint}
}

func (m *Event) Pb() *eventpb.Event {
	return &eventpb.Event{
		Type: (m.Type).Pb(),
		Next: model.ConvertSlice(m.Next, func(t *Event) *eventpb.Event {
			return (t).Pb()
		}),
		DestModule: m.DestModule,
	}
}

func EventFromPb(pb *eventpb.Event) *Event {
	return &Event{
		Type: Event_TypeFromPb(pb.Type),
		Next: model.ConvertSlice(pb.Next, func(t *eventpb.Event) *Event {
			return EventFromPb(t)
		}),
		DestModule: pb.DestModule,
	}
}
