package syntax

import (
	"errors"
	"math"
)

type TrackingRuneIter struct {
	buf     *Buf
	curPos  Pos
	limit   uint64
	numRead uint64
	maxRead *uint64
}

func NewTrackingRuneIter(pos Pos, buf *Buf) TrackingRuneIter {
	var maxRead uint64
	return TrackingRuneIter{
		curPos:  pos,
		buf:     buf,
		limit:   math.MaxUint64,
		maxRead: &maxRead,
	}
}

func (iter *TrackingRuneIter) NextRune() (a rune, e error) {
	if iter.numRead >= iter.limit || iter.curPos.Row >= len(iter.buf.runes) {
		return 0, errors.New("eof")
	}

	line := iter.buf.runes[iter.curPos.Row]
	switch {
	case iter.curPos.Col < len(line):
		a = line[iter.curPos.Col]
		iter.curPos.Col++
	case iter.curPos.Col == len(line) && iter.curPos.Row < len(iter.buf.runes)-1:
		a = '\n'
		iter.curPos.Row++
		iter.curPos.Col = 0
	default:
		return 0, errors.New("eof")
	}

	iter.numRead++
	if iter.maxRead != nil && iter.numRead > *iter.maxRead {
		*iter.maxRead = iter.numRead
	}
	return a, nil
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

func (iter *TrackingRuneIter) MaxRead() uint64 {
	if iter.maxRead == nil {
		return iter.numRead
	}
	return *iter.maxRead
}

// Limit sets the maximum number of runes this reader can produce.
func (iter *TrackingRuneIter) Limit(n uint64) {
	iter.limit = iter.numRead + n
}
