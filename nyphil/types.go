package nyphil

import "errors"

var (
	YearInds  = loadYears()
	NoDataErr = errors.New("Not enough data to write file")
)

//create a composer type
type composer struct {
	Nprog    int
	Programs []uint16
	Seasons  []string
}


