package syntax

//import "log/slog"

// MapFn maps a successful parse result to another parse result.
type MapFn func(Result) Result

// Map maps a successful parse result to another parse result using mapFn.
func (f Func) Map(mapFn MapFn) Func {
	return func(iter TrackingRuneIter, state State) Result {
		result := f(iter, state)
		if result.IsFailure() {
			return FailedResult
		}
		return mapFn(result)
	}
}

// MapWithInputFn maps a successful parse result to another parse result,
// using the original input (iter + state) that produced the first parse result.
type MapWithInputFn func(Result, TrackingRuneIter, State) Result

// MapWithInput maps a successful parse to another parse result according to mapFn.
// The input iterator will output only runes consumed by the result before returning EOF.
func (f Func) MapWithInput(mapFn MapWithInputFn) Func {
	return func(iter TrackingRuneIter, state State) Result {
		result := f(iter, state)
		if result.IsFailure() {
			return FailedResult
		}
		iter.Limit(result.NumConsumed)
		return mapFn(result, iter, state)
	}
}

// Then produces a parse func that succeeds if both `f` and `nextFn` succeed.
func (f Func) Then(nextFn Func) Func {
	return func(iter TrackingRuneIter, state State) Result {
		result := f(iter, state)
		if result.IsFailure() {
			return FailedResult
		}

		iter.Skip(result.NumConsumed)
		nextResult := nextFn(iter, result.NextState)
		if nextResult.IsFailure() {
			return FailedResult
		}

		return combineSeqResults(result, nextResult)
	}
}

func (f Func) ThenMaybe(nextFn Func) Func {
	return func(iter TrackingRuneIter, state State) Result {
		result := f(iter, state)
		if result.IsFailure() {
			return FailedResult
		}

		//slog.Info("ThenMaybe Skip", slog.Any("num", result.NumConsumed))

		iter.Skip(result.NumConsumed)
		nextResult := nextFn(iter, result.NextState)
		if nextResult.IsFailure() {
			return result
		}

		return combineSeqResults(result, nextResult)
	}
}

func combineSeqResults(r1, r2 Result) Result {
	//return Result{}
	tokens := r1.ComputedTokens
	for _, tok := range r2.ComputedTokens {
		tokens = append(tokens, ComputedToken{
			Offset: r1.NumConsumed + tok.Offset,
			Length: tok.Length,
			Role:   tok.Role,
		})
	}

	result := Result{
		NumConsumed:    r1.NumConsumed + r2.NumConsumed,
		ComputedTokens: tokens,
		NextState:      r2.NextState,
	}

	//slog.Info("combineSeqResults", slog.Any("result", result))
	return result
}

// Or produces a parse func that returns the result of `f` if it succeeds,
// or the result of `nextFn` if `f` fails.
func (f Func) Or(nextFn Func) Func {
	return func(iter TrackingRuneIter, state State) Result {
		result := f(iter, state)
		if result.IsSuccess() {
			return result
		}
		return nextFn(iter, state)
	}
}

func (f Func) recoverFromFailure() Func {
	return func(iter TrackingRuneIter, state State) Result {
		var numSkipped uint64
		for {
			result := f(iter, state)

			// slog.Info("recoverFromFailure",
			// 	slog.Any("result", result))

			if result.IsSuccess() {
				// slog.Info("recoverFromFailure",
				// 	slog.Any("result", result))
				return result.ShiftForward(numSkipped)
			}

			// Recover by skipping one rune ahead.
			// slog.Info("recoverFromFailure Skip",
			// 	slog.Any("num", 1))
			n := iter.Skip(1)
			numSkipped += n
			if n == 0 {
				return Result{
					NumConsumed: numSkipped,
					NextState:   state,
				}
			}
		}
	}
}
