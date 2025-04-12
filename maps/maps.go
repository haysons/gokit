package maps

import "maps"

// Keys 获取一组map的key列表
func Keys[K comparable, V any](in ...map[K]V) []K {
	size := 0
	for i := range in {
		size += len(in[i])
	}
	result := make([]K, 0, size)
	for i := range in {
		for k := range in[i] {
			result = append(result, k)
		}
	}
	return result
}

// UniqKeys 获取一组map的去重后key列表
func UniqKeys[K comparable, V any](in ...map[K]V) []K {
	size := 0
	for i := range in {
		size += len(in[i])
	}
	seen := make(map[K]struct{}, size)
	result := make([]K, 0)
	for i := range in {
		for k := range in[i] {
			if _, exists := seen[k]; exists {
				continue
			}
			seen[k] = struct{}{}
			result = append(result, k)
		}
	}
	return result
}

// Values 获取一组map的value列表
func Values[K comparable, V any](in ...map[K]V) []V {
	size := 0
	for i := range in {
		size += len(in[i])
	}
	result := make([]V, 0, size)

	for i := range in {
		for k := range in[i] {
			result = append(result, in[i][k])
		}
	}

	return result
}

// UniqValues 获取一组map的去重后value列表
func UniqValues[K comparable, V comparable](in ...map[K]V) []V {
	size := 0
	for i := range in {
		size += len(in[i])
	}

	seen := make(map[V]struct{}, size)
	result := make([]V, 0)

	for i := range in {
		for k := range in[i] {
			val := in[i][k]
			if _, exists := seen[val]; exists {
				continue
			}
			seen[val] = struct{}{}
			result = append(result, val)
		}
	}

	return result
}

// Pairs 提取map中的kv，分别返回，一些kv存储需要以此参数传参
func Pairs[K comparable, V any](in map[K]V) ([]K, []V) {
	keys := make([]K, 0, len(in))
	values := make([]V, 0, len(in))
	for k := range in {
		keys = append(keys, k)
		values = append(values, in[k])
	}
	return keys, values
}

// Invert 将map的值作为键，键作为值
func Invert[K comparable, V comparable](in map[K]V) map[V]K {
	out := make(map[V]K, len(in))

	for k := range in {
		out[in[k]] = k
	}

	return out
}

// Merge 将多个map进行合并，多个map存在同一个键时参数靠后map的值将覆盖靠前的map
func Merge[K comparable, V any, Map ~map[K]V](maps ...Map) Map {
	count := 0
	for i := range maps {
		count += len(maps[i])
	}

	out := make(Map, count)
	for i := range maps {
		for k := range maps[i] {
			out[k] = maps[i][k]
		}
	}

	return out
}

// Chunk 针对map进行分块，分块的长度为size
func Chunk[K comparable, V any](m map[K]V, size int) []map[K]V {
	if size <= 0 {
		panic("The chunk size must be greater than 0")
	}

	count := len(m)
	if count == 0 {
		return []map[K]V{}
	}

	chunksNum := count / size
	if count%size != 0 {
		chunksNum += 1
	}

	result := make([]map[K]V, 0, chunksNum)

	for k, v := range m {
		if len(result) == 0 || len(result[len(result)-1]) == size {
			result = append(result, make(map[K]V, size))
		}

		result[len(result)-1][k] = v
	}

	return result
}

// Equal 判断两个map是否相等
func Equal[M1, M2 ~map[K]V, K, V comparable](m1 M1, m2 M2) bool {
	return maps.Equal(m1, m2)
}
