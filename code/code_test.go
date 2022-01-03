package code

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMake(t *testing.T) {
	tests := []struct {
		op       Opcode
		operands []int
		expeced  []byte
	}{
		{OpConstant, []int{65534}, []byte{byte(OpConstant), 255, 254}},
		{OpAdd, []int{}, []byte{byte(OpAdd)}},
	}

	for _, tt := range tests {
		instruction := Make(tt.op, tt.operands...)

		require.Equal(t, len(tt.expeced), len(instruction))
		require.ElementsMatch(t, tt.expeced, instruction)
	}
}

func TestInstructionString(t *testing.T) {
	instructions := []Instructions{
		Make(OpAdd),
		Make(OpConstant, 1),
		Make(OpConstant, 2),
		Make(OpConstant, 65535),
	}

	expected := `0000 OpAdd
0001 OpConstant 1
0004 OpConstant 2
0007 OpConstant 65535
`

	concatenated := Instructions{}
	for _, ins := range instructions {
		concatenated = append(concatenated, ins...)
	}

	require.Equal(t, expected, concatenated.String(), "instructions wrongly formatted")
}

func TestReadOperands(t *testing.T) {
	tests := []struct {
		op        Opcode
		operands  []int
		bytesRead int
	}{
		{OpConstant, []int{65535}, 2},
	}

	for _, tt := range tests {
		instruction := Make(tt.op, tt.operands...)
		def, err := Lookup(byte(tt.op))
		require.Nil(t, err, "definition not found")

		operandsRead, n := ReadOperands(def, instruction[1:])
		require.Equal(t, tt.bytesRead, n, "wrong #bytesReads")

		require.ElementsMatch(t, tt.operands, operandsRead, "operand wrong")
	}
}