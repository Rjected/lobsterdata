package lobsterdata

// Event represents what type of event the lobster row is.
type Event string

const (
	Submission       Event = "1"
	Cancellation     Event = "2"
	Deletion         Event = "3"
	ExecutionVisible Event = "4"
	ExecutionHidden  Event = "5"
	CrossTrade       Event = "6"
	TradingHalt      Event = "7"
)
