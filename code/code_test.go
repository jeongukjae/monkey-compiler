package code

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMake(t *testing.T) {
	tests := []struct {
		op       Opcode
		operadns []int
		expeced  []byte
	}{
		{OpConstant, []int{65534}, []byte{byte(OpConstant), 255, 254}},
	}

	for _, tt := range tests {
		instruction, err := Make(tt.op, tt.operadns...)
		require.Nil(t, err)

		require.Equal(t, len(tt.expeced), len(instruction))
		require.ElementsMatch(t, tt.expeced, instruction)
	}
}
