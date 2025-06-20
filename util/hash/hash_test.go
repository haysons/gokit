package hash

import (
	"github.com/stretchr/testify/assert"
	"hash/crc32"
	"strconv"
	"testing"
)

func TestHashInRange(t *testing.T) {
	rangeMin, rangeMax := uint64(0), uint64(999)
	hashCounts := make(map[uint64]int)

	var totalInRange int
	for i := 0; i < 100000; i++ {
		data := []byte("test data 1234 ,./+" + strconv.Itoa(i))
		hashInRange := RangeHash(data, rangeMin, rangeMax)
		if hashInRange >= rangeMin && hashInRange <= rangeMax {
			totalInRange++
		}
		hashCounts[hashInRange]++
	}
	assert.Equal(t, totalInRange, 100000, "All hashes should be within the range")

	// 查找命中最少和最多的哈希值及其数量
	var minCount, maxCount int
	var minHash, maxHash uint64
	for hash, count := range hashCounts {
		if count > maxCount {
			maxCount = count
			maxHash = hash
		}
		if minCount == 0 || count < minCount {
			minCount = count
			minHash = hash
		}
	}
	t.Logf("Number of unique hashes: %d", len(hashCounts))
	t.Logf("Most frequent hash: %d with count: %d", maxHash, maxCount)
	t.Logf("Least frequent hash: %d with count: %d", minHash, minCount)
	assert.LessOrEqual(t, maxCount-minCount, 100, "The difference between the most and least frequent hashes should not exceed 100")
	assert.Equal(t, len(hashCounts), int(rangeMax-rangeMin+1), "The number of unique hashes should be equal with the range")
}

func TestChecksum(t *testing.T) {
	data1 := []byte("Hello, world!")
	checksum1 := Checksum(data1)
	expectedChecksum1 := crc32.ChecksumIEEE(data1)
	assert.Equal(t, checksum1, expectedChecksum1, "Checksum should be the same")

	data2 := []byte("")
	checksum2 := Checksum(data2)
	expectedChecksum2 := crc32.ChecksumIEEE(data2)
	assert.Equal(t, checksum2, expectedChecksum2, "Checksum should be the same for empty data")

	data3 := []byte("Different data!")
	checksum3 := Checksum(data3)
	checksum4 := Checksum([]byte("Different data!!"))
	assert.NotEqual(t, checksum3, checksum4, "Checksums should not be equal for different data")
}
