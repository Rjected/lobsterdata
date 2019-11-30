package lobsterdata

import (
	"encoding/json"
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
func (lh *LOBSTERExecutionHidden) UnmarshalCsvLOBSTER(eventFields []string) (err error) {
	if len(eventFields) != 6 {
		err = fmt.Errorf("Error unmarshalling LOBSTER line, data does not have 6 columns")
		return
	}

	if Event(eventFields[1]) != ExecutionHidden {
		err = fmt.Errorf("Trying to unmarshal a LOBSTER hidden execution from an event that is not a hidden execution is invalid")
		return
	}

	// Adding a "seconds" to the first field because we want to parse
	// it as a duration
	eventFields[0] += "s"
	if lh.EventSinceMidnight, err = time.ParseDuration(eventFields[0]); err != nil {
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

	if lh.Size, err = strconv.ParseUint(eventFields[3], 10, 64); err != nil {
		err = fmt.Errorf("Error parsing size field in LOBSTER data as uint64: %s", err)
		return
	}

	if lh.Price, err = strconv.ParseUint(eventFields[4], 10, 64); err != nil {
		err = fmt.Errorf("Error parsing price field in LOBSTER data as uint64: %s", err)
		return
	}

	if lh.Direction, err = strconv.ParseInt(eventFields[5], 10, 64); err != nil {
		err = fmt.Errorf("Error parsing direction field in LOBSTER data as int64: %s", err)
		return
	}
	return
}

// MarshalLOBSTER marshals a LOBSTERExecutionHidden into a set of strings
// that can be written using encoding/csv.
func (lh *LOBSTERExecutionHidden) MarshalCsvLOBSTER() (eventFields []string, err error) {
	eventFields = make([]string, 6)
	eventFields[0] = fmt.Sprintf("%f", lh.EventSinceMidnight.Seconds())
	eventFields[1] = fmt.Sprintf("%s", ExecutionHidden)
	eventFields[2] = "0"
	eventFields[3] = fmt.Sprintf("%d", lh.Size)
	eventFields[4] = fmt.Sprintf("%d", lh.Price)
	eventFields[5] = fmt.Sprintf("%d", lh.Direction)
	return
}

// MarshalJSON implements the JSONMarshaler interface for this struct.
func (lh *LOBSTERExecutionHidden) MarshalJSON() (jsonBytes []byte, err error) {
	return json.Marshal(struct {
		TheMainEvent LOBSTERExecutionHidden `json:"event"`
		EventType    Event                  `json:"eventtype"`
	}{
		TheMainEvent: LOBSTERExecutionHidden(*lh),
		EventType:    ExecutionHidden,
	})
}
