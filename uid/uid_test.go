package uid

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUUID(t *testing.T) {
	uid := UUID()
	assert.Len(t, uid, 36, "UUID length should be 36 characters")
}

func TestUUIDHex(t *testing.T) {
	uidHex := UUIDHex()
	assert.Len(t, uidHex, 32, "UUIDHex length should be 32 characters")
}

func TestUUIDBase32(t *testing.T) {
	uidBase32 := UUIDBase32()
	assert.Len(t, uidBase32, 26, "UUIDBase32 length should be 26 characters")
}

func TestUUIDBase58(t *testing.T) {
	uidBase58 := UUIDBase58()
	assert.NotEmpty(t, uidBase58, "UUIDBase58 should not be empty")
}

func TestXID(t *testing.T) {
	xid := XID()
	assert.Len(t, xid, 20, "XID length should be 20 characters")
}

func TestNumericUID(t *testing.T) {
	id := uint64(123456)
	uid := NumericUID(id)
	assert.True(t, uid >= 100000000 && uid <= 999999999, "NumericUID should be within the valid range")
}

func TestNumericUIDNano(t *testing.T) {
	uid := NumericUIDNano()
	assert.True(t, uid >= 100000000 && uid <= 999999999, "NumericUIDNano should be within the valid range")
}

func TestNumericUIDGenerator(t *testing.T) {
	// 测试正常范围
	numericUIDGen, err := NewNumericUIDGenerator(1000, 1000, 1379472, 5426620)
	assert.NoError(t, err)
	id := uint64(123456)
	uid := numericUIDGen.Generate(id)
	t.Log(uid)
	// 验证生成的UID在有效范围内，且min和max相同时，UID应一致
	assert.Equal(t, uid, numericUIDGen.minUID, "Generated UID should equal minUID when minUID == maxUID")

	// 测试其他正常范围
	numericUIDGen, err = NewNumericUIDGenerator(1000, 10000, 1379472, 5426620)
	assert.NoError(t, err)
	uid = numericUIDGen.Generate(id)
	t.Log(uid)
	assert.True(t, uid >= numericUIDGen.minUID && uid <= numericUIDGen.maxUID, "Generated UID should be within the valid range")
}

func TestNumericUIDGenerator_InvalidRange(t *testing.T) {
	// 测试生成器初始化时的无效范围
	_, err := NewNumericUIDGenerator(1000, 99, 123, 456)
	assert.Error(t, err, "Expected error when minUID > maxUID")
}
