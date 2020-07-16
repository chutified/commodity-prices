package data

import (
	"testing"
	"time"
)

var expCommodities int = 39

func TestCommoditiesData(t *testing.T) {

	// New >>>>>>>>>>>>>>>
	cd := New()

	// Update >>>>>>>>>>>>>>>
	err := cd.Update()
	if err != nil {
		t.Fatalf("unexpacted source error: %v", err)
	}
	if l := len(cd.Commodities); l != expCommodities {
		t.Errorf("expected %d of commodities data, got %d", expCommodities, l)
	}

	// MonitorData >>>>>>>>>>>>>>>
	updCh, errCh := cd.MonitorData(1 * time.Second)

	err = nil
	go func() {
		err = <-errCh
	}()
	_ = <-updCh

	if err != nil {
		t.Errorf("unexpected err occured: %v", err)
	}
}
