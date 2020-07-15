package data

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/chutified/resource-finder/models"
)

// CommoditiesData is a data controller.
type CommoditiesData struct {
	log         *log.Logger
	Commodities map[string]models.Commodity
}

// New constructs a new data service.
func New(l *log.Logger) *CommoditiesData {
	cd := &CommoditiesData{
		log:         l,
		Commodities: make(map[string]models.Commodity),
	}

	return cd
}

// Update updates the commodities data.
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
func (cd *CommoditiesData) MonitorData(interval time.Duration) chan struct{} {

	// channel for the notification
	ret := make(chan struct{})

	go func() {
		ticker := time.NewTicker(interval)
		cache := make(map[string]models.Commodity)

		// check every tick
		for range ticker.C {

			// update
			err := cd.Update()
			if err != nil {
				cd.log.Printf("[ERROR] updating data: %v", err)
			}

			// compare
			if !(reflect.DeepEqual(cache, cd.Commodities)) {

				// update cache
				for k, v := range cd.Commodities {
					cache[k] = v
				}

				// inform
				cd.log.Printf("Data updated.")
				ret <- struct{}{}
			}
		}
	}()

	return ret
}
