package lobsterdata

type HaltReason int64

const (
	HaltTrading   HaltReason = -1
	ResumeQuoting HaltReason = 0
	ResumeTrading HaltReason = 1
)
