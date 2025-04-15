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
	bucketNum := max - min + 1
	return JumpHash(hash, bucketNum) + min
}

// JumpHash google jump hash，一种高性能的一致性哈希算法，此函数将返回key所在的桶的位置，桶的位置位于左闭右开区间[0, bucketNum)之中。
func JumpHash(key uint64, bucketNum uint64) uint64 {
	if bucketNum < 1 {
		panic("bucketNum must be >= 1")
	}

	var b, j int64

	num := int64(bucketNum)
	for j < num {
		b = j
		key = key*2862933555777941757 + 1
		j = int64(float64(b+1) * (float64(int64(1)<<31) / float64((key>>33)+1)))
	}

	return uint64(b)
}

// Checksum 计算数据的校验和，用于判断数据完整性
func Checksum(data []byte) uint32 {
	return crc32.ChecksumIEEE(data)
}
