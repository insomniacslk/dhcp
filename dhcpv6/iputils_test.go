package dhcpv6

import (
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var ErrDummy = errors.New("dummy error")

type MatchingAddressTestSuite struct {
	suite.Suite
	m mock.Mock

	ips   []net.IP
	addrs []net.Addr
}

func (s *MatchingAddressTestSuite) InterfaceAddresses(name string) ([]net.Addr, error) {
	args := s.m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if ifaddrs, ok := args.Get(0).([]net.Addr); ok {
		return ifaddrs, args.Error(1)
	}
	panic(fmt.Sprintf("assert: arguments: InterfaceAddresses(0) failed because object wasn't correct type: %v", args.Get(0)))
}

func (s *MatchingAddressTestSuite) Match(ip net.IP) bool {
	args := s.m.Called(ip)
	return args.Bool(0)
}

func (s *MatchingAddressTestSuite) SetupTest() {
	InterfaceAddresses = s.InterfaceAddresses
	s.ips = []net.IP{
		net.ParseIP("2401:db00:3020:70e1:face:0:7e:0"),
		net.ParseIP("2803:6080:890c:847e::1"),
		net.ParseIP("fe80::4a57:ddff:fe04:d8e9"),
	}
	s.addrs = []net.Addr{}
	for _, ip := range s.ips {
		s.addrs = append(s.addrs, &net.IPNet{IP: ip})
	}
}

func (s *MatchingAddressTestSuite) TestGetMatchingAddr() {
	// Check if error from InterfaceAddresses immidately returns error
	s.m.On("InterfaceAddresses", "eth0").Return(nil, ErrDummy).Once()
	_, err := getMatchingAddr("eth0", s.Match)
	s.Assert().Equal(ErrDummy, err)
	s.m.AssertExpectations(s.T())
	// Check if the looping is stopped after finding a matching address
	s.m.On("InterfaceAddresses", "eth0").Return(s.addrs, nil).Once()
	s.m.On("Match", s.ips[0]).Return(false).Once()
	s.m.On("Match", s.ips[1]).Return(true).Once()
	ip, err := getMatchingAddr("eth0", s.Match)
	s.Require().NoError(err)
	s.Assert().Equal(s.ips[1], ip)
	s.m.AssertExpectations(s.T())
	// Check if the looping skips not matching addresses
	s.m.On("InterfaceAddresses", "eth0").Return(s.addrs, nil).Once()
	s.m.On("Match", s.ips[0]).Return(false).Once()
	s.m.On("Match", s.ips[1]).Return(false).Once()
	s.m.On("Match", s.ips[2]).Return(true).Once()
	ip, err = getMatchingAddr("eth0", s.Match)
	s.Require().NoError(err)
	s.Assert().Equal(s.ips[2], ip)
	s.m.AssertExpectations(s.T())
	// Check if the error is returned if no matching address is found
	s.m.On("InterfaceAddresses", "eth0").Return(s.addrs, nil).Once()
	s.m.On("Match", s.ips[0]).Return(false).Once()
	s.m.On("Match", s.ips[1]).Return(false).Once()
	s.m.On("Match", s.ips[2]).Return(false).Once()
	_, err = getMatchingAddr("eth0", s.Match)
	s.Assert().EqualError(err, "no matching address found for interface eth0")
	s.m.AssertExpectations(s.T())
}

func (s *MatchingAddressTestSuite) TestGetLinkLocalAddr() {
	s.m.On("InterfaceAddresses", "eth0").Return(s.addrs, nil).Once()
	ip, err := GetLinkLocalAddr("eth0")
	s.Require().NoError(err)
	s.Assert().Equal(s.ips[2], ip)
	s.m.AssertExpectations(s.T())
}

func (s *MatchingAddressTestSuite) TestGetGlobalAddr() {
	s.m.On("InterfaceAddresses", "eth0").Return(s.addrs, nil).Once()
	ip, err := GetGlobalAddr("eth0")
	s.Require().NoError(err)
	s.Assert().Equal(s.ips[0], ip)
	s.m.AssertExpectations(s.T())
}

func TestMatchingAddressTestSuite(t *testing.T) {
	suite.Run(t, new(MatchingAddressTestSuite))
}
