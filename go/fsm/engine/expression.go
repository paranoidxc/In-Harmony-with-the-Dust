package engine

// Event is an input event.
// This usually represents a keypress, but the compiled state machine doesn't assume
// that the events have any particular meaning.
type Event int64

// CaptureId is an identifier for a subsequence of events.
type CaptureId uint64

// Expr is a regular expression that matches input events.
type Expr interface{}
