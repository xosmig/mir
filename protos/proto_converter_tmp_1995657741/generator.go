
package main

import (
	"log"
	"os"
	"reflect"

	"github.com/filecoin-project/mir/codegen/proto-converter/codegen"
	pkg_ "github.com/filecoin-project/mir/pkg/pb/eventpb"
)

func main() {
	err := codegen.GenerateAll(
		[]reflect.Type{
			
			reflect.TypeOf((*pkg_.Event)(nil)),
			
			reflect.TypeOf((*pkg_.Event_Init)(nil)),
			
			reflect.TypeOf((*pkg_.Event_Tick)(nil)),
			
			reflect.TypeOf((*pkg_.Event_WalAppend)(nil)),
			
			reflect.TypeOf((*pkg_.Event_WalEntry)(nil)),
			
			reflect.TypeOf((*pkg_.Event_WalTruncate)(nil)),
			
			reflect.TypeOf((*pkg_.Event_NewRequests)(nil)),
			
			reflect.TypeOf((*pkg_.Event_HashRequest)(nil)),
			
			reflect.TypeOf((*pkg_.Event_HashResult)(nil)),
			
			reflect.TypeOf((*pkg_.Event_SignRequest)(nil)),
			
			reflect.TypeOf((*pkg_.Event_SignResult)(nil)),
			
			reflect.TypeOf((*pkg_.Event_VerifyNodeSigs)(nil)),
			
			reflect.TypeOf((*pkg_.Event_NodeSigsVerified)(nil)),
			
			reflect.TypeOf((*pkg_.Event_RequestReady)(nil)),
			
			reflect.TypeOf((*pkg_.Event_SendMessage)(nil)),
			
			reflect.TypeOf((*pkg_.Event_MessageReceived)(nil)),
			
			reflect.TypeOf((*pkg_.Event_Deliver)(nil)),
			
			reflect.TypeOf((*pkg_.Event_Iss)(nil)),
			
			reflect.TypeOf((*pkg_.Event_VerifyRequestSig)(nil)),
			
			reflect.TypeOf((*pkg_.Event_RequestSigVerified)(nil)),
			
			reflect.TypeOf((*pkg_.Event_StoreVerifiedRequest)(nil)),
			
			reflect.TypeOf((*pkg_.Event_AppSnapshotRequest)(nil)),
			
			reflect.TypeOf((*pkg_.Event_AppSnapshot)(nil)),
			
			reflect.TypeOf((*pkg_.Event_AppRestoreState)(nil)),
			
			reflect.TypeOf((*pkg_.Event_TimerDelay)(nil)),
			
			reflect.TypeOf((*pkg_.Event_TimerRepeat)(nil)),
			
			reflect.TypeOf((*pkg_.Event_TimerGarbageCollect)(nil)),
			
			reflect.TypeOf((*pkg_.Event_Bcb)(nil)),
			
			reflect.TypeOf((*pkg_.Event_Mempool)(nil)),
			
			reflect.TypeOf((*pkg_.Event_Availability)(nil)),
			
			reflect.TypeOf((*pkg_.Event_NewEpoch)(nil)),
			
			reflect.TypeOf((*pkg_.Event_NewConfig)(nil)),
			
			reflect.TypeOf((*pkg_.Event_TestingString)(nil)),
			
			reflect.TypeOf((*pkg_.Event_TestingUint)(nil)),
			
			reflect.TypeOf((*pkg_.Init)(nil)),
			
			reflect.TypeOf((*pkg_.Tick)(nil)),
			
			reflect.TypeOf((*pkg_.NewRequests)(nil)),
			
			reflect.TypeOf((*pkg_.HashRequest)(nil)),
			
			reflect.TypeOf((*pkg_.HashResult)(nil)),
			
			reflect.TypeOf((*pkg_.HashOrigin)(nil)),
			
			reflect.TypeOf((*pkg_.HashOrigin_ContextStore)(nil)),
			
			reflect.TypeOf((*pkg_.HashOrigin_Request)(nil)),
			
			reflect.TypeOf((*pkg_.HashOrigin_Iss)(nil)),
			
			reflect.TypeOf((*pkg_.HashOrigin_Dsl)(nil)),
			
			reflect.TypeOf((*pkg_.SignRequest)(nil)),
			
			reflect.TypeOf((*pkg_.SignResult)(nil)),
			
			reflect.TypeOf((*pkg_.SignOrigin)(nil)),
			
			reflect.TypeOf((*pkg_.SignOrigin_ContextStore)(nil)),
			
			reflect.TypeOf((*pkg_.SignOrigin_Iss)(nil)),
			
			reflect.TypeOf((*pkg_.SignOrigin_Dsl)(nil)),
			
			reflect.TypeOf((*pkg_.SigVerData)(nil)),
			
			reflect.TypeOf((*pkg_.VerifyNodeSigs)(nil)),
			
			reflect.TypeOf((*pkg_.NodeSigsVerified)(nil)),
			
			reflect.TypeOf((*pkg_.SigVerOrigin)(nil)),
			
			reflect.TypeOf((*pkg_.SigVerOrigin_ContextStore)(nil)),
			
			reflect.TypeOf((*pkg_.SigVerOrigin_Iss)(nil)),
			
			reflect.TypeOf((*pkg_.SigVerOrigin_Dsl)(nil)),
			
			reflect.TypeOf((*pkg_.RequestReady)(nil)),
			
			reflect.TypeOf((*pkg_.SendMessage)(nil)),
			
			reflect.TypeOf((*pkg_.MessageReceived)(nil)),
			
			reflect.TypeOf((*pkg_.WALAppend)(nil)),
			
			reflect.TypeOf((*pkg_.WALEntry)(nil)),
			
			reflect.TypeOf((*pkg_.WALTruncate)(nil)),
			
			reflect.TypeOf((*pkg_.WALLoadAll)(nil)),
			
			reflect.TypeOf((*pkg_.Deliver)(nil)),
			
			reflect.TypeOf((*pkg_.VerifyRequestSig)(nil)),
			
			reflect.TypeOf((*pkg_.RequestSigVerified)(nil)),
			
			reflect.TypeOf((*pkg_.StoreVerifiedRequest)(nil)),
			
			reflect.TypeOf((*pkg_.AppSnapshotRequest)(nil)),
			
			reflect.TypeOf((*pkg_.AppSnapshot)(nil)),
			
			reflect.TypeOf((*pkg_.AppRestoreState)(nil)),
			
			reflect.TypeOf((*pkg_.TimerDelay)(nil)),
			
			reflect.TypeOf((*pkg_.TimerRepeat)(nil)),
			
			reflect.TypeOf((*pkg_.TimerGarbageCollect)(nil)),
			
			reflect.TypeOf((*pkg_.NewEpoch)(nil)),
			
			reflect.TypeOf((*pkg_.NewConfig)(nil)),
			
		})

	if err != nil {
		log.Printf("Error: %v\n", err)
		os.Exit(2)
	}
}
