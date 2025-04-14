package uid

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"github.com/haysons/gokit/encode"
	"os"
	"sync"
)

var machineIDGenerator = NewMachineIDGenerator()

// MachineID 获取机器id，length为id的字节数
func MachineID(length int) []byte {
	mid, err := machineIDGenerator.Generate(length)
	if err != nil {
		panic(fmt.Sprintf("machineIDGenerator.Generate(%d): %v", length, err))
	}
	return mid
}

// MachineIDStr 获取机器id，length为id的字节数，使用16进制编码
func MachineIDStr(length int) string {
	mid, err := machineIDGenerator.GenerateStr(length)
	if err != nil {
		panic(fmt.Sprintf("machineIDGenerator.GenerateStr(%d): %v", length, err))
	}
	return mid
}

// MachineIDInt 获取机器id，length为id的字节数，以整数的方式返回机器id
func MachineIDInt(length int) uint64 {
	mid, err := machineIDGenerator.GenerateInt(length)
	if err != nil {
		panic(fmt.Sprintf("machineIDGenerator.GenerateInt(%d): %v", length, err))
	}
	return mid
}

// MachineIDGenerator 生成机器id，优先使用特定平台的机器唯一标识，次优先使用hostname，均不存在则使用随机值，机器id最大长度为8字节
// 目前仅在内存中缓存了生成的机器id，每次生成机器id时有可能发生变化
type MachineIDGenerator struct {
	idStore sync.Map // 保存各长度机器id值
}

func NewMachineIDGenerator() *MachineIDGenerator {
	return &MachineIDGenerator{
		idStore: sync.Map{},
	}
}

// Generate 生成机器id，length为id的字节数
func (m *MachineIDGenerator) Generate(length int) ([]byte, error) {
	if length > 8 {
		length = 8
	}
	// 缓存机器id，至少保证运行期间机器id不变
	value, ok := m.idStore.Load(length)
	if ok {
		return value.([]byte), nil
	}

	id := make([]byte, length)
	mid, err := readPlatformMachineID()
	if err != nil || len(mid) == 0 {
		mid, err = os.Hostname()
	}
	if err == nil && len(mid) != 0 {
		hw := sha256.New()
		hw.Write([]byte(mid))
		copy(id, hw.Sum(nil))
	} else {
		if _, randErr := rand.Reader.Read(id); randErr != nil {
			return nil, fmt.Errorf("machine id: cannot get hostname nor generate a random number: %v; %v", err, randErr)
		}
	}

	m.idStore.Store(length, id)
	return id, nil
}

// GenerateStr 生成机器id，使用16进制编码
func (m *MachineIDGenerator) GenerateStr(length int) (string, error) {
	machineID, err := m.Generate(length)
	if err != nil {
		return "", err
	}
	return encode.HexEncode(machineID), nil
}

// GenerateInt 生成机器id，使用小端序转化为uint64，可直接作为机器编号
func (m *MachineIDGenerator) GenerateInt(length int) (uint64, error) {
	machineID, err := m.Generate(length)
	if err != nil {
		return 0, err
	}
	var buf [8]byte
	copy(buf[:], machineID)
	return binary.LittleEndian.Uint64(buf[:]), nil
}
