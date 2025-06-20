package encode

import "github.com/mr-tron/base58"

// SanitizeEncoder 使用msgpack对任意类型进行序列化，之后使用自定义码集的base58进行编码.
// 如果数据不想直接展示出来，可使用这种方式 进行编解码，从而起到脱敏的效果，且编码出的内容比较短，
// 注意这只起到脱敏的效果，并不是安全地加解密操作
type SanitizeEncoder struct {
	alphabet *base58.Alphabet // 编码码集
}

// NewSanitizeEncoder 通过自定义码集创建一个脱敏编码对象
func NewSanitizeEncoder(alphabet string) *SanitizeEncoder {
	return &SanitizeEncoder{
		alphabet: base58.NewAlphabet(alphabet),
	}
}

func (enc *SanitizeEncoder) Encode(v any) (string, error) {
	bytes, err := MsgpackMarshal(v)
	if err != nil {
		return "", err
	}
	return base58.FastBase58EncodingAlphabet(bytes, enc.alphabet), nil
}

func (enc *SanitizeEncoder) Decode(encoded string, v any) error {
	bytes, err := base58.FastBase58DecodingAlphabet(encoded, enc.alphabet)
	if err != nil {
		return err
	}
	return MsgpackUnmarshal(bytes, v)
}
