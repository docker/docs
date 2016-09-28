package api

const (
	// DefaultPerPageLimit is the default number of results to return in a
	// paginated request when a limit has not been specified or is
	// unparseable.
	DefaultPerPageLimit = 10
	// MaxPerPageLimit is the maximum number of results to return in a
	// paginated request if a specified limit exceeds this amount.
	MaxPerPageLimit = 1 << 16
)
