package dhcpv4

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptionGenericCode(t *testing.T) {
	o := OptionGeneric{
		OptionCode: OptionDHCPMessageType,
		Data:       []byte{byte(MessageTypeDiscover)},
	}
	require.Equal(t, OptionDHCPMessageType, o.Code())
}

func TestOptionGenericData(t *testing.T) {
	o := OptionGeneric{
		OptionCode: OptionNameServer,
		Data:       []byte{192, 168, 0, 1},
	}
	require.Equal(t, []byte{192, 168, 0, 1}, o.Data)
}

func TestOptionGenericToBytes(t *testing.T) {
	o := OptionGeneric{
		OptionCode: OptionDHCPMessageType,
		Data:       []byte{byte(MessageTypeDiscover)},
	}
	serialized := o.ToBytes()
	expected := []byte{53, 1, 1}
	require.Equal(t, expected, serialized)
}

func TestOptionGenericToBytesZeroOptions(t *testing.T) {
	o := OptionGeneric{OptionCode: OptionEnd}
	serialized := o.ToBytes()
	expected := []byte{255}
	require.Equal(t, expected, serialized)

	o = OptionGeneric{OptionCode: OptionPad}
	serialized = o.ToBytes()
	expected = []byte{0}
	require.Equal(t, expected, serialized)
}

func TestOptionGenericString(t *testing.T) {
	o := OptionGeneric{
		OptionCode: OptionDHCPMessageType,
		Data:       []byte{byte(MessageTypeDiscover)},
	}
	require.Equal(t, "DHCP Message Type -> [1]", o.String())
}

func TestOptionGenericStringUnknown(t *testing.T) {
	o := OptionGeneric{
		OptionCode: 102, // Returend option code.
		Data:       []byte{byte(MessageTypeDiscover)},
	}
	require.Equal(t, "Unknown -> [1]", o.String())
}

func TestOptionGenericLength(t *testing.T) {
	filename := "/path/to/file"
	o := OptionGeneric{
		OptionCode: OptionBootfileName,
		Data:       []byte(filename),
	}
	require.Equal(t, len(filename), o.Length())
}
