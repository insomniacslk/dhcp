package bsdp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseOptGeneric(t *testing.T) {
	// Empty bytestream produces error
	_, err := ParseOptGeneric([]byte{})
	require.Error(t, err, "error from empty bytestream")

	// Good parse
	o, err := ParseOptGeneric([]byte{1, 1, 1})
	require.NoError(t, err)
	require.Equal(t, OptionMessageType, o.Code())
	require.Equal(t, MessageTypeList, MessageType(o.Data[0]))

	// Bad parse
	o, err = ParseOptGeneric([]byte{1, 2, 1})
	require.Error(t, err, "invalid length")
}

func TestOptGenericCode(t *testing.T) {
	o := OptGeneric{
		OptionCode: OptionMessageType,
		Data:       []byte{byte(MessageTypeList)},
	}
	require.Equal(t, OptionMessageType, o.Code())
}

func TestOptGenericData(t *testing.T) {
	o := OptGeneric{
		OptionCode: OptionServerIdentifier,
		Data:       []byte{192, 168, 0, 1},
	}
	require.Equal(t, []byte{192, 168, 0, 1}, o.Data)
}

func TestOptGenericToBytes(t *testing.T) {
	o := OptGeneric{
		OptionCode: OptionServerIdentifier,
		Data:       []byte{192, 168, 0, 1},
	}
	serialized := o.ToBytes()
	expected := []byte{3, 4, 192, 168, 0, 1}
	require.Equal(t, expected, serialized)
}

func TestOptGenericString(t *testing.T) {
	o := OptGeneric{
		OptionCode: OptionServerIdentifier,
		Data:       []byte{192, 168, 0, 1},
	}
	require.Equal(t, "BSDP Server Identifier -> [192 168 0 1]", o.String())
}

func TestOptGenericStringUnknown(t *testing.T) {
	o := OptGeneric{
		OptionCode: 102, // Returend option code.
		Data:       []byte{5},
	}
	require.Equal(t, "Unknown -> [5]", o.String())
}

func TestOptGenericLength(t *testing.T) {
	filename := "some_machine_name"
	o := OptGeneric{
		OptionCode: OptionMachineName,
		Data:       []byte(filename),
	}
	require.Equal(t, len(filename), o.Length())
}
