package dhcpv4

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseOptionGeneric(t *testing.T) {
	// Empty bytestream produces error
	_, err := ParseOptionGeneric(OptionHostName, []byte{})
	require.Error(t, err, "error from empty bytestream")
}

func TestOptionGenericCode(t *testing.T) {
	o := OptionGeneric{
		OptionCode: OptionDHCPMessageType,
		Data:       []byte{byte(MessageTypeDiscover)},
	}
	require.Equal(t, OptionDHCPMessageType, o.Code())
}

func TestOptionGenericToBytes(t *testing.T) {
	o := OptionGeneric{
		OptionCode: OptionDHCPMessageType,
		Data:       []byte{byte(MessageTypeDiscover)},
	}
	serialized := o.ToBytes()
	expected := []byte{1}
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
