package iana

// EntID represents the Enterprise IDs as set by IANA
type EntID int

// See https://www.iana.org/assignments/enterprise-numbers/enterprise-numbers for values
const (
	EntIDCiscoSystems     EntID = 9
	EntIDCienaCorporation EntID = 1271
)

var entIDToStringMap = map[EntID]string{
	EntIDCiscoSystems:     "Cisco Systems",
	EntIDCienaCorporation: "Ciena Corporation",
}

// String returns the vendor name for a given Enterprise ID
func (e EntID) String() string {
	if vendor := entIDToStringMap[e]; vendor != "" {
		return vendor
	}
	return "Unknown"
}
