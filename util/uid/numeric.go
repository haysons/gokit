package uid

import (
	"encoding/binary"
	"fmt"

	"github.com/dchest/siphash"
)

// NumericUIDGenerator 整形uid生成器
type NumericUIDGenerator struct {
	minUID    uint64
	maxUID    uint64
	rangeSize uint64
	key1      uint64
	key2      uint64
}

func NewNumericUIDGenerator(min, max, key1, key2 uint64) (*NumericUIDGenerator, error) {
	if min > max {
		return nil, fmt.Errorf("invalid range [%d, %d]", min, max)
	}
	return &NumericUIDGenerator{
		minUID:    min,
		maxUID:    max,
		rangeSize: max - min + 1,
		key1:      key1,
		key2:      key2,
	}, nil
}

// Generate 基于唯一id生成数字uid，数字uid范围位于闭区间[minUID, maxUID]之间
func (g *NumericUIDGenerator) Generate(ID uint64) uint64 {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], ID)
	hashed := siphash.Hash(g.key1, g.key2, buf[:])
	return (hashed % g.rangeSize) + g.minUID
}
