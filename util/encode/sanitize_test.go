package encode

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSanitizeCustomEncoder(t *testing.T) {
	encoder := NewSanitizeEncoder("Lxv9frHdiqRPuzoG2MV7BmKCyFQN1ZjwApt3Eg4WDUkJY5behs8SanXT6c")
	data := map[string]interface{}{"score": 99, "pass": true}

	encoded, err := encoder.Encode(data)
	assert.NoError(t, err)

	var decoded map[string]interface{}
	err = encoder.Decode(encoded, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, data["score"], int(decoded["score"].(int8))) // msgpack returns int8
	assert.Equal(t, data["pass"], decoded["pass"])
}
