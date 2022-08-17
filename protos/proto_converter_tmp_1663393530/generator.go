
package main

import (
	"log"
	"os"
	"reflect"

	generator_ "github.com/filecoin-project/mir/codegen/generators/mir-std-gen/generator"
	pkg_ "github.com/filecoin-project/mir/pkg/pb/eventpb"
)

func main() {
	generator := generator_.CombinedGenerator{}
	err := generator.Run(
		[]reflect.Type{
			
			reflect.TypeOf((*pkg_.Event)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.Event_Init)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.Event_Tick)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.Event_WalAppend)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.Event_WalEntry)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.Event_WalTruncate)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.Event_NewRequests)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.Event_HashRequest)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.Event_HashResult)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.Event_SignRequest)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.Event_SignResult)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.Event_VerifyNodeSigs)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.Event_NodeSigsVerified)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.Event_RequestReady)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.Event_SendMessage)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.Event_MessageReceived)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.Event_DeliverCert)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.Event_Iss)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.Event_VerifyRequestSig)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.Event_RequestSigVerified)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.Event_StoreVerifiedRequest)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.Event_AppSnapshotRequest)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.Event_AppSnapshot)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.Event_AppRestoreState)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.Event_TimerDelay)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.Event_TimerRepeat)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.Event_TimerGarbageCollect)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.Event_Bcb)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.Event_Mempool)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.Event_Availability)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.Event_NewEpoch)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.Event_NewConfig)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.Event_Factory)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.Event_BatchDb)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.Event_BatchFetcher)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.Event_ThreshCrypto)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.Event_PingPong)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.Event_Checkpoint)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.Event_SbEvent)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.Event_TestingString)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.Event_TestingUint)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.Init)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.Tick)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.NewRequests)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.HashRequest)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.HashResult)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.HashOrigin)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.HashOrigin_ContextStore)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.HashOrigin_Request)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.HashOrigin_Iss)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.HashOrigin_Dsl)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.HashOrigin_Checkpoint)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.HashOrigin_Sb)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.SignRequest)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.SignResult)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.SignOrigin)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.SignOrigin_ContextStore)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.SignOrigin_Dsl)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.SignOrigin_Checkpoint)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.SignOrigin_Sb)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.SigVerData)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.VerifyNodeSigs)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.NodeSigsVerified)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.SigVerOrigin)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.SigVerOrigin_ContextStore)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.SigVerOrigin_Iss)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.SigVerOrigin_Dsl)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.SigVerOrigin_Checkpoint)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.SigVerOrigin_Sb)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.RequestReady)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.SendMessage)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.MessageReceived)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.WALAppend)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.WALEntry)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.WALTruncate)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.WALLoadAll)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.DeliverCert)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.VerifyRequestSig)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.RequestSigVerified)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.StoreVerifiedRequest)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.AppSnapshotRequest)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.AppSnapshot)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.AppRestoreState)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.TimerDelay)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.TimerRepeat)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.TimerGarbageCollect)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.NewEpoch)(nil)).Elem(),
			
			reflect.TypeOf((*pkg_.NewConfig)(nil)).Elem(),
			
		})

	if err != nil {
		log.Printf("Error: %v\n", err)
		os.Exit(2)
	}
}
