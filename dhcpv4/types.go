package dhcpv4

// values from http://www.networksorcery.com/enp/protocol/dhcp.htm and
// http://www.networksorcery.com/enp/protocol/bootp/options.htm

// OpcodeType represents a DHCPv4 opcode.
type OpcodeType uint8

// constants that represent valid values for OpcodeType
const (
	_ OpcodeType = iota // skip 0
	OpcodeBootRequest
	OpcodeBootReply
)

// OpcodeToString maps an OpcodeType to its mnemonic name
var OpcodeToString = map[OpcodeType]string{
	OpcodeBootRequest: "BootRequest",
	OpcodeBootReply:   "BootReply",
}

// DHCPv4 Options
const (
	OptionPad                                        OptionCode = 0
	OptionSubnetMask                                            = 1
	OptionTimeOffset                                            = 2
	OptionRouter                                                = 3
	OptionTimeServer                                            = 4
	OptionNameServer                                            = 5
	OptionDomainNameServer                                      = 6
	OptionLogServer                                             = 7
	OptionQuoteServer                                           = 8
	OptionLPRServer                                             = 9
	OptionImpressServer                                         = 10
	OptionResourceLocationServer                                = 11
	OptionHostName                                              = 12
	OptionBootFileSize                                          = 13
	OptionMeritDumpFile                                         = 14
	OptionDomainName                                            = 15
	OptionSwapServer                                            = 16
	OptionRootPath                                              = 17
	OptionExtensionsPath                                        = 18
	OptionIPForwarding                                          = 19
	OptionNonLocalSourceRouting                                 = 20
	OptionPolicyFilter                                          = 21
	OptionMaximumDatagramAssemblySize                           = 22
	OptionDefaultIPTTL                                          = 23
	OptionPathMTUAgingTimeout                                   = 24
	OptionPathMTUPlateauTable                                   = 25
	OptionInterfaceMTU                                          = 26
	OptionAllSubnetsAreLocal                                    = 27
	OptionBroadcastAddress                                      = 28
	OptionPerformMaskDiscovery                                  = 29
	OptionMaskSupplier                                          = 30
	OptionPerformRouterDiscovery                                = 31
	OptionRouterSolicitationAddress                             = 32
	OptionStaticRoutingTable                                    = 33
	OptionTrailerEncapsulation                                  = 34
	OptionArpCacheTimeout                                       = 35
	OptionEthernetEncapsulation                                 = 36
	OptionDefaulTCPTTL                                          = 37
	OptionTCPKeepaliveInterval                                  = 38
	OptionTCPKeepaliveGarbage                                   = 39
	OptionNetworkInformationServiceDomain                       = 40
	OptionNetworkInformationServers                             = 41
	OptionNTPServers                                            = 42
	OptionVendorSpecificInformation                             = 43
	OptionNetBIOSOverTCPIPNameServer                            = 44
	OptionNetBIOSOverTCPIPDatagramDistributionServer            = 45
	OptionNetBIOSOverTCPIPNodeType                              = 46
	OptionNetBIOSOverTCPIPScope                                 = 47
	OptionXWindowSystemFontServer                               = 48
	OptionXWindowSystemDisplayManger                            = 49
	OptionRequestedIPAddress                                    = 50
	OptionIPAddressLeaseTime                                    = 51
	OptionOptionOverload                                        = 52
	OptionDHCPMessageType                                       = 53
	OptionServerIdentifier                                      = 54
	OptionParameterRequestList                                  = 55
	OptionMessage                                               = 56
	OptionMaximumDHCPMessageSize                                = 57
	OptionRenewTimeValue                                        = 58
	OptionRebindingTimeValue                                    = 59
	OptionClassIdentifier                                       = 60
	OptionClientIdentifier                                      = 61
	OptionNetWareIPDomainName                                   = 62
	OptionNetWareIPInformation                                  = 63
	OptionNetworkInformationServicePlusDomain                   = 64
	OptionNetworkInformationServicePlusServers                  = 65
	OptionTFTPServerName                                        = 66
	OptionBootfileName                                          = 67
	OptionMobileIPHomeAgent                                     = 68
	OptionSimpleMailTransportProtocolServer                     = 69
	OptionPostOfficeProtocolServer                              = 70
	OptionNetworkNewsTransportProtocolServer                    = 71
	OptionDefaultWorldWideWebServer                             = 72
	OptionDefaultFingerServer                                   = 73
	OptionDefaultInternetRelayChatServer                        = 74
	OptionStreetTalkServer                                      = 75
	OptionStreetTalkDirectoryAssistanceServer                   = 76
	OptionUserClassInformation                                  = 77
	OptionSLPDirectoryAgent                                     = 78
	OptionSLPServiceScope                                       = 79
	OptionRapidCommit                                           = 80
	OptionFQDN                                                  = 81
	OptionRelayAgentInformation                                 = 82
	OptionInternetStorageNameService                            = 83
	// Option 84 returned in RFC 3679
	OptionNDSServers                       = 85
	OptionNDSTreeName                      = 86
	OptionNDSContext                       = 87
	OptionBCMCSControllerDomainNameList    = 88
	OptionBCMCSControllerIPv4AddressList   = 89
	OptionAuthentication                   = 90
	OptionClientLastTransactionTime        = 91
	OptionAssociatedIP                     = 92
	OptionClientSystemArchitectureType     = 93
	OptionClientNetworkInterfaceIdentifier = 94
	OptionLDAP                             = 95
	// Option 96 returned in RFC 3679
	OptionClientMachineIdentifier     = 97
	OptionOpenGroupUserAuthentication = 98
	OptionGeoConfCivic                = 99
	OptionIEEE10031TZString           = 100
	OptionReferenceToTZDatabase       = 101
	// Options 102-111 returned in RFC 3679
	OptionNetInfoParentServerAddress = 112
	OptionNetInfoParentServerTag     = 113
	OptionURL                        = 114
	// Option 115 returned in RFC 3679
	OptionAutoConfigure                   = 116
	OptionNameServiceSearch               = 117
	OptionSubnetSelection                 = 118
	OptionDNSDomainSearchList             = 119
	OptionSIPServersDHCPOption            = 120
	OptionClasslessStaticRouteOption      = 121
	OptionCCC                             = 122
	OptionGeoConf                         = 123
	OptionVendorIdentifyingVendorClass    = 124
	OptionVendorIdentifyingVendorSpecific = 125
	// Options 126-127 returned in RFC 3679
	OptionTFTPServerIPAddress                   = 128
	OptionCallServerIPAddress                   = 129
	OptionDiscriminationString                  = 130
	OptionRemoteStatisticsServerIPAddress       = 131
	Option8021PVLANID                           = 132
	Option8021QL2Priority                       = 133
	OptionDiffservCodePoint                     = 134
	OptionHTTPProxyForPhoneSpecificApplications = 135
	OptionPANAAuthenticationAgent               = 136
	OptionLoSTServer                            = 137
	OptionCAPWAPAccessControllerAddresses       = 138
	OptionOPTIONIPv4AddressMoS                  = 139
	OptionOPTIONIPv4FQDNMoS                     = 140
	OptionSIPUAConfigurationServiceDomains      = 141
	OptionOPTIONIPv4AddressANDSF                = 142
	OptionOPTIONIPv6AddressANDSF                = 143
	// Options 144-149 returned in RFC 3679
	OptionTFTPServerAddress = 150
	OptionStatusCode        = 151
	OptionBaseTime          = 152
	OptionStartTimeOfState  = 153
	OptionQueryStartTime    = 154
	OptionQueryEndTime      = 155
	OptionDHCPState         = 156
	OptionDataSource        = 157
	// Options 158-174 returned in RFC 3679
	OptionEtherboot                        = 175
	OptionIPTelephone                      = 176
	OptionEtherbootPacketCableAndCableHome = 177
	// Options 178-207 returned in RFC 3679
	OptionPXELinuxMagicString  = 208
	OptionPXELinuxConfigFile   = 209
	OptionPXELinuxPathPrefix   = 210
	OptionPXELinuxRebootTime   = 211
	OptionOPTION6RD            = 212
	OptionOPTIONv4AccessDomain = 213
	// Options 214-219 returned in RFC 3679
	OptionSubnetAllocation        = 220
	OptionVirtualSubnetAllocation = 221
	// Options 222-223 returned in RFC 3679
	// Options 224-254 are reserved for private use
	OptionEnd = 255
)

