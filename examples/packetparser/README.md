# packetparser

An example of the encoding/decoding capabilities of the `dhcp` library, for both
DHCPv4 and DHCPv6. It can read a `pcap` file, or work as a client (which is the
default). See `./packetparser -h` for details, shown below for convenience:

```
$ ./packetparser -h
Usage of ./packetparser:
  -c int
        Stop after <count> packets (default: 0)
  -debug
        Enable debug output (default: false)
  -etherip
        Enables LayerTypeEtherIP instead of LayerTypeEthernet, use with linux-cooked PCAP files. (default: false)
  -i string
        Network interface to send packets through (default "eth0")
  -live
        Sniff DHCP packets from the network (default: false)
  -r string
        PCAP file to read from. If not specified, try to send an actual DHCP request
  -s int
        Set the snaplen when using -live (default: 0)
  -to string
        Destination to send packets to. If empty, will use [ff02::1:2]:547
  -unpack
        Unpack inner DHCPv6 messages when parsing relay messages
  -v int
        IP version to use (default 6)
```
