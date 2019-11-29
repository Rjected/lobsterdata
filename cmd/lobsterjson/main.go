package main

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"os"

	"github.com/op/go-logging"
	"github.com/rjected/lobsterdata"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app         = kingpin.New("lobsterjson", "A LOBSTER data csv to json tool.")
	verbose     = app.Flag("verbose", "Verbose mode.").Short('v').Bool()
	lobsterpath = app.Flag("path", "Path to LOBSTER csv file").Required().File()
	lobsterout  = app.Flag("output", "Path to output json file").String()
	tostdout    = app.Flag("tostdout", "Send JSON to standard output.").Bool()
	numrows     = app.Flag("numrows", "Number of rows to process.").Uint()

	log = logging.MustGetLogger("lobsterdata")
	// Example format string. Everything except the message has a custom color
	// which is dependent on the log level. Many fields have a custom output
	// formatting too, eg. the time returns the hour down to the milli second.
	format = logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)
)

// TODO: add filtering, for example a flag that says "only process
// entries with this type

func main() {
	app.HelpFlag.Short('h')
	app.Parse(os.Args[1:])

	// For demo purposes, create two backend for os.Stderr.
	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2 := logging.NewLogBackend(os.Stderr, "", 0)

	// For messages written to backend2 we want to add some additional
	// information to the output, including the used log level and the name of
	// the function.
	backend2Formatter := logging.NewBackendFormatter(backend2, format)

	// Only errors and more severe messages should be sent to backend1
	backend1Leveled := logging.AddModuleLevel(backend1)
	backend1Leveled.SetLevel(logging.ERROR, "")

	// Set the backends to be used.
	logging.SetBackend(backend1Leveled, backend2Formatter)

	if verbose != nil && *verbose {
		logging.SetLevel(logging.INFO, "lobsterdata")
	}

	if len(os.Args[1:]) < 2 {
		log.Critical("Need to provide at least two arguments")
		return
	}

	if lobsterout == nil && tostdout == nil {
		log.Critical("Must provide either a path for the json output file, or a flag for standard output")
		return
	}

	if lobsterpath == nil {
		log.Critical("Error, the lobster data file path was not specified, please specify a path")
		return
	}

	var jsonoutput io.Writer
	var jsonfile *os.File
	var err error
	if lobsterout != nil && tostdout == nil {
		log.Info("Creating output json file")
		outfilename := *lobsterout
		if jsonfile, err = os.Create(outfilename); err != nil {
			log.Criticalf("Could not create output JSON file: %s", err)
			return
		}
		jsonoutput = jsonfile
	}

	if tostdout != nil {
		log.Info("Setting json output at stdout")
		jsonoutput = os.Stdout
	}

	var actualPath *os.File = *lobsterpath
	var csvLine []string
	var jsonLine []byte
	csvReader := csv.NewReader(actualPath)
	lineNum := uint(0)

	// initialize all these pointers, we might need em
	currSubmission := new(lobsterdata.LOBSTERSubmission)
	currCancellation := new(lobsterdata.LOBSTERCancellation)
	currDeletion := new(lobsterdata.LOBSTERDeletion)
	currExecutionVisible := new(lobsterdata.LOBSTERExecutionVisible)
	currExecutionHidden := new(lobsterdata.LOBSTERExecutionHidden)
	currCrossTrade := new(lobsterdata.LOBSTERCrossTrade)
	currTradingHalt := new(lobsterdata.LOBSTERTradingHalt)

	log.Info("Starting CSV read")
	for csvLine, err = csvReader.Read(); err == nil; csvLine, err = csvReader.Read() {
		lineNum++
		switch lobsterdata.Event(csvLine[1]) {
		case lobsterdata.Submission:
			if err = currSubmission.UnmarshalCsvLOBSTER(csvLine); err != nil {
				log.Criticalf("Error unmarshalling csv line into submission: %s", err)
				return
			}
			if jsonLine, err = json.MarshalIndent(currSubmission, "", "\t"); err != nil {
				log.Criticalf("Error marshalling data into json: %s", err)
				return
			}
		case lobsterdata.Cancellation:
			if err = currCancellation.UnmarshalCsvLOBSTER(csvLine); err != nil {
				log.Criticalf("Error unmarshalling csv line into cancellation: %s", err)
				return
			}
			if jsonLine, err = json.MarshalIndent(currCancellation, "", "\t"); err != nil {
				log.Criticalf("Error marshalling data into json: %s", err)
				return
			}
		case lobsterdata.Deletion:
			if err = currDeletion.UnmarshalCsvLOBSTER(csvLine); err != nil {
				log.Criticalf("Error unmarshalling csv line into Deletion: %s", err)
				return
			}
			if jsonLine, err = json.MarshalIndent(currDeletion, "", "\t"); err != nil {
				log.Criticalf("Error marshalling data into json: %s", err)
				return
			}
		case lobsterdata.ExecutionVisible:
			if err = currExecutionVisible.UnmarshalCsvLOBSTER(csvLine); err != nil {
				log.Criticalf("Error unmarshalling csv line into ExecutionVisible: %s", err)
				return
			}
			if jsonLine, err = json.MarshalIndent(currExecutionVisible, "", "\t"); err != nil {
				log.Criticalf("Error marshalling data into json: %s", err)
				return
			}
		case lobsterdata.ExecutionHidden:
			if err = currExecutionHidden.UnmarshalCsvLOBSTER(csvLine); err != nil {
				log.Criticalf("Error unmarshalling csv line into ExecutionHidden: %s", err)
				return
			}
			if jsonLine, err = json.MarshalIndent(currExecutionHidden, "", "\t"); err != nil {
				log.Criticalf("Error marshalling data into json: %s", err)
				return
			}
		case lobsterdata.CrossTrade:
			if err = currCrossTrade.UnmarshalCsvLOBSTER(csvLine); err != nil {
				log.Criticalf("Error unmarshalling csv line into CrossTrade: %s", err)
				return
			}
			if jsonLine, err = json.MarshalIndent(currCrossTrade, "", "\t"); err != nil {
				log.Criticalf("Error marshalling data into json: %s", err)
				return
			}
		case lobsterdata.TradingHalt:
			if err = currTradingHalt.UnmarshalCsvLOBSTER(csvLine); err != nil {
				log.Criticalf("Error unmarshalling csv line into TradingHalt: %s", err)
				return
			}
			if jsonLine, err = json.MarshalIndent(currTradingHalt, "", "\t"); err != nil {
				log.Criticalf("Error marshalling data into json: %s", err)
				return
			}
		default:
			log.Error("Encountered invalid data in csv file on line %d", lineNum)
			continue
		}

		// there should be something in jsonLine now, also don't worry
		// about the number of bytes for now
		if _, err = jsonoutput.Write(jsonLine); err != nil {
			log.Criticalf("Error writing json to output: %s", err)
			return
		}

		if numrows != nil && lineNum == *numrows {
			log.Info("Done processing data!")
			break
		}
	}

	if err != nil && err != io.EOF {
		log.Criticalf("Error after reading csv file is not EOF: %s", err)
		return
	}

	log.Info("Done reading csv, closing csv file")
	if err = actualPath.Close(); err != nil {
		log.Criticalf("Error closing csv file: %s", err)
		return
	}

	if jsonfile != nil && tostdout == nil {
		log.Info("Done writing to json file, closing")
		if err = jsonfile.Close(); err != nil {
			log.Criticalf("Error closing json file after writing: %s", err)
			return
		}
	}

	return
}
