package lobsterdata

import (
	"fmt"
	"strconv"
	"time"
)

// LOBSTERExecutionHidden is a struct that represents the hidden
// execution event type from the LOBSTER data set.
type LOBSTERExecutionHidden struct {
	EventSinceMidnight time.Duration `json:"timesincemidnight"`
	Size               uint64        `json:"size"`
	Price              uint64        `json:"price"`
	Direction          int64         `json:"side"`
}

// UnmarshalCsvLOBSTER unmarshals a list of strings into a
// LOBSTERExecutionHidden, given they are parsed from encoding/csv.
func (ls *LOBSTERExecutionHidden) UnmarshalCsvLOBSTER(eventFields []string) (err error) {
	if len(eventFields) != 6 {
		err = fmt.Errorf("Error unmarshalling LOBSTER submission, data does not have 6 columns")
		return
	}

	eventType, err := strconv.ParseUint(eventFields[1], 10, 8)
	if err != nil {
		err = fmt.Errorf("Error parsing eventType field in LOBSTER data as uint8: %s", err)
		return
	}
	if eventType != 1 {
		err = fmt.Errorf("Trying to unmarshal a LOBSTER submission from an event that is not a submission is invalid")
		return
	}

	// Adding a "seconds" to the first field because we want to parse
	// it as a duration
	eventFields[0] += "s"
	if ls.EventSinceMidnight, err = time.ParseDuration(eventFields[0]); err != nil {
		err = fmt.Errorf("Error parsing the time field in LOBSTER data as a duration: %s", err)
		return
	}

	var tmpOrderID uint64
	if tmpOrderID, err = strconv.ParseUint(eventFields[2], 10, 64); err != nil {
		err = fmt.Errorf("Error parsing orderid field in LOBSTER data as uint64: %s", err)
		return
	}
	if tmpOrderID != 0 {
		err = fmt.Errorf("Error, orderid field on hidden executions should be 0")
		return
	}

	if ls.Size, err = strconv.ParseUint(eventFields[3], 10, 64); err != nil {
		err = fmt.Errorf("Error parsing size field in LOBSTER data as uint64: %s", err)
		return
	}

	if ls.Price, err = strconv.ParseUint(eventFields[4], 10, 64); err != nil {
		err = fmt.Errorf("Error parsing price field in LOBSTER data as uint64: %s", err)
		return
	}

	if ls.Direction, err = strconv.ParseInt(eventFields[5], 10, 64); err != nil {
		err = fmt.Errorf("Error parsing direction field in LOBSTER data as int64: %s", err)
		return
	}
	return
}

// MarshalLOBSTER marshals a LOBSTERExecutionHidden into a set of strings
// that can be written using encoding/csv.
func (ls *LOBSTERExecutionHidden) MarshalCsvLOBSTER() (eventFields []string, err error) {
	eventFields = make([]string, 6)
	eventFields[0] = fmt.Sprintf("%f", ls.EventSinceMidnight.Seconds())
	eventFields[1] = "1"
	eventFields[2] = "0"
	eventFields[3] = fmt.Sprintf("%d", ls.Size)
	eventFields[4] = fmt.Sprintf("%d", ls.Price)
	eventFields[5] = fmt.Sprintf("%d", ls.Direction)
	return
}
