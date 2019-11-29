package lobsterdata

// LOBSTERData is an interface for all types of LOBSTER Data, it
// specifies that LOBSTER Data should marshal and unmarshal from an
// output of a csv
type LOBSTERData interface {
	// UnmarshalCsvLOBSTER should unmarshal data from a list of
	// strings, which are assumed to be received by the encoding/csv
	// output from parsing a LOBSTER dataset csv.
	UnmarshalCsvLOBSTER([]string) error

	// MarshalCsvLOBSTER should marshal data from a LOBSTERData,
	// outputting a list of strings that can be written to a csv file
	// by encoding/csv.
	MarshalCsvLOBSTER() ([]string, error)
}
