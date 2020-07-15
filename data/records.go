package data

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chutified/resource-finder/models"
)

// record defines the raw data from the website.
type record struct {
	name  string
	price string
	chp   string
	chn   string
	unit  string
	time  string
}

// getRecords crawls the website and provides the records.
func getRecords() ([]record, error) {

	// get the root selector
	doc, err := goquery.NewDocument("https://markets.businessinsider.com/commodities")
	if err != nil {
		return nil, fmt.Errorf("goquery getting  HTML response: %w", err)
	}

	// list of records
	var records []record
	// query the commodities rows
	doc.Find(".row-hover").Each(func(i int, s *goquery.Selection) {
		children := s.Children()
		if children.HasClass("tdFixed") {

			// defines values
			name := children.First()
			price := name.Next()
			chp := price.Next()
			chn := chp.Next()
			unit := chn.Next()
			time := unit.Next()

			// save the record
			r := record{
				name:  name.Text(),
				price: price.Text(),
				chp:   chp.Text(),
				chn:   chn.Text(),
				unit:  unit.Text(),
				time:  time.Text(),
			}
			records = append(records, r)
		}
	})

	// success
	return records, nil
}

// parseRecords parses the records and returns the slice of commodities.
func parseRecords(rs []record) ([]models.Commodity, error) {

	// list of commodities
	var cmds []models.Commodity
	// parse each record
	for _, r := range rs {

		// parse name
		n := strings.TrimSpace(r.name)
		name := strings.ToLower(n) // Name

		// parse price
		pStr := strings.TrimSpace(r.price)
		pStr = strings.ReplaceAll(pStr, ",", "")
		priceF64, err := strconv.ParseFloat(pStr, 32)
		if err != nil {
			return nil, fmt.Errorf("parsing commodity price: %w", err)
		}
		price := float32(priceF64) // Price

		// parse chp
		chpStr := strings.ReplaceAll(r.chp, "%", "")
		chpStr = strings.TrimSpace(chpStr)
		chpF64, err := strconv.ParseFloat(chpStr, 32)
		if err != nil {
			return nil, fmt.Errorf("parsing commodity percentages change: %w", err)
		}
		changeP := float32(chpF64) // ChangeP

		// parse chn
		chnStr := strings.ReplaceAll(r.chn, "%", "")
		chnStr = strings.TrimSpace(chnStr)
		chnF64, err := strconv.ParseFloat(chnStr, 32)
		if err != nil {
			return nil, fmt.Errorf("parsing commodity number change: %w", err)
		}
		changeN := float32(chnF64) // ChangeN

		// parse unit
		ss := strings.Split(r.unit, "per")
		if len(ss) != 2 {
			return nil, fmt.Errorf("unexpected unit format: \"%s\"", r.unit)
		}
		currency := strings.TrimSpace(ss[0]) // Currency
		wu := strings.TrimSpace(ss[1])
		weightUnit := strings.ToLower(wu) // WeightUnit

		// define time layouts
		const timeTime = "3:04:05 PM"
		const timeDate = "1/2/2006"
		// parse time
		t, err := time.Parse(timeTime, r.time)
		if err != nil {
			t, err = time.Parse(timeDate, r.time)
			if err != nil {
				return nil, fmt.Errorf("unexpected time format: %w", err)
			}
		} else {
			now := time.Now()
			t = time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.UTC)
		}
		lastUpdate := t // LastUpdate

		// append to the other commodities
		cmd := models.Commodity{
			Name:       name,
			Price:      price,
			Currency:   currency,
			WeightUnit: weightUnit,
			ChangeP:    changeP,
			ChangeN:    changeN,
			LastUpdate: lastUpdate,
		}
		cmds = append(cmds, cmd)
	}

	// success
	return cmds, nil
}
