package lobsterdata

// HaltReason represents the reason for a trading halt.
// Its values are defined according to the corresponding Price field
// in the LOBSTER csv.
type HaltReason int64

const (
	HaltTrading   HaltReason = -1
	ResumeQuoting HaltReason = 0
	ResumeTrading HaltReason = 1
)
