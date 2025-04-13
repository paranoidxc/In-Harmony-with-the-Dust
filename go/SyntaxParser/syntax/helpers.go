package syntax

import (
	//"log/slog"
	"sort"
	"strings"
	"unicode/utf8"
)

// consumeString consumes the characters in `s`.
func consumeString(s string) Func {
	return func(iter TrackingRuneIter, state State) Result {
		var numConsumed uint64
		for _, targetRune := range s {
			r, err := iter.NextRune()
			if err != nil || r != targetRune {
				return FailedResult
			}
			numConsumed++
		}
		return Result{
			NumConsumed: numConsumed,
			NextState:   state,
		}
	}
}

// consumeToString consumes all characters up to and including the string `s`.
func consumeToString(s string) Func {
	f := consumeString(s)
	return func(iter TrackingRuneIter, state State) Result {
		var numSkipped uint64
		for {
			r := f(iter, state)
			if r.IsSuccess() {
				return r.ShiftForward(numSkipped)
			}

			_, err := iter.NextRune()
			if err != nil {
				return FailedResult
			}
			numSkipped++
		}
	}
}

func consumeSingleRuneLike(predicateFn func(rune) bool) Func {
	return func(iter TrackingRuneIter, state State) Result {
		r, err := iter.NextRune()
		if err == nil && predicateFn(r) {
			return Result{
				NumConsumed: 1,
				NextState:   state,
			}
		}
		return FailedResult
	}
}

func consumeRunesLike(predicateFn func(rune) bool) Func {
	return func(iter TrackingRuneIter, state State) Result {
		var numConsumed uint64
		for {
			r, err := iter.NextRune()
			if err != nil || !predicateFn(r) {
				return Result{
					NumConsumed: numConsumed,
					NextState:   state,
				}
			}
			numConsumed++
		}
	}
}

// consumeToEofOrRuneLike consumes up to and including a rune matching a predicate or EOF.
func consumeToEofOrRuneLike(predicate func(r rune) bool) Func {
	return func(iter TrackingRuneIter, state State) Result {
		line := iter.buf.runes[iter.curPos.Row]
		lineLen := len(line)
		numConsumed := uint64(lineLen - iter.curPos.Col)
		/*
			for {
				r, err := iter.NextRune()
				if err == io.EOF {
					break
				} else if err != nil {
					return FailedResult
				}

				numConsumed++

				if predicate(r) {
					break
				}
			}
		*/

		// slog.Info(">>>>>>>>>>>>>>>>>>> consumeToEofOrRuneLike <<<<<<<<<<<<<<<<<<<<< ",
		// 	slog.Any("numConsumed", numConsumed),
		// )

		return Result{
			NumConsumed: numConsumed,
			NextState:   state,
		}
	}
}

// consumeToNextLineFeed consumes up to and including the next newline character or the last character in the document, whichever comes first.
var consumeToNextLineFeed = consumeToEofOrRuneLike(func(r rune) bool {
	return r == '\n'
})

func consumeDigitsAndSeparators(allowLeadingSeparator bool, isDigit func(r rune) bool) Func {
	return func(iter TrackingRuneIter, state State) Result {
		var numConsumed uint64
		var lastWasUnderscore bool
		for {
			r, err := iter.NextRune()
			if err != nil {
				break
			}

			if r == '_' && !lastWasUnderscore && (allowLeadingSeparator || numConsumed > 0) {
				lastWasUnderscore = true
				numConsumed++
				continue
			}

			if isDigit(r) {
				lastWasUnderscore = false
				numConsumed++
				continue
			}

			break
		}

		if lastWasUnderscore {
			numConsumed--
		}

		return Result{
			NumConsumed: numConsumed,
			NextState:   state,
		}
	}
}

// recognizeToken recognizes the consumed characters in the result as a token.
func recognizeToken(tokenRole TokenRole) MapFn {
	return func(result Result) Result {
		token := ComputedToken{
			Length: result.NumConsumed,
			Role:   tokenRole,
		}
		return Result{
			NumConsumed:    result.NumConsumed,
			ComputedTokens: []ComputedToken{token},
			NextState:      result.NextState,
		}
	}
}

