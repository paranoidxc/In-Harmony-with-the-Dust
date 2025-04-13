package syntax

import (
	"errors"
	//"log/slog"
	"math"
)

type TrackingRuneIter struct {
	buf     *Buf
	pos     Pos
	curPos  Pos
	eof     bool
	limit   uint64
	numRead uint64
	maxRead *uint64
}

func NewTrackingRuneIter(pos Pos, buf *Buf) TrackingRuneIter {
	var maxRead uint64
	return TrackingRuneIter{
		pos:     pos,
		curPos:  pos,
		buf:     buf,
		limit:   math.MaxUint64,
		maxRead: &maxRead,
	}
}

func (iter *TrackingRuneIter) NextRune() (a rune, e error) {
	//slog.Info("iter NextRune ==========")
	if iter.curPos.Row < len(iter.buf.runes) {
		line := iter.buf.runes[iter.curPos.Row]
		if iter.curPos.Col < len(line) {
		} else {
			e = errors.New("eof1")
		}
	} else {
		e = errors.New("eof2")
	}

	if e == nil {
		a := iter.buf.runes[iter.curPos.Row][iter.curPos.Col]
		//slog.Info("iter NextRune", slog.Any("rune", string(a)))
		if iter.curPos.Col+1 == len(iter.buf.runes[iter.curPos.Row]) {
			iter.curPos.Col = 0
			iter.curPos.Row += 1
		} else {
			iter.curPos.Col += 1
		}
		return a, nil
	}

	e = errors.New("eof")
	return
}

func (iter *TrackingRuneIter) Skip(n uint64) uint64 {
	for i := uint64(0); i < n; i++ {
		_, err := iter.NextRune()
		if err != nil {
			return i
		}
	}
	return n
}

// Limit sets the maximum number of runes this reader can produce.
func (iter *TrackingRuneIter) Limit(n uint64) {
	iter.limit = n
}
