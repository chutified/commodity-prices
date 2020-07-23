package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"
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

	return c
}

// GetCommodity handles grpc calls.
func (c *Commodities) GetCommodity(ctx context.Context, req *commodity.CommodityRequest) (*commodity.CommodityResponse, error) {

	// handling
	resp, err := c.handleRequest(req)
	if err != nil {
		c.log.Printf("[ERROR] handle request data: %v", err)

		gErr := status.Newf(
			codes.NotFound,
			"Name of the commodity \"%s\" was not found.", req.GetName(),
		)
		gErr, wde := gErr.WithDetails(req)
		if wde != nil {
			return nil, err
		}

		return nil, gErr.Err()
	}

	// success
	c.log.Printf("[SUCCESS] respond to the client's GetCommodity request: %s", req.GetName())
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
			c.log.Printf("[ERROR] invalid request format: %v", err)
			return err
		}

		// validate if error occurs, terminate the request
		n := req.GetName()
		if _, ok := c.data.Commodities[n]; !ok {
			c.log.Printf("[ERROR] commodity %s not found", n)

			// handle grpc error
			gStatus := status.Newf(
				codes.NotFound,
				"Commodity %s was not found.", n,
			)
			gStatus.WithDetails(req)
			srv.Send(&commodity.StreamingCommodityResponse{
				Message: &commodity.StreamingCommodityResponse_Error{
					Error: gStatus.Proto(),
				},
			})
			continue
		}

		// append a subscribtion
		reqs, ok := c.subscribtions[srv]
		if !ok {
			c.subscribtions[srv] = []*commodity.CommodityRequest{}
		}

		// skip if already subscribed
		var validErr *status.Status
		for _, v := range reqs {
			if v.GetName() == req.GetName() {
				c.log.Printf("[ERROR] client already subscribe for: %s", v.GetName())

				// err
				validErr = status.Newf(
					codes.AlreadyExists,
					"Client is already subscribed for commodity:\"%s\"", v.GetName(),
				)
				// add original request
				var wde error
				validErr, wde = validErr.WithDetails(req)
				if wde != nil {
					c.log.Printf("[ERROR] unable to add metadata to error: %v", wde)
				}
				break
			}
		}
		// if validErr exists rettuns error and continue
		if validErr != nil {
			srv.Send(&commodity.StreamingCommodityResponse{
				Message: &commodity.StreamingCommodityResponse_Error{
					Error: validErr.Proto(),
				},
			})
			continue
		}

		// success
		c.subscribtions[srv] = append(reqs, req)
		c.log.Printf("[SUCCESS] client subscribed: %v", req)
	}

	// break
	return nil
}

// HandleUpdates post for each subscribed client the update
func (c *Commodities) HandleUpdates() {

	// range updates
	updates, errs := c.data.MonitorData(15 * time.Second)

	// errors
	go func() {
		for {
			c.log.Printf("[ERROR] %v", <-errs)
		}
	}()

	// updates
	for range updates {

		// inform
		c.log.Printf("[UPDATE] send new values to subscribers")

		// loop over subscribed clients
		for clientSrv, reqs := range c.subscribtions {

			// loop over client's requests
			for _, req := range reqs {

				// handling
				resp, err := c.handleRequest(req)
				if err != nil {
					c.log.Printf("[ERROR] handle request data: %v", err)

					// handle grpc error
					gErr := status.Newf(
						codes.NotFound,
						"Name of the commodity \"%s\" was not found.", req.GetName(),
					)
					gErr, wde := gErr.WithDetails(req)
					if wde != nil {
						c.log.Printf("[ERROR] unable to add metadata to error: %v", wde)
					}
					err = clientSrv.Send(&commodity.StreamingCommodityResponse{
						Message: &commodity.StreamingCommodityResponse_Error{
							Error: gErr.Proto(),
						},
					})
					if err != nil {
						c.log.Printf("[ERROR] send request data: %v", err)
						continue
					}

					continue
				}

				// post response
				err = clientSrv.Send(&commodity.StreamingCommodityResponse{
					Message: &commodity.StreamingCommodityResponse_CommodityResponse{
						CommodityResponse: resp,
					},
				})
				if err != nil {
					c.log.Printf("[ERROR] send request data: %v", err)
					continue
				}
			}
		}
	}
}

// handleRequest handles the request and returns the appropriate response.
func (c *Commodities) handleRequest(req *commodity.CommodityRequest) (*commodity.CommodityResponse, error) {

	// get the commodity
	name := req.GetName()
	name = strings.ToLower(name)
	cmd, err := c.data.GetCommodity(name)
	if err != nil {
		return nil, fmt.Errorf("get commodity: %w", err)
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
	fmt.Println("HERE", resp.Currency)
	return resp, nil
}
