package syntax

type Func func(TrackingRuneIter, State) Result

type Result struct {
	NumConsumed    uint64
	ReadLength     uint64
	ComputedTokens []ComputedToken
	NextState      State
}

type Pos struct {
	Row int
	Col int
}

var FailedResult = Result{}

func (r Result) IsSuccess() bool {
	return r.NumConsumed > 0
}

func (r Result) IsFailure() bool {
	return !r.IsSuccess()
}

func (r Result) ShiftForward(n uint64) Result {
	if n > 0 {
		r.NumConsumed += n
		for i := 0; i < len(r.ComputedTokens); i++ {
			r.ComputedTokens[i].Offset += n
		}
	}
	return r
}

type computation struct {
	readLength     uint64
	consumedLength uint64
	startState     State
	endState       State
	tokens         []ComputedToken
}

type Edit struct {
	Offset      uint64
	NumInserted uint64
	NumDeleted  uint64
}

type P struct {
	parseFunc    Func
	computations []computation
}

func New(f Func) *P {
	f = f.recoverFromFailure()
	return &P{parseFunc: f}
}

const minInitialConsumedLen = 1024

func (p *P) ParseAll(buf *Buf) {
	state := State(EmptyState{})
	var totalOffset uint64
	pos := Pos{Row: 0, Col: 0}
	n := totalChars(buf)

	p.computations = p.computations[:0]

	for totalOffset < n {
		result := p.runParseFunc(pos, buf, state)
		if result.NumConsumed == 0 {
			break
		}

		c := computation{
			readLength:     result.ReadLength,
			consumedLength: result.NumConsumed,
			startState:     state,
			endState:       result.NextState,
			tokens:         result.ComputedTokens,
		}

		if len(p.computations) > 0 {
			last := &p.computations[len(p.computations)-1]
			if last.consumedLength < minInitialConsumedLen {
				for _, tok := range c.tokens {
					last.tokens = append(last.tokens, ComputedToken{
						Offset: last.consumedLength + tok.Offset,
						Length: tok.Length,
						Role:   tok.Role,
					})
				}
				readLength := last.consumedLength + c.readLength
				if readLength > last.readLength {
					last.readLength = readLength
				}
				last.consumedLength += c.consumedLength
				last.endState = c.endState
			} else {
				p.computations = append(p.computations, c)
			}
		} else {
			p.computations = append(p.computations, c)
		}

		totalOffset += result.NumConsumed
		state = result.NextState
		pos = advancePos(buf, pos, result.NumConsumed)
	}
}

// Tokens returns all tokens with absolute offsets from the document start.
func (p *P) Tokens() []ComputedToken {
	var tokens []ComputedToken
	var offset uint64
	for _, c := range p.computations {
		for _, tok := range c.tokens {
			tokens = append(tokens, ComputedToken{
				Offset: offset + tok.Offset,
				Length: tok.Length,
				Role:   tok.Role,
			})
		}
		offset += c.consumedLength
	}
	return tokens
}

func (p *P) ParseAfterEdit(buf *Buf, edit Edit) {
	oldComputations := append([]computation(nil), p.computations...)
	oldOffsetByIndex := make([]uint64, len(oldComputations))
	var oldOffset uint64
	for i, c := range oldComputations {
		oldOffsetByIndex[i] = oldOffset
		oldOffset += c.consumedLength
	}

	p.computations = p.computations[:0]
	state := State(EmptyState{})
	pos := Pos{Row: 0, Col: 0}
	var totalOffset uint64
	n := totalChars(buf)

	for totalOffset < n {
		if c, ok := reusableComputation(oldComputations, oldOffsetByIndex, edit, totalOffset, state); ok {
			p.computations = append(p.computations, c)
			totalOffset += c.consumedLength
			state = c.endState
			pos = advancePos(buf, pos, c.consumedLength)
			continue
		}

		result := p.runParseFunc(pos, buf, state)
		if result.NumConsumed == 0 {
			break
		}

		c := computation{
			readLength:     result.ReadLength,
			consumedLength: result.NumConsumed,
			startState:     state,
			endState:       result.NextState,
			tokens:         result.ComputedTokens,
		}
		p.computations = append(p.computations, c)
		totalOffset += result.NumConsumed
		state = result.NextState
		pos = advancePos(buf, pos, result.NumConsumed)
	}
}

func reusableComputation(computations []computation, offsets []uint64, edit Edit, newOffset uint64, state State) (computation, bool) {
	oldOffset, ok := oldOffsetAfterEdit(edit, newOffset)
	if !ok {
		return computation{}, false
	}

	for i, c := range computations {
		if offsets[i] != oldOffset || c.startState != state {
			continue
		}
		if computationOverlapsEdit(offsets[i], c, edit) {
			return computation{}, false
		}
		return c, true
	}
	return computation{}, false
}

func oldOffsetAfterEdit(edit Edit, newOffset uint64) (uint64, bool) {
	if newOffset < edit.Offset {
		return newOffset, true
	}
	if newOffset < edit.Offset+edit.NumInserted {
		return 0, false
	}
	return newOffset - edit.NumInserted + edit.NumDeleted, true
}

func computationOverlapsEdit(offset uint64, c computation, edit Edit) bool {
	readEnd := offset + c.readLength
	deleteEnd := edit.Offset + edit.NumDeleted
	if edit.NumDeleted > 0 && offset < deleteEnd && readEnd > edit.Offset {
		return true
	}
	if edit.NumInserted > 0 && offset < edit.Offset && readEnd > edit.Offset {
		return true
	}
	return false
}

func (p *P) runParseFunc(pos Pos, buf *Buf, state State) Result {
	trackingIter := NewTrackingRuneIter(pos, buf)
	result := p.parseFunc(trackingIter, state)
	result.ReadLength = trackingIter.MaxRead()
	if result.ReadLength < result.NumConsumed {
		result.ReadLength = result.NumConsumed
	}
	return result
}

func totalChars(buf *Buf) uint64 {
	if len(buf.runes) == 0 {
		return 0
	}

	var n uint64
	for _, line := range buf.runes {
		n += uint64(len(line))
	}
	return n + uint64(len(buf.runes)-1)
}

func advancePos(buf *Buf, pos Pos, n uint64) Pos {
	row := pos.Row
	col := pos.Col
	remaining := n

	for remaining > 0 && row < len(buf.runes) {
		lineLen := len(buf.runes[row])
		if col < lineLen {
			available := uint64(lineLen - col)
			if remaining <= available {
				col += int(remaining)
				break
			}
			remaining -= available
			col = lineLen
		}

		if remaining > 0 {
			if row >= len(buf.runes)-1 {
				break
			}
			remaining--
			row++
			col = 0
		}
	}
	return Pos{Row: row, Col: col}
}
