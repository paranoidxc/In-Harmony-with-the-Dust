package syntax

import (
	"strings"
	"unicode/utf8"
)

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

// func recognizeKeywordOrConsume(keywords []string) MapWithInputFn {
// 	// Calculate the length of the longest keyword to limit how much
// 	// of the input needs to be reprocessed.
// 	maxLength := maxStrLen(keywords)
// 	return func(result Result, iter TrackingRuneIter, state parser.State) parser.Result {
// 		if result.NumConsumed > maxLength {
// 			return result
// 		}
//
// 		s := readInputString(iter, result.NumConsumed)
// 		for _, kw := range keywords {
// 			if kw == s {
// 				token := parser.ComputedToken{
// 					Role:   parser.TokenRoleKeyword,
// 					Length: result.NumConsumed,
// 				}
// 				return parser.Result{
// 					NumConsumed:    result.NumConsumed,
// 					ComputedTokens: []parser.ComputedToken{token},
// 					NextState:      state,
// 				}
// 			}
// 		}
//
// 		return result
// 	}
// }
