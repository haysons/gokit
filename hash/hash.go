package hash

import (
	"github.com/cespare/xxhash/v2"
	"hash/crc32"
)

// RangeHash 对data取哈希值，哈希值最终位于闭区间[min, max]之间，非密码学安全，主要用于哈希分片相关
func RangeHash(data []byte, min, max uint64) uint64 {
	if min > max {
		panic("min should be less than max")
	}
	hash := xxhash.Sum64(data)
	rangeSize := max - min + 1
	hashInRange := min + (hash % rangeSize)
	return hashInRange
}

// Checksum 计算数据的校验和，用于判断数据完整性
func Checksum(data []byte) uint32 {
	return crc32.ChecksumIEEE(data)
}
