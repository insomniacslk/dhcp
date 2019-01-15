package dhcpv4

import (
	"fmt"
)

// OptDomainName implements the domain name option described in RFC 2132,
// Section 3.17.
type OptDomainName struct {
	DomainName string
}

// ParseOptDomainName returns a new OptDomainName from a byte stream, or error
// if any.
func ParseOptDomainName(data []byte) (*OptDomainName, error) {
	return &OptDomainName{DomainName: string(data)}, nil
}

// Code returns the option code.
func (o *OptDomainName) Code() OptionCode {
	return OptionDomainName
}

// ToBytes returns a serialized stream of bytes for this option.
func (o *OptDomainName) ToBytes() []byte {
	return []byte(o.DomainName)
}

// String returns a human-readable string.
func (o *OptDomainName) String() string {
	return fmt.Sprintf("Domain Name -> %v", o.DomainName)
}

// OptHostName implements the host name option described by RFC 2132, Section
// 3.14.
type OptHostName struct {
	HostName string
}

// ParseOptHostName returns a new OptHostName from a byte stream, or error if
// any.
func ParseOptHostName(data []byte) (*OptHostName, error) {
	return &OptHostName{HostName: string(data)}, nil
}

// Code returns the option code.
func (o *OptHostName) Code() OptionCode {
	return OptionHostName
}

// ToBytes returns a serialized stream of bytes for this option.
func (o *OptHostName) ToBytes() []byte {
	return []byte(o.HostName)
}

// String returns a human-readable string.
func (o *OptHostName) String() string {
	return fmt.Sprintf("Host Name -> %v", o.HostName)
}

// OptRootPath implements the root path option described by RFC 2132, Section
// 3.19.
type OptRootPath struct {
	Path string
}

// ParseOptRootPath constructs an OptRootPath struct from a sequence of  bytes
// and returns it, or an error.
func ParseOptRootPath(data []byte) (*OptRootPath, error) {
	return &OptRootPath{Path: string(data)}, nil
}

// Code returns the option code.
func (o *OptRootPath) Code() OptionCode {
	return OptionRootPath
}

// ToBytes returns a serialized stream of bytes for this option.
func (o *OptRootPath) ToBytes() []byte {
	return []byte(o.Path)
}

// String returns a human-readable string for this option.
func (o *OptRootPath) String() string {
	return fmt.Sprintf("Root Path -> %v", o.Path)
}

// OptBootfileName implements the bootfile name option described in RFC 2132,
// Section 9.5.
type OptBootfileName struct {
	BootfileName string
}

// Code returns the option code
func (op *OptBootfileName) Code() OptionCode {
	return OptionBootfileName
}

// ToBytes serializes the option and returns it as a sequence of bytes
func (op *OptBootfileName) ToBytes() []byte {
	return []byte(op.BootfileName)
}

func (op *OptBootfileName) String() string {
	return fmt.Sprintf("Bootfile Name -> %s", op.BootfileName)
}

// ParseOptBootfileName returns a new OptBootfile from a byte stream or error if any
func ParseOptBootfileName(data []byte) (*OptBootfileName, error) {
	return &OptBootfileName{BootfileName: string(data)}, nil
}

// OptTFTPServerName implements the TFTP server name option described by RFC
// 2132, Section 9.4.
type OptTFTPServerName struct {
	TFTPServerName string
}

// Code returns the option code
func (op *OptTFTPServerName) Code() OptionCode {
	return OptionTFTPServerName
}

// ToBytes serializes the option and returns it as a sequence of bytes
func (op *OptTFTPServerName) ToBytes() []byte {
	return []byte(op.TFTPServerName)
}

func (op *OptTFTPServerName) String() string {
	return fmt.Sprintf("TFTP Server Name -> %s", op.TFTPServerName)
}

// ParseOptTFTPServerName returns a new OptTFTPServerName from a byte stream or error if any
func ParseOptTFTPServerName(data []byte) (*OptTFTPServerName, error) {
	return &OptTFTPServerName{TFTPServerName: string(data)}, nil
}

// OptClassIdentifier implements the vendor class identifier option described
// in RFC 2132, Section 9.13.
type OptClassIdentifier struct {
	Identifier string
}

// ParseOptClassIdentifier constructs an OptClassIdentifier struct from a sequence of
// bytes and returns it, or an error.
func ParseOptClassIdentifier(data []byte) (*OptClassIdentifier, error) {
	return &OptClassIdentifier{Identifier: string(data)}, nil
}

// Code returns the option code.
func (o *OptClassIdentifier) Code() OptionCode {
	return OptionClassIdentifier
}

// ToBytes returns a serialized stream of bytes for this option.
func (o *OptClassIdentifier) ToBytes() []byte {
	return []byte(o.Identifier)
}

// String returns a human-readable string for this option.
func (o *OptClassIdentifier) String() string {
	return fmt.Sprintf("Class Identifier -> %v", o.Identifier)
}
