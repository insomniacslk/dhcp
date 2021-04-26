package iana

// EntID represents the Enterprise IDs as set by IANA
// https://www.iana.org/assignments/enterprise-numbers/enterprise-numbers
type EntID uint32

// EntIDCiscoSystems is IANA Enterprise ID for Cisco Systems
const EntIDCiscoSystems EntID = 0x9

var entIDToStringMap = map[EntID]string{
	EntIDCiscoSystems: "Cisco Systems",
}

// String returns the vendor name for a given Enterprise ID
func (e EntID) String() string {
	if vendor := entIDToStringMap[e]; vendor != "" {
		return vendor
	}
	return "Unknown"
}