func maxStrLen(ss []string) uint64 {
	maxLength := uint64(0)
	for _, s := range ss {
		length := uint64(utf8.RuneCountInString(s))
		if length > maxLength {
			maxLength = length
		}
	}
	return maxLength
}

// consumeLongestMatchingOption consumes the longest matching option from a set of options.
func consumeLongestMatchingOption(options []string) Func {
	// Sort options descending by length.
	sort.SliceStable(options, func(i, j int) bool {
		return len(options[i]) > len(options[j])
	})

	// Allocate buffer for lookahead runes (shared across func invocations).
	buf := make([]rune, maxStrLen(options))
	return func(iter TrackingRuneIter, state State) Result {
		// Lookahead up to the length of the longest option.
		var n uint64
		for i := 0; i < len(buf); i++ {
			r, err := iter.NextRune()
			if err != nil {
				break
			}
			buf[i] = r
			n++
		}

		// Look for longest matching option.
		// We can return the first one that matches b/c options
		// are sorted descending by length.
		for _, opt := range options {
			var i uint64
			matched := true
			for _, r := range opt {
				if r != buf[i] || i >= n {
					matched = false
					break
				}
				i++
			}
			if matched {

				// slog.Info(">>>>>>>>>>>>>>>>>>> consumeLongestMatchingOption <<<<<<<<<<<<<<<<<<<<< ",
				// 	slog.Any("numConsumed", i),
				// )

				return Result{
					NumConsumed: i,
					NextState:   state,
				}
			}
		}
		return FailedResult
	}
}

// readInputString reads a string from the text up to `n` characters long.
func readInputString(iter TrackingRuneIter, n uint64) string {
	var sb strings.Builder
	for i := uint64(0); i < n; i++ {
		r, err := iter.NextRune()
		if err != nil {
			break
		}
		if _, err := sb.WriteRune(r); err != nil {
			panic(err)
		}
	}
	return sb.String()
}

func recognizeKeywordOrConsume(keywords []string) MapWithInputFn {
	// Calculate the length of the longest keyword to limit how much
	// of the input needs to be reprocessed.
	maxLength := maxStrLen(keywords)
	return func(result Result, iter TrackingRuneIter, state State) Result {
		if result.NumConsumed > maxLength {
			return result
		}

		s := readInputString(iter, result.NumConsumed)

		// slog.Info("recognizeKeywordOrConsume----------------",
		// 	slog.Any("input string", s),
		// )

		for _, kw := range keywords {
			if kw == s {
				token := ComputedToken{
					Role:   TokenRoleKeyword,
					Length: result.NumConsumed,
				}

				result := Result{
					NumConsumed:    result.NumConsumed,
					ComputedTokens: []ComputedToken{token},
					NextState:      state,
				}
				// slog.Info(">>>>>>>>>>>>>>>>>>> recognizeKeywordOrConsume <<<<<<<<<<<<<<<<<<<<< ",
				// 	slog.Any("keyword", kw),
				// 	slog.Any("token", token),
				// 	slog.Any("result", result),
				// )
				return result
			}
		}

		return result
	}
}

// consumeCStyleString consumes a string with characters escaped by a backslash.
func consumeCStyleString(quoteRune rune, allowLineBreaks bool) Func {
	return func(iter TrackingRuneIter, state State) Result {
		var n uint64
		r, err := iter.NextRune()
		if err != nil || r != quoteRune {
			return FailedResult
		}
		n++

		var inEscapeSeq bool
		for {
			r, err = iter.NextRune()
			if err != nil || (!allowLineBreaks && r == '\n') {
				return FailedResult
			}
			n++

			if r == quoteRune && !inEscapeSeq {
				return Result{
					NumConsumed: n,
					ComputedTokens: []ComputedToken{
						{Length: n},
					},
					NextState: state,
				}
			}

			if r == '\\' && !inEscapeSeq {
				inEscapeSeq = true
				continue
			}

			if inEscapeSeq {
				inEscapeSeq = false
			}
		}
	}
}

// parseCStyleString parses a string with characters escaped by a backslash.
func parseCStyleString(quoteRune rune, allowLineBreaks bool) Func {
	return consumeCStyleString(quoteRune, allowLineBreaks).
		Map(recognizeToken(TokenRoleString))
}
