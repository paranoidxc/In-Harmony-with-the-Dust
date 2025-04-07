package syntax

import (
	"log/slog"
)

type Func func(TrackingRuneIter, State) Result

type Result struct {
	NumConsumed uint64
	NextState   State
}

type Pos struct {
	Row int
	Col int
}

var FailedResult = Result{}

// IsSuccess returns whether the parse succeeded.
func (r Result) IsSuccess() bool {
	return r.NumConsumed > 0
}

// IsFailure returns whether the parse failed.
func (r Result) IsFailure() bool {
	return !r.IsSuccess()
}

// ShiftForward shifts the result offsets forward by the specified number of positions.
func (r Result) ShiftForward(n uint64) Result {
	if n > 0 {
		r.NumConsumed += n
		// for i := 0; i < len(r.ComputedTokens); i++ {
		// 	r.ComputedTokens[i].Offset += n
		// }
	}
	return r
}

type P struct {
	parseFunc Func
	//lastComputation *computation
}

func New(f Func) *P {
	// This ensures that the parse func always makes progress.
	f = f.recoverFromFailure()
	return &P{parseFunc: f}
}

// ParseAll parses the entire document.
// func (p *P) ParseAll(buf *mgr.Buf) {
func (p *P) ParseAll(buf *Buf) {
	// var prevComputation *computation
	state := State(EmptyState{})
	// leafComputations := make([]*computation, 0)
	// n := tree.NumChars()
	slog.Info("ParseAll")
	pos := Pos{}
	//for pos.Row <= tree.Rows {
	//for pos.Col <= tree.RowCols {
	p.runParseFunc(pos, buf, state)
	//slog.Info("runParseFunc c", slog.Any("c", c))
	//p.runParseFunc(tree, pos, state)
	//slog.Info("break")
	//break
	//}
	//break
	// 	pos += c.ConsumedLength()
	// 	state = c.EndState()
	//
	// 	if prevComputation != nil && prevComputation.ConsumedLength() < minInitialConsumedLen {
	// 		// For the initial parse, combine small leaves. This saves memory by reducing both
	// 		// the number of leaves and parent nodes we need to allocate.
	// 		combineLeaves(prevComputation, c)
	// 	} else {
	// 		leafComputations = append(leafComputations, c)
	// 		prevComputation = c
	// 	}
	//}
	// c := concatLeafComputations(leafComputations)
	// p.lastComputation = c
}

func (p *P) runParseFunc(pos Pos, buf *Buf, state State) {
	trackingIter := NewTrackingRuneIter(pos, buf)
	p.parseFunc(trackingIter, state)
}
