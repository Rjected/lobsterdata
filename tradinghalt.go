package lobsterdata

import (
	"fmt"
	"strconv"
	"time"
)

// LOBSTERTradingHalt is a struct that represents the trading halt
// event type from the LOBSTER data set.
type LOBSTERTradingHalt struct {
	EventSinceMidnight time.Duration `json:"timesincemidnight"`
	HaltType           HaltReason    `json:"halttype"`
}

// UnmarshalCsvLOBSTER unmarshals a list of strings into a
// LOBSTERTradingHalt, given they are parsed from encoding/csv.
func (ls *LOBSTERTradingHalt) UnmarshalCsvLOBSTER(eventFields []string) (err error) {
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
		err = fmt.Errorf("Invalid trading halt, orderid field must be 0")
		return
	}

	var tmpSize uint64
	if tmpSize, err = strconv.ParseUint(eventFields[3], 10, 64); err != nil {
		err = fmt.Errorf("Error parsing size field in LOBSTER data as uint64: %s", err)
		return
	}
	if tmpSize != 0 {
		err = fmt.Errorf("Invalid trading halt, size field must be 0")
		return
	}

	var tmpHaltType int64
	if tmpHaltType, err = strconv.ParseInt(eventFields[4], 10, 64); err != nil {
		err = fmt.Errorf("Error parsing price field in LOBSTER data as int64: %s", err)
		return
	}
	ls.HaltType = HaltReason(tmpHaltType)

	var tmpDirection int64
	if tmpDirection, err = strconv.ParseInt(eventFields[5], 10, 64); err != nil {
		err = fmt.Errorf("Error parsing direction field in LOBSTER data as int64: %s", err)
		return
	}
	if tmpDirection != -1 {
		err = fmt.Errorf("Invalid trading halt, direction field must be -1")
		return
	}
	return
}

// MarshalLOBSTER marshals a LOBSTERTradingHalt into a set of strings
// that can be written using encoding/csv.
func (ls *LOBSTERTradingHalt) MarshalCsvLOBSTER() (eventFields []string, err error) {
	eventFields = make([]string, 6)
	eventFields[0] = fmt.Sprintf("%f", ls.EventSinceMidnight.Seconds())
	eventFields[1] = "7"
	eventFields[2] = "0"
	eventFields[3] = "0"
	eventFields[4] = fmt.Sprintf("%d", ls.HaltType)
	eventFields[5] = "-1"
	return
}
