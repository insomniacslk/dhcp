package dhcpv6

// FIXME: rename all the options to have a consistent name, e.g. OPT_<NAME>
const (
	OPTION_CLIENTID                             OptionCode = 1
	OPTION_SERVERID                             OptionCode = 2
	OPTION_IA_NA                                OptionCode = 3
	OPTION_IA_TA                                OptionCode = 4
	OPTION_IAADDR                               OptionCode = 5
	OPTION_ORO                                  OptionCode = 6
	OPTION_PREFERENCE                           OptionCode = 7
	OPTION_ELAPSED_TIME                         OptionCode = 8
	OPTION_RELAY_MSG                            OptionCode = 9
	_                                                      // skip 10
	OPTION_AUTH                                 OptionCode = 11
	OPTION_UNICAST                              OptionCode = 12
	OPTION_STATUS_CODE                          OptionCode = 13
	OPTION_RAPID_COMMIT                         OptionCode = 14
	OPTION_USER_CLASS                           OptionCode = 15
	OPTION_VENDOR_CLASS                         OptionCode = 16
	OPTION_VENDOR_OPTS                          OptionCode = 17
	OPTION_INTERFACE_ID                         OptionCode = 18
	OPTION_RECONF_MSG                           OptionCode = 19
	OPTION_RECONF_ACCEPT                        OptionCode = 20
	SIP_SERVERS_DOMAIN_NAME_LIST                OptionCode = 21
	SIP_SERVERS_IPV6_ADDRESS_LIST               OptionCode = 22
	DNS_RECURSIVE_NAME_SERVER                   OptionCode = 23
	DOMAIN_SEARCH_LIST                          OptionCode = 24
	OPTION_IA_PD                                OptionCode = 25
	OPTION_IAPREFIX                             OptionCode = 26
	OPTION_NIS_SERVERS                          OptionCode = 27
	OPTION_NISP_SERVERS                         OptionCode = 28
	OPTION_NIS_DOMAIN_NAME                      OptionCode = 29
	OPTION_NISP_DOMAIN_NAME                     OptionCode = 30
	SNTP_SERVER_LIST                            OptionCode = 31
	INFORMATION_REFRESH_TIME                    OptionCode = 32
	BCMCS_CONTROLLER_DOMAIN_NAME_LIST           OptionCode = 33
	BCMCS_CONTROLLER_IPV6_ADDRESS_LIST          OptionCode = 34
	_                                                      // skip 35
	OPTION_GEOCONF_CIVIC                        OptionCode = 36
	OPTION_REMOTE_ID                            OptionCode = 37
	RELAY_AGENT_SUBSCRIBER_ID                   OptionCode = 38
	FQDN                                        OptionCode = 39
	PANA_AUTHENTICATION_AGENT                   OptionCode = 40
	OPTION_NEW_POSIX_TIMEZONE                   OptionCode = 41
	OPTION_NEW_TZDB_TIMEZONE                    OptionCode = 42
	ECHO_REQUEST                                OptionCode = 43
	OPTION_LQ_QUERY                             OptionCode = 44
	OPTION_CLIENT_DATA                          OptionCode = 45
	OPTION_CLT_TIME                             OptionCode = 46
	OPTION_LQ_RELAY_DATA                        OptionCode = 47
	OPTION_LQ_CLIENT_LINK                       OptionCode = 48
	MIPV6_HOME_NETWORK_ID_FQDN                  OptionCode = 49
	MIPV6_VISITED_HOME_NETWORK_INFORMATION      OptionCode = 50
	LOST_SERVER                                 OptionCode = 51
	CAPWAP_ACCESS_CONTROLLER_ADDRESSES          OptionCode = 52
	RELAY_ID                                    OptionCode = 53
	OPTION_IPV6_ADDRESS_MOS                     OptionCode = 54
	OPTION_IPV6_FQDN_MOS                        OptionCode = 55
	OPTION_NTP_SERVER                           OptionCode = 56
	OPTION_V6_ACCESS_DOMAIN                     OptionCode = 57
	OPTION_SIP_UA_CS_LIST                       OptionCode = 58
	OPT_BOOTFILE_URL                            OptionCode = 59
	OPT_BOOTFILE_PARAM                          OptionCode = 60
	OPTION_CLIENT_ARCH_TYPE                     OptionCode = 61
	OPTION_NII                                  OptionCode = 62
	OPTION_GEOLOCATION                          OptionCode = 63
	OPTION_AFTR_NAME                            OptionCode = 64
	OPTION_ERP_LOCAL_DOMAIN_NAME                OptionCode = 65
	OPTION_RSOO                                 OptionCode = 66
	OPTION_PD_EXCLUDE                           OptionCode = 67
	VIRTUAL_SUBNET_SELECTION                    OptionCode = 68
	MIPV6_IDENTIFIED_HOME_NETWORK_INFORMATION   OptionCode = 69
	MIPV6_UNRESTRICTED_HOME_NETWORK_INFORMATION OptionCode = 70
	MIPV6_HOME_NETWORK_PREFIX                   OptionCode = 71
	MIPV6_HOME_AGENT_ADDRESS                    OptionCode = 72
	MIPV6_HOME_AGENT_FQDN                       OptionCode = 73
)

