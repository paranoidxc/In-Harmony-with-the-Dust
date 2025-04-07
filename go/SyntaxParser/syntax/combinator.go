package syntax

import "log/slog"

func (f Func) recoverFromFailure() Func {
	return func(iter TrackingRuneIter, state State) Result {
		var numSkipped uint64
		for {
			result := f(iter, state)
			slog.Info("recoverFromFailure",
				slog.Any("result", result))
			if result.IsSuccess() {
				return result.ShiftForward(numSkipped)
			}

			// Recover by skipping one rune ahead.
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
