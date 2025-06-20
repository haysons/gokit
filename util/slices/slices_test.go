package slices

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUniq(t *testing.T) {
	input := []int{1, 2, 2, 3, 1, 4}
	expected := []int{1, 2, 3, 4}

	result := Uniq(input)

	assert.ElementsMatch(t, expected, result, "Uniq result does not match expected unique elements")
}

func TestChunk(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	size := 2

	result := Chunk(input, size)

	expected := [][]int{
		{1, 2},
		{3, 4},
		{5},
	}

	assert.Equal(t, expected, result, "Chunk result is not as expected")
}

func TestChunkPanic(t *testing.T) {
	assert.Panics(t, func() {
		Chunk([]int{1, 2}, 0)
	}, "Expected panic when chunk size is less than or equal to 0")
}

func TestShuffle(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	original := make([]int, len(input))
	copy(original, input)

	result := Shuffle(input)
	assert.ElementsMatch(t, original, result, "Shuffled result should contain the same elements")
	assert.Len(t, result, len(original), "Shuffled result should have the same length")
}

func TestReverse(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	expected := []int{5, 4, 3, 2, 1}

	Reverse(input)

	assert.Equal(t, expected, input, "Reversed slice does not match expected result")
}

func TestSet(t *testing.T) {
	input := []string{"a", "b", "a", "c"}
	result := Set(input)

	expected := map[string]struct{}{
		"a": {},
		"b": {},
		"c": {},
	}

	assert.Equal(t, expected, result, "Set result is incorrect")
}
