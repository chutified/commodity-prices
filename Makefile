.PHONY: protogen

protogen:
	protoc -I protos/ --go_out=plugins=grpc:protos/commodity/ protos/commodity.proto
