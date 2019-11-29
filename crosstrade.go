package lobsterdata

import (
	"fmt"
	"strconv"
	"time"
)

// LOBSTERCrossTrade is a struct that represents the cross trade event
// type from the LOBSTER data set.
type LOBSTERCrossTrade struct {
	EventSinceMidnight time.Duration `json:"timesincemidnight"`
	OrderID            uint64        `json:"orderid"`
	Size               uint64        `json:"size"`
	Price              uint64        `json:"price"`
	Direction          int64         `json:"side"`
}

// UnmarshalCsvLOBSTER unmarshals a list of strings into a
// LOBSTERCrossTrade, given they are parsed from encoding/csv.
func (lt *LOBSTERCrossTrade) UnmarshalCsvLOBSTER(eventFields []string) (err error) {
	if len(eventFields) != 6 {
		err = fmt.Errorf("Error unmarshalling LOBSTER line, data does not have 6 columns")
		return
	}

	if Event(eventFields[1]) != CrossTrade {
		err = fmt.Errorf("Trying to unmarshal a LOBSTER cross trade from an event that is not a cross trade is invalid")
		return
	}

	// Adding a "seconds" to the first field because we want to parse
	// it as a duration
	eventFields[0] += "s"
	if lt.EventSinceMidnight, err = time.ParseDuration(eventFields[0]); err != nil {
		err = fmt.Errorf("Error parsing the time field in LOBSTER data as a duration: %s", err)
		return
	}

	if lt.OrderID, err = strconv.ParseUint(eventFields[2], 10, 64); err != nil {
		err = fmt.Errorf("Error parsing orderid field in LOBSTER data as uint64: %s", err)
		return
	}

	if lt.Size, err = strconv.ParseUint(eventFields[3], 10, 64); err != nil {
		err = fmt.Errorf("Error parsing size field in LOBSTER data as uint64: %s", err)
		return
	}

	if lt.Price, err = strconv.ParseUint(eventFields[4], 10, 64); err != nil {
		err = fmt.Errorf("Error parsing price field in LOBSTER data as uint64: %s", err)
		return
	}

	if lt.Direction, err = strconv.ParseInt(eventFields[5], 10, 64); err != nil {
		err = fmt.Errorf("Error parsing direction field in LOBSTER data as int64: %s", err)
		return
	}
	return
}

// MarshalLOBSTER marshals a LOBSTERCrossTrade into a set of strings
// that can be written using encoding/csv.
func (lt *LOBSTERCrossTrade) MarshalCsvLOBSTER() (eventFields []string, err error) {
	eventFields = make([]string, 6)
	eventFields[0] = fmt.Sprintf("%f", lt.EventSinceMidnight.Seconds())
	eventFields[1] = fmt.Sprintf("%s", CrossTrade)
	eventFields[2] = fmt.Sprintf("%d", lt.OrderID)
	eventFields[3] = fmt.Sprintf("%d", lt.Size)
	eventFields[4] = fmt.Sprintf("%d", lt.Price)
	eventFields[5] = fmt.Sprintf("%d", lt.Direction)
	return
}
