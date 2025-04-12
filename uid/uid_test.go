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
