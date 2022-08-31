/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package protos

// Build the custom code generators.
//go:generate go build -o ../codegen/proto-converter/proto-converter.bin ../codegen/proto-converter
//go:generate go build -o ../codegen/protoc-plugin/protoc-gen-mir ../codegen/protoc-plugin

// Define some helpful shorthands.
//go:generate -command protoc-events protoc --proto_path=. --go_out=../pkg/pb/ --go_opt=paths=source_relative --plugin=../codegen/protoc-plugin/protoc-gen-mir --mir_out=../pkg/pb --mir_opt=paths=source_relative
//go:generate -command proto-converter ../codegen/proto-converter/proto-converter.bin

//go:generate protoc-events mir/plugin.proto
//go:generate protoc-events commonpb/commonpb.proto
//go:generate protoc-events messagepb/messagepb.proto
//go:generate protoc-events requestpb/requestpb.proto
//go:generate protoc-events eventpb/eventpb.proto
//go:generate protoc-events recordingpb/recordingpb.proto
//go:generate protoc-events isspb/isspb.proto
//go:generate protoc-events bcbpb/bcbpb.proto
//go:generate protoc-events isspbftpb/isspbftpb.proto
//go:generate protoc-events contextstorepb/contextstorepb.proto
//go:generate protoc-events dslpb/dslpb.proto
//go:generate protoc-events mempoolpb/mempoolpb.proto
//go:generate protoc-events availabilitypb/availabilitypb.proto
//go:generate protoc-events availabilitypb/mscpb/mscpb.proto

//go:generate proto-converter "github.com/filecoin-project/mir/pkg/pb/eventpb"
//go:generate proto-converter "github.com/filecoin-project/mir/pkg/pb/contextstorepb"
//go:generate proto-converter "github.com/filecoin-project/mir/pkg/pb/availabilitypb"

//go:generate protoc --proto_path=. --go_out=:../pkg/ --go_opt=paths=source_relative simplewal/simplewal.proto
//go:generate protoc --proto_path=. --go_out=:../samples/ --go_opt=paths=source_relative chat-demo/chatdemo.proto
//go:generate protoc --go_out=../pkg/ --go_opt=paths=source_relative --go-grpc_out=../pkg/ --go-grpc_opt=paths=source_relative requestreceiver/requestreceiver.proto
//go:generate protoc --go_out=../pkg/ --go_opt=paths=source_relative --go-grpc_out=../pkg/ --go-grpc_opt=paths=source_relative net/grpc/grpctransport.proto
//xgo:generate protoc --proto_path=. --go_out=plugins=grpc:../pkg/ --go_opt=paths=source_relative grpctransport/grpctransport.proto
//xgo:generate protoc --proto_path=. --go_out=plugins=grpc:../pkg/ --go_opt=paths=source_relative requestreceiver/requestreceiver.proto
