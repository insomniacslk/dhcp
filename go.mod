module github.com/vitrifi/dhcp

go 1.20

require (
	github.com/google/go-cmp v0.5.9
	github.com/hugelgupf/socketpair v0.0.0-20190730060125-05d35a94e714
	github.com/insomniacslk/dhcp v0.0.0-00010101000000-000000000000
	github.com/jsimonetti/rtnetlink v1.3.5
	github.com/mdlayher/netlink v1.7.2
	github.com/mdlayher/packet v1.1.2
	github.com/stretchr/testify v1.6.1
	github.com/u-root/uio v0.0.0-20230220225925-ffce2a382923
	golang.org/x/net v0.23.0
	golang.org/x/sys v0.18.0
)

require (
	github.com/davecgh/go-spew v1.1.0 // indirect
	github.com/josharian/native v1.1.0 // indirect
	github.com/mdlayher/socket v0.4.1 // indirect
	github.com/pierrec/lz4/v4 v4.1.14 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.1.0 // indirect
	golang.org/x/sync v0.3.0 // indirect
	gopkg.in/yaml.v3 v3.0.0 // indirect
)

replace github.com/insomniacslk/dhcp => github.com/vitrifi/dhcp v0.0.0-20240829085014-a3a4c1f04475
