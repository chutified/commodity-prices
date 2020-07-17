.PHONY: clean, protogen, build, run

protogen:
	protoc -I protos/ --go_out=plugins=grpc:protos/commodity/ protos/commodity.proto

build:
	docker build -t commodity_prices .

run:
	docker run -it --network="host" --name commoditysrv --rm commodity_prices
