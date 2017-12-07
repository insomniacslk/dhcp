package dhcpv6

import (
	"fmt"
	"strings"
)

func LabelsFromBytes(buf []byte) ([]string, error) {
	var (
		pos     = 0
		domains = make([]string, 0)
		label   = ""
	)
	for {
		if pos >= len(buf) {
			return domains, nil
		}
		length := int(buf[pos])
		pos++
		if length == 0 {
			domains = append(domains, label)
			label = ""
		}
		if len(buf)-pos < length {
			return nil, fmt.Errorf("DomainNamesFromBytes: invalid short label length")
		}
		if label != "" {
			label += "."
		}
		label += string(buf[pos : pos+length])
		pos += length
	}
	return domains, nil
}

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

func LabelsToBytes(labels []string) []byte {
	var encodedLabels []byte
	for _, label := range labels {
		encodedLabels = append(encodedLabels, LabelToBytes(label)...)
	}
	return encodedLabels
}
