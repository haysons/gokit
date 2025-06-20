package slices

import "math/rand/v2"

// Uniq 对切片元素进行去重
func Uniq[T comparable, Slice ~[]T](collection Slice) Slice {
	result := make(Slice, 0, len(collection))
	seen := make(map[T]struct{}, len(collection))

	for i := range collection {
		if _, ok := seen[collection[i]]; ok {
			continue
		}

		seen[collection[i]] = struct{}{}
		result = append(result, collection[i])
	}

	return result
}

// Chunk 对切片进行分块，分块的长度为size
func Chunk[T any, Slice ~[]T](collection Slice, size int) []Slice {
	if size <= 0 {
		panic("Second parameter must be greater than 0")
	}

	chunksNum := len(collection) / size
	if len(collection)%size != 0 {
		chunksNum += 1
	}

	result := make([]Slice, 0, chunksNum)

	for i := 0; i < chunksNum; i++ {
		last := (i + 1) * size
		if last > len(collection) {
			last = len(collection)
		}
		result = append(result, collection[i*size:last:last])
	}

	return result
}

// Shuffle 打乱切片元素顺序
func Shuffle[T any, Slice ~[]T](collection Slice) Slice {
	rand.Shuffle(len(collection), func(i, j int) {
		collection[i], collection[j] = collection[j], collection[i]
	})
	return collection
}

// Reverse 对切片进行逆序
func Reverse[T any, Slice ~[]T](collection Slice) {
	length := len(collection)
	half := length / 2

	for i := 0; i < half; i = i + 1 {
		j := length - 1 - i
		collection[i], collection[j] = collection[j], collection[i]
	}
}

// Set 使用切片元素构造集合
func Set[T comparable, Slice ~[]T](collection Slice) map[T]struct{} {
	result := make(map[T]struct{}, len(collection))

	for _, item := range collection {
		result[item] = struct{}{}
	}

	return result
}
