package uid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
