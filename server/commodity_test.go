package server

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"testing"

	data "github.com/chutified/resource-finder/data"
	commodity "github.com/chutified/resource-finder/protos/commodity"
	"gopkg.in/go-playground/assert.v1"
)

func TestCommodities(t *testing.T) {

	l := log.New(bytes.NewBuffer([]byte{}), "", log.LUTC)
	cd := data.New()
	err := cd.Update()
	if err != nil {
		t.Fatalf("unexpacted source error: %v", err)
	}

	// New >>>>>>>>>>>>>>>
	cs := New(l, cd)

	// GetCommodity >>>>>>>>>>>>>>>
	// ok
	{
		req := &commodity.CommodityRequest{
			Name: "coal",
		}
		_, err := cs.GetCommodity(context.Background(), req)

		assert.Equal(t, err, nil)
	}
	// commodity not found
	{
		req := &commodity.CommodityRequest{
			Name: "invalid",
		}
		_, err := cs.GetCommodity(context.Background(), req)

		exp := fmt.Sprintf(".*%s.*", "Name of the commodity .* was not found.")
		if err == nil {
			t.Fatal("expected error not found")
		}
		assert.MatchRegex(t, err.Error(), exp)
	}

	// SubscribeCommodity >>>>>>>>>>>>>>>
	// TODO

	// HandleUpdates()
	// TODO
}
