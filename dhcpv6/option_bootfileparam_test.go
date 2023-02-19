package dhcpv6

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	testBootfileParams0Compiled = "\x00\x0eroot=/dev/sda1\x00\x00\x00\x02rw"
	testBootfileParams1         = []string{
		"initrd=http://myserver.mycompany.local/initrd.xz",
		"",
		"root=/dev/sda1",
		"rw",
		"netconsole=..:\000:.something\000here.::..",
		string(make([]byte, (1<<16)-1)),
	}
)

// compileTestBootfileParams is an independent implementation of bootfile param encoder
func compileTestBootfileParams(t *testing.T, params []string) []byte {
	var length [2]byte
	buf := bytes.Buffer{}
	for _, param := range params {
		if len(param) >= 1<<16 {
			panic("a too long parameter")
		}
		binary.BigEndian.PutUint16(length[:], uint16(len(param)))
		_, err := buf.Write(length[:])
		require.NoError(t, err)
		_, err = buf.WriteString(param)
		require.NoError(t, err)
	}

	return buf.Bytes()
}

func TestOptBootFileParam(t *testing.T) {
	expected := string(compileTestBootfileParams(t, testBootfileParams1))
	var opt optBootFileParam
	if err := opt.FromBytes([]byte(expected)); err != nil {
		t.Fatal(err)
	}
	if string(opt.ToBytes()) != expected {
		t.Fatalf("Invalid boot file parameter. Expected %v, got %v", expected, opt)
	}
}

func TestParsedTypeOptBootFileParam(t *testing.T) {
	tryParse := func(compiled []byte, expected []string) {
		var opt optBootFileParam
		err := opt.FromBytes([]byte(compiled))
		require.NoError(t, err)
		require.Equal(t, compiled, opt.ToBytes())
		require.Equal(t, expected, opt.params)
	}

	tryParse(
		[]byte(testBootfileParams0Compiled),
		[]string{"root=/dev/sda1", "", "rw"},
	)
	tryParse(
		compileTestBootfileParams(t, testBootfileParams1),
		testBootfileParams1,
	)
}
