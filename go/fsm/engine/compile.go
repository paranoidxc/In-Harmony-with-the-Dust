package engine

type CmdId uint64

type CmdExpr struct {
	CmdId CmdId
	Expr  Expr
}

// Compile transforms a set of expressions (one for each input command) to a state machine.
func Compile(cmdExprs []CmdExpr) (*StateMachine, error) {
	return nil, nil
}
