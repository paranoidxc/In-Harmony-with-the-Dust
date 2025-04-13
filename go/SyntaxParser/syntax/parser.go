package syntax

import (
	"log/slog"
)

type Func func(TrackingRuneIter, State) Result

type Result struct {
	NumConsumed    uint64
	ComputedTokens []ComputedToken
	NextState      State
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
	//slog.Info(">>>> ShiftForward", slog.Any("n", n))
	if n > 0 {
		r.NumConsumed += n
		// for i := 0; i < len(r.ComputedTokens); i++ {
		// 	r.ComputedTokens[i].Offset += n
		// }
	}

	//slog.Info(">>>> ShiftForward", slog.Any("shiftNum", n), slog.Any("result", r))
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
	slog.Info("ParseAll", slog.Any("line_count", len(buf.runes)))
	pos := Pos{}

	for pos.Row < len(buf.runes) {
		for pos.Col < len(buf.runes[pos.Row]) {
			result := p.runParseFunc(pos, buf, state)

			// slog.Info(">>>>> runParseFunc result >>>>>",
			// 	slog.Any("source pos", pos),
			// 	slog.Any("result", result),
			// )
			startPos := pos
			endPos := pos

			leftNum := -1
			numConsumed := int(result.NumConsumed)

			for leftNum == -1 {
				line := buf.runes[pos.Row]
				lineLen := len(line)
				newLineCol := (pos.Col + numConsumed)
				if lineLen-newLineCol == 0 {
					endPos.Col = lineLen - 1
					pos.Col = 0
					pos.Row += 1
					break
				} else if lineLen-newLineCol > 0 {
					pos.Col += numConsumed
					endPos.Col = pos.Col - 1
					break
				} else {
					numConsumed -= (lineLen - pos.Col)

					pos.Row += 1
					pos.Col = 0

					if numConsumed <= 0 {
						endPos.Row = pos.Row - 1
						endPos.Col = newLineCol - 1
						break
					}
				}
			}

			if result.ComputedTokens != nil {
				slog.Info(">>>>>>>>>>>> Syntax Found <<<<<<<<<<<<",
					slog.Any("result", result),
					slog.Any("start pos", startPos),
					slog.Any("end pos", endPos),
					slog.Any("new pos", pos),
				)
			}

			//slog.Info("runParseFunc new pos >>>>> ", slog.Any("new pos", pos))
			if pos.Col == 0 {
				break
			}
		}

		if pos.Row < len(buf.runes) && len(buf.runes[pos.Row]) == 0 {
			pos.Row += 1
		}
	}
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

func (p *P) runParseFunc(pos Pos, buf *Buf, state State) Result {
	trackingIter := NewTrackingRuneIter(pos, buf)
	result := p.parseFunc(trackingIter, state)
	return result
}
