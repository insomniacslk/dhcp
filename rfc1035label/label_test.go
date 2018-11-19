package rfc1035label

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLabelsFromBytes(t *testing.T) {
	labels, err := LabelsFromBytes([]byte{
		0x9, 's', 'l', 'a', 'c', 'k', 'w', 'a', 'r', 'e',
		0x2, 'i', 't',
		0x0,
	})
	require.NoError(t, err)
	require.Equal(t, 1, len(labels))
	require.Equal(t, "slackware.it", labels[0])
}

func TestLabelsFromBytesZeroLength(t *testing.T) {
	labels, err := LabelsFromBytes([]byte{})
	require.NoError(t, err)
	require.Equal(t, 0, len(labels))
}

func TestLabelsFromBytesInvalidLength(t *testing.T) {
	labels, err := LabelsFromBytes([]byte{0x5, 0xaa, 0xbb}) // short length
	require.Error(t, err)
	require.Equal(t, 0, len(labels))
}

func TestLabelsFromBytesInvalidLengthOffByOne(t *testing.T) {
	labels, err := LabelsFromBytes([]byte{0x3, 0xaa, 0xbb}) // short length
	require.Error(t, err)
	require.Equal(t, 0, len(labels))
}

func TestLabelToBytes(t *testing.T) {
	encodedLabel := LabelToBytes("slackware.it")
	expected := []byte{
		0x9, 's', 'l', 'a', 'c', 'k', 'w', 'a', 'r', 'e',
		0x2, 'i', 't',
		0x0,
	}
	require.Equal(t, expected, encodedLabel)
}

func TestLabelsToBytes(t *testing.T) {
	expected := []byte{
		9, 's', 'l', 'a', 'c', 'k', 'w', 'a', 'r', 'e',
		2, 'i', 't',
		0,
		9, 'i', 'n', 's', 'o', 'm', 'n', 'i', 'a', 'c',
		9, 's', 'l', 'a', 'c', 'k', 'w', 'a', 'r', 'e',
		2, 'i', 't',
		0,
	}
	encodedLabels := LabelsToBytes([]string{"slackware.it", "insomniac.slackware.it"})
	require.Equal(t, expected, encodedLabels)
}

func TestLabelToBytesZeroLength(t *testing.T) {
	encodedLabel := LabelToBytes("")
	expected := []byte{0}
	require.Equal(t, expected, encodedLabel)
}

func TestCompressedLabel(t *testing.T) {
	data := []byte{
		// slackware.it
		9, 's', 'l', 'a', 'c', 'k', 'w', 'a', 'r', 'e',
		2, 'i', 't',
		0,
		// insomniac.slackware.it
		9, 'i', 'n', 's', 'o', 'm', 'n', 'i', 'a', 'c',
		192, 0,
		// mail.systemboot.org
		4, 'm', 'a', 'i', 'l',
		10, 's', 'y', 's', 't', 'e', 'm', 'b', 'o', 'o', 't',
		3, 'o', 'r', 'g',
		0,
		// systemboot.org
		192, 31,
	}
	expected := []string{
		"slackware.it",
		"insomniac.slackware.it",
		"mail.systemboot.org",
		"systemboot.org",
	}

	labels, err := LabelsFromBytes(data)
	require.NoError(t, err)
	require.Equal(t, expected, labels)
}

func TestShortCompressedLabel(t *testing.T) {
	data := []byte{
		// slackware.it
		9, 's', 'l', 'a', 'c', 'k', 'w', 'a', 'r', 'e',
		2, 'i', 't',
		0,
		// insomniac.slackware.it
		9, 'i', 'n', 's', 'o', 'm', 'n', 'i', 'a', 'c',
		192,
	}

	_, err := LabelsFromBytes(data)
	require.Error(t, err)
}

func TestNestedCompressedLabel(t *testing.T) {
	data := []byte{
		// it
		3, 'i', 't',
		0,
		// slackware.it
		9, 's', 'l', 'a', 'c', 'k', 'w', 'a', 'r', 'e',
		192, 0,
		// insomniac.slackware.it
		9, 'i', 'n', 's', 'o', 'm', 'n', 'i', 'a', 'c',
		192, 5,
	}
	_, err := LabelsFromBytes(data)
	require.Error(t, err)
}
