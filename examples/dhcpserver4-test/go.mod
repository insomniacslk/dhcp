module dhcpserver4-test

go 1.23.0

require github.com/insomniacslk/dhcp v0.0.0

require (
	github.com/josharian/native v1.1.0 // indirect
	github.com/pierrec/lz4/v4 v4.1.14 // indirect
	github.com/u-root/uio v0.0.0-20230220225925-ffce2a382923 // indirect
	golang.org/x/sys v0.31.0 // indirect
)

replace github.com/insomniacslk/dhcp => ../../../dhcp
