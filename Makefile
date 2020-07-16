.PHONY: clean, protogen, build, run

protogen:
	protoc -I protos/ --go_out=plugins=grpc:protos/commodity/ protos/commodity.proto

build:
	docker build -t commodity_prices .

run:
	docker run -it -p 10501:10501 --name commoditysrv --rm commodity_prices
