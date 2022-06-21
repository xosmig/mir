package protoutil

//go:generate protoc --proto_path=. --go_out=testpb/ --go_opt=paths=source_relative protoutil_testpb.proto
