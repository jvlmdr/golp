package lp

// The result of a pivot operation.
type Pivot struct {
	Enter     int
	Leave     int
	Final     bool
	Unbounded bool
}
