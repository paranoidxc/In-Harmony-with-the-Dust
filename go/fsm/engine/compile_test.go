package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompile(t *testing.T) {
	testCases := []struct {
		name     string
		cmdExprs []CmdExpr
		expected *StateMachine
	}{
		// {
		// 	name:     "empty",
		// 	cmdExprs: []CmdExpr{},
		// 	expected: &StateMachine{
		// 		numStates:   1,
		// 		acceptCmd:   map[stateId]CmdId{},
		// 		transitions: map[stateId][]transition{},
		// 	},
		// },
		// -------------------------
		// {
		// 	name: "EventExpr",
		// 	cmdExprs: []CmdExpr{
		// 		{
		// 			CmdId: 0,
		// 			Expr:  EventExpr{Event: 99},
		// 		},
		// 	},
		// 	expected: &StateMachine{
		// 		numStates: 2,
		// 		acceptCmd: map[stateId]CmdId{
		// 			1: 0,
		// 		},
		// 		transitions: map[stateId][]transition{
		// 			0: {
		// 				{
		// 					eventRange: eventRange{start: 99, end: 99},
		// 					nextState:  1,
		// 				},
		// 			},
		// 		},
		// 	},
		// },
		// {
		// 	name: "EventRangeExpr",
		// 	cmdExprs: []CmdExpr{
		// 		{
		// 			CmdId: 0,
		// 			Expr: EventRangeExpr{
		// 				StartEvent: 23,
		// 				EndEvent:   79,
		// 			},
		// 		},
		// 	},
		// 	expected: &StateMachine{
		// 		numStates: 2,
		// 		acceptCmd: map[stateId]CmdId{
		// 			1: 0,
		// 		},
		// 		transitions: map[stateId][]transition{
		// 			0: {
		// 				{
		// 					eventRange: eventRange{start: 23, end: 79},
		// 					nextState:  1,
		// 				},
		// 			},
		// 		},
		// 	},
		// },
		// -------------
		// {
		// 	name: "ConcatExpr",
		// 	cmdExprs: []CmdExpr{
		// 		{
		// 			CmdId: 0,
		// 			Expr: ConcatExpr{
		// 				Children: []Expr{
		// 					EventExpr{Event: 12},
		// 					EventExpr{Event: 34},
		// 				},
		// 			},
		// 		},
		// 	},
		// 	expected: &StateMachine{
		// 		numStates: 3,
		// 		acceptCmd: map[stateId]CmdId{
		// 			1: 0,
		// 		},
		// 		transitions: map[stateId][]transition{
		// 			0: {
		// 				{
		// 					eventRange: eventRange{start: 12, end: 12},
		// 					nextState:  2,
		// 				},
		// 			},
		// 			2: {
		// 				{
		// 					eventRange: eventRange{start: 34, end: 34},
		// 					nextState:  1,
		// 				},
		// 			},
		// 		},
		// 	},
		// },
		// -------------------
		// {
		// 	name: "AltExpr",
		// 	cmdExprs: []CmdExpr{
		// 		{
		// 			CmdId: 0,
		// 			Expr: AltExpr{
		// 				Children: []Expr{
		// 					EventExpr{Event: 12},
		// 					EventExpr{Event: 34},
		// 				},
		// 			},
		// 		},
		// 	},
		// 	expected: &StateMachine{
		// 		numStates: 2,
		// 		acceptCmd: map[stateId]CmdId{
		// 			1: 0,
		// 		},
		// 		transitions: map[stateId][]transition{
		// 			0: {
		// 				{
		// 					eventRange: eventRange{start: 12, end: 12},
		// 					nextState:  1,
		// 				},
		// 				{
		// 					eventRange: eventRange{start: 34, end: 34},
		// 					nextState:  1,
		// 				},
		// 			},
		// 		},
		// 	},
		// },

		{
			name: "OptionExpr",
			cmdExprs: []CmdExpr{
				{
					CmdId: 0,
					Expr: OptionExpr{
						Child: EventExpr{Event: 99},
					},
				},
			},
			expected: &StateMachine{
				numStates: 2,
				acceptCmd: map[stateId]CmdId{
					0: 0,
					1: 0,
				},
				transitions: map[stateId][]transition{
					0: {
						{
							eventRange: eventRange{start: 99, end: 99},
							nextState:  1,
						},
					},
				},
			},
		},
		{
			name: "StarExpr",
			cmdExprs: []CmdExpr{
				{
					CmdId: 0,
					Expr: StarExpr{
						Child: EventExpr{Event: 99},
					},
				},
			},
			expected: &StateMachine{
				numStates: 1,
				acceptCmd: map[stateId]CmdId{
					0: 0,
				},
				transitions: map[stateId][]transition{
					0: {
						{
							eventRange: eventRange{start: 99, end: 99},
							nextState:  0,
						},
					},
				},
			},
		},
		{
			name: "CaptureExpr",
			cmdExprs: []CmdExpr{
				{
					CmdId: 0,
					Expr: CaptureExpr{
						Child: EventExpr{Event: 99},
					},
				},
			},
			expected: &StateMachine{
				numStates: 2,
				acceptCmd: map[stateId]CmdId{
					1: 0,
				},
				transitions: map[stateId][]transition{
					0: {
						{
							eventRange: eventRange{start: 99, end: 99},
							nextState:  1,
							captures: map[CmdId]CaptureId{
								0: 0,
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sm, err := Compile(tc.cmdExprs)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, sm)
		})
	}
}
