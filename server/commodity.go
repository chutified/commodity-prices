package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	data "github.com/chutified/resource-finder/data"
	commodity "github.com/chutified/resource-finder/protos/commodity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Commodities defines the commodity server.
type Commodities struct {
	log           *log.Logger
	data          *data.CommoditiesData
	subscribtions map[commodity.Commodity_SubscribeCommodityServer][]*commodity.CommodityRequest
}

// New creates a new commodity server.
func New(l *log.Logger, cd *data.CommoditiesData) *Commodities {
	c := &Commodities{
		log:           l,
		data:          cd,
		subscribtions: make(map[commodity.Commodity_SubscribeCommodityServer][]*commodity.CommodityRequest),
	}

	go c.handleUpdates()

	return c
}

// handleUpdates post for each subscribed client the update
func (c *Commodities) handleUpdates() {

	// range updates
	updates, errs := c.data.MonitorData(45 * time.Second)

	// errors
	go func() {
		for {
			c.log.Printf("[ERROR] %v", <-errs)
		}
	}()

	// updates
	for range updates {

		// inform
		c.log.Printf("[UPDATE] new values")

		// loop over subscribed clients
		for clientSrv, reqs := range c.subscribtions {

			// loop over client's requests
			for _, req := range reqs {

				// handling
				resp, err := c.handleRequest(req)
				if err != nil {
					c.log.Printf("[ERROR] handling request data: %v", err)
					continue
				}

				// post response
				err = clientSrv.Send(resp)
				if err != nil {
					c.log.Printf("[ERROR] sending request data: %v", err)
					continue
				}
			}
		}
	}
}

// GetCommodity handles grpc calls.
func (c *Commodities) GetCommodity(ctx context.Context, req *commodity.CommodityRequest) (*commodity.CommodityResponse, error) {

	// handling
	resp, err := c.handleRequest(req)
	if err != nil {
		c.log.Printf("[ERROR] handling request data: %v", err)

		grpcErr := status.Newf(
			codes.NotFound,
			"Name of the commodity \"%s\" was not found.", req.GetName(),
		)
		grpcErr, wde := grpcErr.WithDetails(req)
		if wde != nil {
			return nil, err
		}

		return nil, grpcErr.Err()
	}

	// success
	c.log.Printf("[HANDLE] client request: %v", req)
	return resp, nil
}

// SubscribeCommodity handles grpc subscription.
func (c *Commodities) SubscribeCommodity(srv commodity.Commodity_SubscribeCommodityServer) error { // satisfy CommodityServer interface

	// handling requests
	for {

		// get request
		req, err := srv.Recv()
		if err == io.EOF {
			c.log.Printf("[EXIT] client closed connection")
			break
		}
		if err != nil {
			c.log.Printf("[ERROR] invalid request: %v", err)
			return err // TODO handle error
		}

		// validate
		n := req.GetName()
		if _, ok := c.data.Commodities[n]; !ok {
			c.log.Printf("[ERROR] commodity %s not found", n)
			continue
		}

		// success
		c.log.Printf("[HANDLE] client subscribtion: %v", req)

		// append a subscribtion
		reqs, ok := c.subscribtions[srv]
		if !ok {
			c.subscribtions[srv] = []*commodity.CommodityRequest{}
		}

		subscribed := false
		for _, r := range reqs {
			if r.GetName() == req.GetName() {
				subscribed = true
				break
			}
		}

		// skip if already subscribed
		if !subscribed {
			c.subscribtions[srv] = append(reqs, req)
		}
	}

	// break
	return nil
}

// handleRequest handles the request and returns the appropriate response.
func (c *Commodities) handleRequest(req *commodity.CommodityRequest) (*commodity.CommodityResponse, error) {

	// search
	name := req.GetName()
	cmd, ok := c.data.Commodities[name]
	if !ok {
		return nil, fmt.Errorf("commodity %s not found", name)
	}

	// success
	resp := &commodity.CommodityResponse{
		Name:       cmd.Name,
		Price:      cmd.Price,
		Currency:   cmd.Currency,
		WeightUnit: cmd.WeightUnit,
		ChangeP:    cmd.ChangeP,
		ChangeN:    cmd.ChangeN,
		LastUpdate: cmd.LastUpdate.Unix(),
	}
	return resp, nil
}
