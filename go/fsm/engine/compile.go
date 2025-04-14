package engine

import "fmt"

type CmdId uint64

type CmdExpr struct {
	CmdId CmdId
	Expr  Expr
}

// Compile transforms a set of expressions (one for each input command) to a state machine.
func Compile(cmdExprs []CmdExpr) (*StateMachine, error) {
	if err := validateCmdExprs(cmdExprs); err != nil {
		return nil, err
	}
	return nil, nil
}

func validateCmdExprs(cmdExprs []CmdExpr) error {
	if err := validateCmdIdsUnique(cmdExprs); err != nil {
		return err
	}

	for _, cmdExpr := range cmdExprs {
		err := validateExpr(cmdExpr.Expr, false)
		if err != nil {
			return fmt.Errorf("Invalid expression for cmd %d: %w", cmdExpr.CmdId, err)
		}
	}

	return nil
}

func validateCmdIdsUnique(cmdExprs []CmdExpr) error {
	cmdIds := make(map[CmdId]struct{}, len(cmdExprs))
	for _, cmdExpr := range cmdExprs {
		_, exists := cmdIds[cmdExpr.CmdId]
		if exists {
			return fmt.Errorf("Duplicate command ID detected: %d", cmdExpr.CmdId)
		}
		cmdIds[cmdExpr.CmdId] = struct{}{}
	}
	return nil
}

func validateExpr(expr Expr, inCapture bool) error {
	switch expr := expr.(type) {
	case EventExpr:
		break
	default:
		return fmt.Errorf("Invalid expression type %T", expr)
	}

	return nil
}
