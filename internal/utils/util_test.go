package utils

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeUint64(t *testing.T) {
	number := uint64(19041293)

	hex := EncodeUint64(number)
	assert.Equal(t, "0x1228c0d", hex)
}

func TestHexUInt64_UnmarshalJSON(t *testing.T) {
	data := map[string]string{
		// 19041293
		"number": "0x1228c0d",
	}

	b, err := json.Marshal(data)
	assert.NoError(t, err)

	expect := struct {
		Number HexUint64 `json:"number"`
	}{}
	err = json.Unmarshal(b, &expect)
	assert.NoError(t, err)
	assert.Equal(t, uint64(19041293), uint64(expect.Number))

	b, err = json.Marshal(expect)
	assert.NoError(t, err)
	t.Logf("%s", b)
}
