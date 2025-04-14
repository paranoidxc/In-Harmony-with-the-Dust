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
		{
			name:     "empty",
			cmdExprs: []CmdExpr{},
			expected: &StateMachine{
				numStates:   1,
				acceptCmd:   map[stateId]CmdId{},
				transitions: map[stateId][]transition{},
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
