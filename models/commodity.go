package models

import "time"

// Comodity defines and holds a commodity data values.
type Commodity struct {
	Name       string    // commodity name
	Price      float32   // current commodity price
	Currency   string    // currency of the price
	WeightUnit string    // weight unit for the price
	ChangeP    float32   // last change in percentages
	ChangeN    float32   // last change in a number
	LastUpdate time.Time // updated at
}
