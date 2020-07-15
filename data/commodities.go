package data

import (
	"fmt"

	"github.com/chutified/resource-finder/models"
)

// CommoditiesData is a data controller.
type CommoditiesData struct {
	Commodities map[string]models.Commodity
}

// New constructs a new data service.
func New() *CommoditiesData {
	return &CommoditiesData{
		Commodities: make(map[string]models.Commodity),
	}
}

// Update updates the commodities data.
func (s *CommoditiesData) Update() error {

	// get current records
	rs, err := getRecords()
	if err != nil {
		return fmt.Errorf("getting records: %w", err)
	}

	// parse into slice of commodities
	cmds, err := parseRecords(rs)
	if err != nil {
		return err
	}

	// parsu into map of commodities
	mcmds := mapCommodities(cmds)

	// success
	s.Commodities = mcmds
	return nil
}

// mapCommodities maps each commodity with its name.
func mapCommodities(cmds []models.Commodity) map[string]models.Commodity {
	m := make(map[string]models.Commodity)
	for _, c := range cmds {
		m[c.Name] = c
	}
	return m
}
