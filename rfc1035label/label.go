package rfc1035label

import (
	"errors"
	"strings"
)

// This implements RFC 1035 labels, including compression.
// https://tools.ietf.org/html/rfc1035#section-4.1.4

// LabelsFromBytes decodes a serialized stream and returns a list of labels
func LabelsFromBytes(buf []byte) ([]string, error) {
	var (
		pos, oldPos     int
		labels          = make([]string, 0)
		label           string
		handlingPointer bool
	)

	for {
		if pos >= len(buf) {
			break
		}
		length := int(buf[pos])
		pos++
		var chunk string
		if length == 0 {
			labels = append(labels, label)
			label = ""
			if handlingPointer {
				pos = oldPos
				handlingPointer = false
			}
		} else if length&0xc0 == 0xc0 {
			// compression pointer
			if handlingPointer {
				return nil, errors.New("rfc1035label: cannot handle nested pointers")
			}
			handlingPointer = true
			if pos+1 > len(buf) {
				return nil, errors.New("rfc1035label: pointer buffer too short")
			}
			off := int(buf[pos-1]&^0xc0)<<8 + int(buf[pos])
			oldPos = pos + 1
			pos = off
		} else {
			if pos+length > len(buf) {
				return nil, errors.New("rfc1035label: buffer too short")
			}
			chunk = string(buf[pos : pos+length])
			if label != "" {
				label += "."
			}
			label += chunk
			pos += length
		}
	}
	return labels, nil
}

// LabelToBytes encodes a label and returns a serialized stream of bytes
func LabelToBytes(label string) []byte {
	var encodedLabel []byte
	if len(label) == 0 {
		return []byte{0}
	}
	for _, part := range strings.Split(label, ".") {
		encodedLabel = append(encodedLabel, byte(len(part)))
		encodedLabel = append(encodedLabel, []byte(part)...)
	}
	return append(encodedLabel, 0)
}

// LabelsToBytes encodes a list of labels and returns a serialized stream of
// bytes
func LabelsToBytes(labels []string) []byte {
	var encodedLabels []byte
	for _, label := range labels {
		encodedLabels = append(encodedLabels, LabelToBytes(label)...)
	}
	return encodedLabels
}