var OptionCodeToString = map[OptionCode]string{
	OPTION_CLIENTID:                    "OPTION_CLIENTID",
	OPTION_SERVERID:                    "OPTION_SERVERID",
	OPTION_IA_NA:                       "OPTION_IA_NA",
	OPTION_IA_TA:                       "OPTION_IA_TA",
	OPTION_IAADDR:                      "OPTION_IAADDR",
	OPTION_ORO:                         "OPTION_ORO",
	OPTION_PREFERENCE:                  "OPTION_PREFERENCE",
	OPTION_ELAPSED_TIME:                "OPTION_ELAPSED_TIME",
	OPTION_RELAY_MSG:                   "OPTION_RELAY_MSG",
	OPTION_AUTH:                        "OPTION_AUTH",
	OPTION_UNICAST:                     "OPTION_UNICAST",
	OPTION_STATUS_CODE:                 "OPTION_STATUS_CODE",
	OPTION_RAPID_COMMIT:                "OPTION_RAPID_COMMIT",
	OPTION_USER_CLASS:                  "OPTION_USER_CLASS",
	OPTION_VENDOR_CLASS:                "OPTION_VENDOR_CLASS",
	OPTION_VENDOR_OPTS:                 "OPTION_VENDOR_OPTS",
	OPTION_INTERFACE_ID:                "OPTION_INTERFACE_ID",
	OPTION_RECONF_MSG:                  "OPTION_RECONF_MSG",
	OPTION_RECONF_ACCEPT:               "OPTION_RECONF_ACCEPT",
	SIP_SERVERS_DOMAIN_NAME_LIST:       "SIP Servers Domain Name List",
	SIP_SERVERS_IPV6_ADDRESS_LIST:      "SIP Servers IPv6 Address List",
	DNS_RECURSIVE_NAME_SERVER:          "DNS Recursive Name Server",
	DOMAIN_SEARCH_LIST:                 "Domain Search List",
	OPTION_IA_PD:                       "OPTION_IA_PD",
	OPTION_IAPREFIX:                    "OPTION_IAPREFIX",
	OPTION_NIS_SERVERS:                 "OPTION_NIS_SERVERS",
	OPTION_NISP_SERVERS:                "OPTION_NISP_SERVERS",
	OPTION_NIS_DOMAIN_NAME:             "OPTION_NIS_DOMAIN_NAME",
	OPTION_NISP_DOMAIN_NAME:            "OPTION_NISP_DOMAIN_NAME",
	SNTP_SERVER_LIST:                   "SNTP Server List",
	INFORMATION_REFRESH_TIME:           "Information Refresh Time",
	BCMCS_CONTROLLER_DOMAIN_NAME_LIST:  "BCMCS Controller Domain Name List",
	BCMCS_CONTROLLER_IPV6_ADDRESS_LIST: "BCMCS Controller IPv6 Address List",
	OPTION_GEOCONF_CIVIC:               "OPTION_GEOCONF",
	OPTION_REMOTE_ID:                   "OPTION_REMOTE_ID",
	RELAY_AGENT_SUBSCRIBER_ID:          "Relay-Agent Subscriber ID",
	FQDN: "FQDN",
	PANA_AUTHENTICATION_AGENT:              "PANA Authentication Agent",
	OPTION_NEW_POSIX_TIMEZONE:              "OPTION_NEW_POSIX_TIME_ZONE",
	OPTION_NEW_TZDB_TIMEZONE:               "OPTION_NEW_TZDB_TIMEZONE",
	ECHO_REQUEST:                           "Echo Request",
	OPTION_LQ_QUERY:                        "OPTION_LQ_QUERY",
	OPTION_CLIENT_DATA:                     "OPTION_CLIENT_DATA",
	OPTION_CLT_TIME:                        "OPTION_CLT_TIME",
	OPTION_LQ_RELAY_DATA:                   "OPTION_LQ_RELAY_DATA",
	OPTION_LQ_CLIENT_LINK:                  "OPTION_LQ_CLIENT_LINK",
	MIPV6_HOME_NETWORK_ID_FQDN:             "MIPv6 Home Network ID FQDN",
	MIPV6_VISITED_HOME_NETWORK_INFORMATION: "MIPv6 Visited Home Network Information",
	LOST_SERVER:                            "LoST Server",
	CAPWAP_ACCESS_CONTROLLER_ADDRESSES:     "CAPWAP Access Controller Addresses",
	RELAY_ID:                                    "RELAY_ID",
	OPTION_IPV6_ADDRESS_MOS:                     "OPTION-IPv6_Address-MoS",
	OPTION_IPV6_FQDN_MOS:                        "OPTION-IPv6-FQDN-MoS",
	OPTION_NTP_SERVER:                           "OPTION_NTP_SERVER",
	OPTION_V6_ACCESS_DOMAIN:                     "OPTION_V6_ACCESS_DOMAIN",
	OPTION_SIP_UA_CS_LIST:                       "OPTION_SIP_UA_CS_LIST",
	OPT_BOOTFILE_URL:                            "OPT_BOOTFILE_URL",
	OPT_BOOTFILE_PARAM:                          "OPT_BOOTFILE_PARAM",
	OPTION_CLIENT_ARCH_TYPE:                     "OPTION_CLIENT_ARCH_TYPE",
	OPTION_NII:                                  "OPTION_NII",
	OPTION_GEOLOCATION:                          "OPTION_GEOLOCATION",
	OPTION_AFTR_NAME:                            "OPTION_AFTR_NAME",
	OPTION_ERP_LOCAL_DOMAIN_NAME:                "OPTION_ERP_LOCAL_DOMAIN_NAME",
	OPTION_RSOO:                                 "OPTION_RSOO",
	OPTION_PD_EXCLUDE:                           "OPTION_PD_EXCLUDE",
	VIRTUAL_SUBNET_SELECTION:                    "Virtual Subnet Selection",
	MIPV6_IDENTIFIED_HOME_NETWORK_INFORMATION:   "MIPv6 Identified Home Network Information",
	MIPV6_UNRESTRICTED_HOME_NETWORK_INFORMATION: "MIPv6 Unrestricted Home Network Information",
	MIPV6_HOME_NETWORK_PREFIX:                   "MIPv6 Home Network Prefix",
	MIPV6_HOME_AGENT_ADDRESS:                    "MIPv6 Home Agent Address",
	MIPV6_HOME_AGENT_FQDN:                       "MIPv6 Home Agent FQDN",
}
