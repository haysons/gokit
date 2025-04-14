package uid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMachineID(t *testing.T) {
	length := 4
	id1 := MachineID(length)
	id2 := MachineID(length)

	assert.Equal(t, length, len(id1), "MachineID should return a byte slice of specified length")
	assert.Equal(t, id1, id2, "MachineID should cache the result for the same length")
}

func TestMachineIDStr(t *testing.T) {
	length := 4
	str1 := MachineIDStr(length)
	str2 := MachineIDStr(length)

	assert.Equal(t, length*2, len(str1), "Hex encoded machine ID string length should be twice the number of bytes")
	assert.Equal(t, str1, str2, "MachineIDStr should cache the result for the same length")
}

func TestMachineIDInt(t *testing.T) {
	length := 4
	id1 := MachineIDInt(length)
	id2 := MachineIDInt(length)

	assert.Equal(t, id1, id2, "MachineIDInt should return the same integer value for the same length")
}

func TestLengthGreaterThan8(t *testing.T) {
	id := MachineID(10)

	assert.Equal(t, 8, len(id), "If the requested length is greater than 8, the machine ID should be trimmed to 8 bytes")
}
