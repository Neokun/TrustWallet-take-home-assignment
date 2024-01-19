package utils

import (
	"fmt"
	"strconv"
)

// EncodeUint64 encodes i as a hex string with 0x prefix.
func EncodeUint64(i uint64) string {
	enc := make([]byte, 2, 10)
	copy(enc, "0x")
	return string(strconv.AppendUint(enc, i, 16))
}

// HexUint64 is a custom type based on uint64 that can json unmarshal hex string to uint64.
type HexUint64 uint64

func (h *HexUint64) UnmarshalJSON(input []byte) error {
	if !isString(input) {
		return fmt.Errorf("input must be a string")
	}

	return h.UnmarshalText(input[1 : len(input)-1])
}

// UnmarshalText implements encoding.TextUnmarshaler
func (h *HexUint64) UnmarshalText(input []byte) error {
	raw, err := checkNumberText(input)
	if err != nil {
		return err
	}
	if len(raw) > 16 {
		return fmt.Errorf("not in uint64 range")
	}
	var dec uint64
	for _, byte := range raw {
		nib := decodeNibble(byte)
		if nib == badNibble {
			return fmt.Errorf("invalid hex string")
		}
		dec *= 16
		dec += nib
	}
	*h = HexUint64(dec)
	return nil
}

func checkNumberText(input []byte) (raw []byte, err error) {
	if len(input) == 0 {
		return nil, nil // empty strings are allowed
	}
	if !bytesHave0xPrefix(input) {
		return nil, fmt.Errorf("missing 0x prefix")
	}
	input = input[2:]
	if len(input) == 0 {
		return nil, fmt.Errorf("missing 0x")
	}
	if len(input) > 1 && input[0] == '0' {
		return nil, fmt.Errorf("leading zero")
	}
	return input, nil
}

func isString(input []byte) bool {
	return len(input) >= 2 && input[0] == '"' && input[len(input)-1] == '"'
}

func bytesHave0xPrefix(input []byte) bool {
	return len(input) >= 2 && input[0] == '0' && (input[1] == 'x' || input[1] == 'X')
}

const badNibble = ^uint64(0)

func decodeNibble(in byte) uint64 {
	switch {
	case in >= '0' && in <= '9':
		return uint64(in - '0')
	case in >= 'A' && in <= 'F':
		return uint64(in - 'A' + 10)
	case in >= 'a' && in <= 'f':
		return uint64(in - 'a' + 10)
	default:
		return badNibble
	}
}
