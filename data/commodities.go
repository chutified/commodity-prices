package data

import (
	"fmt"
	"reflect"
	"time"

	models "github.com/chutified/commodity-prices/models"
)

// CommoditiesData is a data controller.
type CommoditiesData struct {
	Commodities map[string]models.Commodity
}

// New constructs a new data service.
func New() *CommoditiesData {
	cd := &CommoditiesData{
		Commodities: make(map[string]models.Commodity),
	}
	return cd
}

// GetCommodity retrieves the commodtiy from the memory.
func (cd *CommoditiesData) GetCommodity(name string) (*models.Commodity, error) {

	// search
	cmd, ok := cd.Commodities[name]
	if !ok {
		return nil, fmt.Errorf("commodity %s not found", name)
	}

	// success
	return &cmd, nil
}

// Update updates the Commodities data.
func (cd *CommoditiesData) Update() error {

	cmds, err := getCommodities()
	if err != nil {
		return fmt.Errorf("fetching data: %w", err)
	}

	// success
	cd.Commodities = cmds
	return nil
}

// MonitorData returns a channel which can notify if any data modification occurs.
func (cd *CommoditiesData) MonitorData(interval time.Duration) (chan struct{}, chan error) {

	// channel for the notification
	updCh := make(chan struct{})
	errCh := make(chan error)

	go func() {
		ticker := time.NewTicker(interval)
		cache := make(map[string]models.Commodity)

		// check every tick
		for range ticker.C {

			// update
			err := cd.Update()
			if err != nil {
				errCh <- fmt.Errorf("[ERROR] updating data: %w", err)
			}

			// compare
			if !(reflect.DeepEqual(cache, cd.Commodities)) {

				// update cache
				for k, v := range cd.Commodities {
					cache[k] = v
				}

				// inform
				updCh <- struct{}{}
			}
		}
	}()

	return updCh, errCh
}
