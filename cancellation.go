package lobsterdata

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// LOBSTERCancellation is a struct that represents the Order
// Cancellation event type from the LOBSTER data set.
type LOBSTERCancellation struct {
	EventSinceMidnight time.Duration `json:"timesincemidnight"`
	OrderID            uint64        `json:"orderid"`
	Size               uint64        `json:"size"`
	Price              uint64        `json:"price"`
	Direction          int64         `json:"side"`
}

// UnmarshalCsvLOBSTER unmarshals a list of strings into a
// LOBSTERCancellation, given they are parsed from encoding/csv.
func (lc *LOBSTERCancellation) UnmarshalCsvLOBSTER(eventFields []string) (err error) {
	if len(eventFields) != 6 {
		err = fmt.Errorf("Error unmarshalling LOBSTER line, data does not have 6 columns")
		return
	}

	if Event(eventFields[1]) != Cancellation {
		err = fmt.Errorf("Trying to unmarshal a LOBSTER cancellation from an event that is not a cancellation is invalid")
		return
	}

	// Adding a "seconds" to the first field because we want to parse
	// it as a duration
	eventFields[0] += "s"
	if lc.EventSinceMidnight, err = time.ParseDuration(eventFields[0]); err != nil {
		err = fmt.Errorf("Error parsing the time field in LOBSTER data as a duration: %s", err)
		return
	}

	if lc.OrderID, err = strconv.ParseUint(eventFields[2], 10, 64); err != nil {
		err = fmt.Errorf("Error parsing orderid field in LOBSTER data as uint64: %s", err)
		return
	}

	if lc.Size, err = strconv.ParseUint(eventFields[3], 10, 64); err != nil {
		err = fmt.Errorf("Error parsing size field in LOBSTER data as uint64: %s", err)
		return
	}

	if lc.Price, err = strconv.ParseUint(eventFields[4], 10, 64); err != nil {
		err = fmt.Errorf("Error parsing price field in LOBSTER data as uint64: %s", err)
		return
	}

	if lc.Direction, err = strconv.ParseInt(eventFields[5], 10, 64); err != nil {
		err = fmt.Errorf("Error parsing direction field in LOBSTER data as int64: %s", err)
		return
	}
	return
}

// MarshalLOBSTER marshals a LOBSTERCancellation into a set of strings
// that can be written using encoding/csv.
func (lc *LOBSTERCancellation) MarshalCsvLOBSTER() (eventFields []string, err error) {
	eventFields = make([]string, 6)
	eventFields[0] = fmt.Sprintf("%f", lc.EventSinceMidnight.Seconds())
	eventFields[1] = fmt.Sprintf("%s", Cancellation)
	eventFields[2] = fmt.Sprintf("%d", lc.OrderID)
	eventFields[3] = fmt.Sprintf("%d", lc.Size)
	eventFields[4] = fmt.Sprintf("%d", lc.Price)
	eventFields[5] = fmt.Sprintf("%d", lc.Direction)
	return
}

// MarshalJSON implements the JSONMarshaler interface for this struct.
func (lc *LOBSTERCancellation) MarshalJSON() (jsonBytes []byte, err error) {
	return json.Marshal(struct {
		TheMainEvent LOBSTERCancellation `json:"event"`
		EventType    Event               `json:"eventtype"`
	}{
		TheMainEvent: LOBSTERCancellation(*lc),
		EventType:    Cancellation,
	})
}
