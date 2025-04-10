package constraints

import "cmp"

// Signed 有符号整型
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Unsigned 无符号整型
type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Integer 整型
type Integer interface {
	Signed | Unsigned
}

// Float 浮点型
type Float interface {
	~float32 | ~float64
}

// Number 数字类型
type Number interface {
	Integer | Float
}

// Ordered 可排序类型
type Ordered = cmp.Ordered
