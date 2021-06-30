package main

import (
	"fmt"
	"log"
	"net"
	"os"

	config "github.com/chutommy/commodity-prices/config"
	data "github.com/chutommy/commodity-prices/data"
	commodity "github.com/chutommy/commodity-prices/protos/commodity"
	server "github.com/chutommy/commodity-prices/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	// define logger
	l := log.New(os.Stdout, "[COMMODITY SERVICE] ", log.LstdFlags)

	// get config from the file
	cfg, err := config.GetConfig("config.yaml")
	if err != nil {
		l.Fatalf("[ERROR] get configuration: %v", err)
	}

	// data service
	ds := data.New()
	err = ds.Update()
	if err != nil {
		l.Fatalf("[ERROR] can not update data: %v", err)
	}

	// service server
	cmdSrv := server.New(l, ds)
	go cmdSrv.HandleUpdates()

	// grpc server
	grpcSrv := grpc.NewServer()

	// register
	commodity.RegisterCommodityServer(grpcSrv, cmdSrv)
	reflection.Register(grpcSrv) // support reflection

	// define listen
	lst, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
	if err != nil {
		l.Fatalf("[ERROR] define listening: %v", err)
	}

	// start listening
	l.Printf("[SUCCESS] Listening on %s:%d", cfg.Host, cfg.Port)
	err = grpcSrv.Serve(lst)
	if err != nil {
		l.Panicf("listening: %v", err)
	}
}
