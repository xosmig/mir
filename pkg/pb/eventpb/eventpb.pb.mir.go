package eventpb

import (
	reflect "reflect"
)

func (*Event) ReflectTypeOptions() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*Event_Init)(nil)),
		reflect.TypeOf((*Event_Tick)(nil)),
		reflect.TypeOf((*Event_WalAppend)(nil)),
		reflect.TypeOf((*Event_WalEntry)(nil)),
		reflect.TypeOf((*Event_WalTruncate)(nil)),
		reflect.TypeOf((*Event_NewRequests)(nil)),
		reflect.TypeOf((*Event_HashRequest)(nil)),
		reflect.TypeOf((*Event_HashResult)(nil)),
		reflect.TypeOf((*Event_SignRequest)(nil)),
		reflect.TypeOf((*Event_SignResult)(nil)),
		reflect.TypeOf((*Event_VerifyNodeSigs)(nil)),
		reflect.TypeOf((*Event_NodeSigsVerified)(nil)),
		reflect.TypeOf((*Event_RequestReady)(nil)),
		reflect.TypeOf((*Event_SendMessage)(nil)),
		reflect.TypeOf((*Event_MessageReceived)(nil)),
		reflect.TypeOf((*Event_Deliver)(nil)),
		reflect.TypeOf((*Event_Iss)(nil)),
		reflect.TypeOf((*Event_VerifyRequestSig)(nil)),
		reflect.TypeOf((*Event_RequestSigVerified)(nil)),
		reflect.TypeOf((*Event_StoreVerifiedRequest)(nil)),
		reflect.TypeOf((*Event_AppSnapshotRequest)(nil)),
		reflect.TypeOf((*Event_AppSnapshot)(nil)),
		reflect.TypeOf((*Event_AppRestoreState)(nil)),
		reflect.TypeOf((*Event_TimerDelay)(nil)),
		reflect.TypeOf((*Event_TimerRepeat)(nil)),
		reflect.TypeOf((*Event_TimerGarbageCollect)(nil)),
		reflect.TypeOf((*Event_Bcb)(nil)),
		reflect.TypeOf((*Event_Mempool)(nil)),
		reflect.TypeOf((*Event_Availability)(nil)),
		reflect.TypeOf((*Event_NewEpoch)(nil)),
		reflect.TypeOf((*Event_NewConfig)(nil)),
		reflect.TypeOf((*Event_TestingString)(nil)),
		reflect.TypeOf((*Event_TestingUint)(nil)),
	}
}