// OptionCodeToString maps an OptionCode to its mnemonic name
var OptionCodeToString = map[OptionCode]string{
	OptionPad:                                        "Pad",
	OptionSubnetMask:                                 "Subnet Mask",
	OptionTimeOffset:                                 "Time Offset",
	OptionRouter:                                     "Router",
	OptionTimeServer:                                 "Time Server",
	OptionNameServer:                                 "Name Server",
	OptionDomainNameServer:                           "Domain Name Server",
	OptionLogServer:                                  "Log Server",
	OptionQuoteServer:                                "Quote Server",
	OptionLPRServer:                                  "LPR Server",
	OptionImpressServer:                              "Impress Server",
	OptionResourceLocationServer:                     "Resource Location Server",
	OptionHostName:                                   "Host Name",
	OptionBootFileSize:                               "Boot File Size",
	OptionMeritDumpFile:                              "Merit Dump File",
	OptionDomainName:                                 "Domain Name",
	OptionSwapServer:                                 "Swap Server",
	OptionRootPath:                                   "Root Path",
	OptionExtensionsPath:                             "Extensions Path",
	OptionIPForwarding:                               "IP Forwarding enable/disable",
	OptionNonLocalSourceRouting:                      "Non-local Source Routing enable/disable",
	OptionPolicyFilter:                               "Policy Filter",
	OptionMaximumDatagramAssemblySize:                "Maximum Datagram Reassembly Size",
	OptionDefaultIPTTL:                               "Default IP Time-to-live",
	OptionPathMTUAgingTimeout:                        "Path MTU Aging Timeout",
	OptionPathMTUPlateauTable:                        "Path MTU Plateau Table",
	OptionInterfaceMTU:                               "Interface MTU",
	OptionAllSubnetsAreLocal:                         "All Subnets Are Local",
	OptionBroadcastAddress:                           "Broadcast Address",
	OptionPerformMaskDiscovery:                       "Perform Mask Discovery",
	OptionMaskSupplier:                               "Mask Supplier",
	OptionPerformRouterDiscovery:                     "Perform Router Discovery",
	OptionRouterSolicitationAddress:                  "Router Solicitation Address",
	OptionStaticRoutingTable:                         "Static Routing Table",
	OptionTrailerEncapsulation:                       "Trailer Encapsulation",
	OptionArpCacheTimeout:                            "ARP Cache Timeout",
	OptionEthernetEncapsulation:                      "Ethernet Encapsulation",
	OptionDefaulTCPTTL:                               "Default TCP TTL",
	OptionTCPKeepaliveInterval:                       "TCP Keepalive Interval",
	OptionTCPKeepaliveGarbage:                        "TCP Keepalive Garbage",
	OptionNetworkInformationServiceDomain:            "Network Information Service Domain",
	OptionNetworkInformationServers:                  "Network Information Servers",
	OptionNTPServers:                                 "NTP Servers",
	OptionVendorSpecificInformation:                  "Vendor Specific Information",
	OptionNetBIOSOverTCPIPNameServer:                 "NetBIOS over TCP/IP Name Server",
	OptionNetBIOSOverTCPIPDatagramDistributionServer: "NetBIOS over TCP/IP Datagram Distribution Server",
	OptionNetBIOSOverTCPIPNodeType:                   "NetBIOS over TCP/IP Node Type",
	OptionNetBIOSOverTCPIPScope:                      "NetBIOS over TCP/IP Scope",
	OptionXWindowSystemFontServer:                    "X Window System Font Server",
	OptionXWindowSystemDisplayManger:                 "X Window System Display Manager",
	OptionRequestedIPAddress:                         "Requested IP Address",
	OptionIPAddressLeaseTime:                         "IP Addresses Lease Time",
	OptionOptionOverload:                             "Option Overload",
	OptionDHCPMessageType:                            "DHCP Message Type",
	OptionServerIdentifier:                           "Server Identifier",
	OptionParameterRequestList:                       "Parameter Request List",
	OptionMessage:                                    "Message",
	OptionMaximumDHCPMessageSize:                     "Maximum DHCP Message Size",
	OptionRenewTimeValue:                             "Renew Time Value",
	OptionRebindingTimeValue:                         "Rebinding Time Value",
	OptionClassIdentifier:                            "Class Identifier",
	OptionClientIdentifier:                           "Client identifier",
	OptionNetWareIPDomainName:                        "NetWare/IP Domain Name",
	OptionNetWareIPInformation:                       "NetWare/IP Information",
	OptionNetworkInformationServicePlusDomain:        "Network Information Service+ Domain",
	OptionNetworkInformationServicePlusServers:       "Network Information Service+ Servers",
	OptionTFTPServerName:                             "TFTP Server Name",
	OptionBootfileName:                               "Bootfile Name",
	OptionMobileIPHomeAgent:                          "Mobile IP Home Agent",
	OptionSimpleMailTransportProtocolServer:          "SMTP Server",
	OptionPostOfficeProtocolServer:                   "POP Server",
	OptionNetworkNewsTransportProtocolServer:         "NNTP Server",
	OptionDefaultWorldWideWebServer:                  "Default WWW Server",
	OptionDefaultFingerServer:                        "Default Finger Server",
	OptionDefaultInternetRelayChatServer:             "Default IRC Server",
	OptionStreetTalkServer:                           "StreetTalk Server",
	OptionStreetTalkDirectoryAssistanceServer:        "StreetTalk Directory Assistance Server",
	OptionUserClassInformation:                       "User Class Information",
	OptionSLPDirectoryAgent:                          "SLP DIrectory Agent",
	OptionSLPServiceScope:                            "SLP Service Scope",
	OptionRapidCommit:                                "Rapid Commit",
	OptionFQDN:                                       "FQDN",
	OptionRelayAgentInformation:                      "Relay Agent Information",
	OptionInternetStorageNameService:                 "Internet Storage Name Service",
	// Option 84 returned in RFC 3679
	OptionNDSServers:                       "NDS Servers",
	OptionNDSTreeName:                      "NDS Tree Name",
	OptionNDSContext:                       "NDS Context",
	OptionBCMCSControllerDomainNameList:    "BCMCS Controller Domain Name List",
	OptionBCMCSControllerIPv4AddressList:   "BCMCS Controller IPv4 Address List",
	OptionAuthentication:                   "Authentication",
	OptionClientLastTransactionTime:        "Client Last Transaction Time",
	OptionAssociatedIP:                     "Associated IP",
	OptionClientSystemArchitectureType:     "Client System Architecture Type",
	OptionClientNetworkInterfaceIdentifier: "Client Network Interface Identifier",
	OptionLDAP:                             "LDAP",
	// Option 96 returned in RFC 3679
	OptionClientMachineIdentifier:     "Client Machine Identifier",
	OptionOpenGroupUserAuthentication: "OpenGroup's User Authentication",
	OptionGeoConfCivic:                "GEOCONF_CIVIC",
	OptionIEEE10031TZString:           "IEEE 1003.1 TZ String",
	OptionReferenceToTZDatabase:       "Reference to the TZ Database",
	// Options 102-111 returned in RFC 3679
	OptionNetInfoParentServerAddress: "NetInfo Parent Server Address",
	OptionNetInfoParentServerTag:     "NetInfo Parent Server Tag",
	OptionURL:                        "URL",
	// Option 115 returned in RFC 3679
	OptionAutoConfigure:                   "Auto-Configure",
	OptionNameServiceSearch:               "Name Service Search",
	OptionSubnetSelection:                 "Subnet Selection",
	OptionDNSDomainSearchList:             "DNS Domain Search List",
	OptionSIPServersDHCPOption:            "SIP Servers DHCP Option",
	OptionClasslessStaticRouteOption:      "Classless Static Route Option",
	OptionCCC:                             "CCC, CableLabs Client Configuration",
	OptionGeoConf:                         "GeoConf",
	OptionVendorIdentifyingVendorClass:    "Vendor-Identifying Vendor Class",
	OptionVendorIdentifyingVendorSpecific: "Vendor-Identifying Vendor-Specific",
	// Options 126-127 returned in RFC 3679
	OptionTFTPServerIPAddress:                   "TFTP Server IP Address",
	OptionCallServerIPAddress:                   "Call Server IP Address",
	OptionDiscriminationString:                  "Discrimination String",
	OptionRemoteStatisticsServerIPAddress:       "RemoteStatistics Server IP Address",
	Option8021PVLANID:                           "802.1P VLAN ID",
	Option8021QL2Priority:                       "802.1Q L2 Priority",
	OptionDiffservCodePoint:                     "Diffserv Code Point",
	OptionHTTPProxyForPhoneSpecificApplications: "HTTP Proxy for phone-specific applications",
	OptionPANAAuthenticationAgent:               "PANA Authentication Agent",
	OptionLoSTServer:                            "LoST Server",
	OptionCAPWAPAccessControllerAddresses:       "CAPWAP Access Controller Addresses",
	OptionOPTIONIPv4AddressMoS:                  "OPTION-IPv4_Address-MoS",
	OptionOPTIONIPv4FQDNMoS:                     "OPTION-IPv4_FQDN-MoS",
	OptionSIPUAConfigurationServiceDomains:      "SIP UA Configuration Service Domains",
	OptionOPTIONIPv4AddressANDSF:                "OPTION-IPv4_Address-ANDSF",
	OptionOPTIONIPv6AddressANDSF:                "OPTION-IPv6_Address-ANDSF",
	// Options 144-149 returned in RFC 3679
	OptionTFTPServerAddress: "TFTP Server Address",
	OptionStatusCode:        "Status Code",
	OptionBaseTime:          "Base Time",
	OptionStartTimeOfState:  "Start Time of State",
	OptionQueryStartTime:    "Query Start Time",
	OptionQueryEndTime:      "Query End Time",
	OptionDHCPState:         "DHCP Staet",
	OptionDataSource:        "Data Source",
	// Options 158-174 returned in RFC 3679
	OptionEtherboot:                        "Etherboot",
	OptionIPTelephone:                      "IP Telephone",
	OptionEtherbootPacketCableAndCableHome: "Etherboot / PacketCable and CableHome",
	// Options 178-207 returned in RFC 3679
	OptionPXELinuxMagicString:  "PXELinux Magic String",
	OptionPXELinuxConfigFile:   "PXELinux Config File",
	OptionPXELinuxPathPrefix:   "PXELinux Path Prefix",
	OptionPXELinuxRebootTime:   "PXELinux Reboot Time",
	OptionOPTION6RD:            "OPTION_6RD",
	OptionOPTIONv4AccessDomain: "OPTION_V4_ACCESS_DOMAIN",
	// Options 214-219 returned in RFC 3679
	OptionSubnetAllocation:        "Subnet Allocation",
	OptionVirtualSubnetAllocation: "Virtual Subnet Selection",
	// Options 222-223 returned in RFC 3679
	// Options 224-254 are reserved for private use

	OptionEnd: "End",
}
