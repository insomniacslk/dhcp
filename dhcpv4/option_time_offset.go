package dhcpv4

import (
	"time"

	"github.com/u-root/uio/uio"
)

// TimeOffset is the signed UTC offset described by RFC 2132, Section 3.4.
// A positive value is east of UTC and a negative value is west of it.
// The wire format is a two's complement 32-bit number of seconds.
type TimeOffset time.Duration

// FromBytes parses a time offset from a byte stream according to RFC 2132, Section 3.4.
func (t *TimeOffset) FromBytes(data []byte) error {
	buf := uio.NewBigEndianBuffer(data)
	secs := int32(buf.Read32())
	*t = TimeOffset(time.Duration(secs) * time.Second)
	return buf.FinError()
}

// ToBytes returns a serialized stream of bytes for this option.
func (t TimeOffset) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)
	secs := int32(time.Duration(t) / time.Second)
	buf.Write32(uint32(secs))
	return buf.Data()
}

// String returns a human-readable string for this option.
func (t TimeOffset) String() string {
	return time.Duration(t).String()
}

// OptTimeOffset returns a new Time Offset option.
//
// The time offset option is described by RFC 2132, Section 3.4.
func OptTimeOffset(d time.Duration) Option {
	return Option{Code: OptionTimeOffset, Value: TimeOffset(d)}
}
