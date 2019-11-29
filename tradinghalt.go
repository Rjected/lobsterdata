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
func (lth *LOBSTERTradingHalt) UnmarshalCsvLOBSTER(eventFields []string) (err error) {
	if len(eventFields) != 6 {
		err = fmt.Errorf("Error unmarshalling LOBSTER line, data does not have 6 columns")
		return
	}

	if Event(eventFields[1]) != TradingHalt {
		err = fmt.Errorf("Trying to unmarshal a LOBSTER trading halt from an event that is not a trading halt is invalid")
		return
	}

	// Adding a "seconds" to the first field because we want to parse
	// it as a duration
	eventFields[0] += "s"
	if lth.EventSinceMidnight, err = time.ParseDuration(eventFields[0]); err != nil {
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
	lth.HaltType = HaltReason(tmpHaltType)

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
func (lth *LOBSTERTradingHalt) MarshalCsvLOBSTER() (eventFields []string, err error) {
	eventFields = make([]string, 6)
	eventFields[0] = fmt.Sprintf("%f", lth.EventSinceMidnight.Seconds())
	eventFields[1] = fmt.Sprintf("%s", TradingHalt)
	eventFields[2] = "0"
	eventFields[3] = "0"
	eventFields[4] = fmt.Sprintf("%d", lth.HaltType)
	eventFields[5] = "-1"
	return
}
