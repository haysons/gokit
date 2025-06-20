package maps

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeys(t *testing.T) {
	m1 := map[string]int{"a": 1, "b": 2}
	m2 := map[string]int{"c": 3}
	keys := Keys(m1, m2)
	assert.ElementsMatch(t, []string{"a", "b", "c"}, keys)
}

func TestUniqKeys(t *testing.T) {
	m1 := map[string]int{"a": 1, "b": 2}
	m2 := map[string]int{"b": 3, "c": 4}
	keys := UniqKeys(m1, m2)
	assert.ElementsMatch(t, []string{"a", "b", "c"}, keys)
}

func TestValues(t *testing.T) {
	m1 := map[string]int{"a": 1, "b": 2}
	m2 := map[string]int{"c": 3}
	values := Values(m1, m2)
	assert.ElementsMatch(t, []int{1, 2, 3}, values)
}

func TestUniqValues(t *testing.T) {
	m1 := map[string]int{"a": 1, "b": 2}
	m2 := map[string]int{"c": 2, "d": 3}
	values := UniqValues(m1, m2)
	assert.ElementsMatch(t, []int{1, 2, 3}, values)
}

func TestPairs(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	k, v := Pairs(m)
	assert.ElementsMatch(t, []string{"a", "b"}, k)
	assert.ElementsMatch(t, []int{1, 2}, v)
}

func TestInvert(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	inverted := Invert(m)
	assert.Equal(t, map[int]string{1: "a", 2: "b"}, inverted)
}

func TestMerge(t *testing.T) {
	m1 := map[string]int{"a": 1, "b": 2}
	m2 := map[string]int{"b": 3, "c": 4}
	merged := Merge(m1, m2)
	assert.Equal(t, map[string]int{"a": 1, "b": 3, "c": 4}, merged)
}

func TestChunk(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
	chunks := Chunk(m, 2)
	assert.Len(t, chunks, 2)

	total := 0
	for _, chunk := range chunks {
		total += len(chunk)
	}
	assert.Equal(t, 4, total)
}

func TestChunk_Empty(t *testing.T) {
	var m map[string]int
	chunks := Chunk(m, 2)
	assert.Empty(t, chunks)
}

func TestChunk_Panic(t *testing.T) {
	assert.Panics(t, func() {
		Chunk(map[string]int{"a": 1}, 0)
	})
}

func TestEqual(t *testing.T) {
	m1 := map[string]int{"a": 1, "b": 2}
	m2 := map[string]int{"b": 2, "a": 1}
	m3 := map[string]int{"a": 1, "b": 3}
	assert.True(t, Equal(m1, m2))
	assert.False(t, Equal(m1, m3))
}
