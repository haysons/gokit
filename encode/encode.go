package encode

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"github.com/mr-tron/base58"
	"github.com/vmihailenco/msgpack/v5"
)

// HexEncode 16进制编码
func HexEncode(bytes []byte) string {
	return hex.EncodeToString(bytes)
}

// HexDecode 16进制解码
func HexDecode(encoded string) ([]byte, error) {
	return hex.DecodeString(encoded)
}

// Base32Encode base32编码
func Base32Encode(bytes []byte) string {
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(bytes)
}

// Base32Decode base32解码
func Base32Decode(encoded string) ([]byte, error) {
	return base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(encoded)
}

// Base64Encode base64编码
func Base64Encode(bytes []byte) string {
	return base64.StdEncoding.EncodeToString(bytes)
}

// Base64Decode base64解码
func Base64Decode(encoded string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encoded)
}

// Base58Encode base58编码
func Base58Encode(bytes []byte) string {
	return base58.Encode(bytes)
}

// Base58Decode base58解码
func Base58Decode(encoded string) ([]byte, error) {
	return base58.Decode(encoded)
}

// Uint32Encode 使用小端序编码uint32
func Uint32Encode(v uint32) []byte {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, v)
	return buf
}

// Uint32Decode 使用小端序解码uint32
func Uint32Decode(bytes []byte) uint32 {
	return binary.LittleEndian.Uint32(bytes)
}

// Uint64Encode 使用小端序编码uint64
func Uint64Encode(v uint64) []byte {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, v)
	return buf
}

// Uint64Decode 使用小端序解码uint64
func Uint64Decode(bytes []byte) uint64 {
	return binary.LittleEndian.Uint64(bytes)
}

// MsgpackMarshal 进行msgpack序列化
func MsgpackMarshal(v any) ([]byte, error) {
	return msgpack.Marshal(v)
}

// MsgpackUnmarshal 进行msgpack反序列化
func MsgpackUnmarshal(bytes []byte, v any) error {
	return msgpack.Unmarshal(bytes, v)
}

// SanitizeEncode 针对于任意类型进行脱敏编码
func SanitizeEncode(v any) (string, error) {
	return sanitizeEncoder.Encode(v)
}

// SanitizeDecode 解码脱敏后的编码值
func SanitizeDecode(encoded string, v any) error {
	return sanitizeEncoder.Decode(encoded, v)
}

var sanitizeEncoder = NewSanitizeEncoder("Lxv9frHdiqg4WDUkJY5behs8SanXT6cRPuzoG2MV7BmKCyFQN1ZjwApt3E")

// SanitizeEncoder 使用msgpack对任意类型进行序列化，之后使用自定义码集的base58进行编码，如果数据不想直接展示出来，可使用这种方式
// 进行编解码，从而起到脱敏的效果，且编码出的内容比较短，注意这只起到脱敏的效果，并不是安全地加解密操作
type SanitizeEncoder struct {
	alphabet *base58.Alphabet // 编码码集
}

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
