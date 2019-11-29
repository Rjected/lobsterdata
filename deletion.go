package lobsterdata

import (
	"fmt"
	"strconv"
	"time"
)

// LOBSTERDeletion is a struct that represents the Order Deletion
// event type from the LOBSTER data set.
type LOBSTERDeletion struct {
	EventSinceMidnight time.Duration `json:"timesincemidnight"`
	OrderID            uint64        `json:"orderid"`
	Size               uint64        `json:"size"`
	Price              uint64        `json:"price"`
	Direction          int64         `json:"side"`
}

// UnmarshalCsvLOBSTER unmarshals a list of strings into a
// LOBSTERDeletion, given they are parsed from encoding/csv.
func (ld *LOBSTERDeletion) UnmarshalCsvLOBSTER(eventFields []string) (err error) {
	if len(eventFields) != 6 {
		err = fmt.Errorf("Error unmarshalling LOBSTER line, data does not have 6 columns")
		return
	}

	if Event(eventFields[1]) != Deletion {
		err = fmt.Errorf("Trying to unmarshal a LOBSTER deletion from an event that is not a deletion is invalid")
		return
	}

	// Adding a "seconds" to the first field because we want to parse
	// it as a duration
	eventFields[0] += "s"
	if ld.EventSinceMidnight, err = time.ParseDuration(eventFields[0]); err != nil {
		err = fmt.Errorf("Error parsing the time field in LOBSTER data as a duration: %s", err)
		return
	}

	if ld.OrderID, err = strconv.ParseUint(eventFields[2], 10, 64); err != nil {
		err = fmt.Errorf("Error parsing orderid field in LOBSTER data as uint64: %s", err)
		return
	}

	if ld.Size, err = strconv.ParseUint(eventFields[3], 10, 64); err != nil {
		err = fmt.Errorf("Error parsing size field in LOBSTER data as uint64: %s", err)
		return
	}

	if ld.Price, err = strconv.ParseUint(eventFields[4], 10, 64); err != nil {
		err = fmt.Errorf("Error parsing price field in LOBSTER data as uint64: %s", err)
		return
	}

	if ld.Direction, err = strconv.ParseInt(eventFields[5], 10, 64); err != nil {
		err = fmt.Errorf("Error parsing direction field in LOBSTER data as int64: %s", err)
		return
	}
	return
}

// MarshalLOBSTER marshals a LOBSTERDeletion into a set of strings
// that can be written using encoding/csv.
func (ld *LOBSTERDeletion) MarshalCsvLOBSTER() (eventFields []string, err error) {
	eventFields = make([]string, 6)
	eventFields[0] = fmt.Sprintf("%f", ld.EventSinceMidnight.Seconds())
	eventFields[1] = fmt.Sprintf("%s", Deletion)
	eventFields[2] = fmt.Sprintf("%d", ld.OrderID)
	eventFields[3] = fmt.Sprintf("%d", ld.Size)
	eventFields[4] = fmt.Sprintf("%d", ld.Price)
	eventFields[5] = fmt.Sprintf("%d", ld.Direction)
	return
}
