package subnetcalc

import (
	"net/netip"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalcSubnetInfo_InvalidPrefix(t *testing.T) {
	subnetInfo, err := CalcSubnetInfo(netip.Prefix{})
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "invalid prefix")
	assert.Equal(t, SubnetInfo{}, subnetInfo)
}

func TestCalcSubnetInfo_IPv6(t *testing.T) {
	subnetInfo, err := CalcSubnetInfo(netip.MustParsePrefix("2001:db8::/64"))
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "IPv6 not supported yet")
	assert.Equal(t, SubnetInfo{}, subnetInfo)
}

func TestCalcSubnetInfo_SingleIP(t *testing.T) {
	prefix := netip.MustParsePrefix("192.168.1.2/32")
	SubnetInfo, err := CalcSubnetInfo(prefix)
	assert.NoError(t, err)
	assert.Equal(t, SubnetInfo.NetworkAddress, prefix.Addr())
	assert.Equal(t, SubnetInfo.BroadcastIP, prefix.Addr())
	assert.Equal(t, SubnetInfo.SubnetMask, netip.MustParseAddr("255.255.255.255"))
	assert.Equal(t, SubnetInfo.TotalIP, uint(1))
}

func TestCalcSubnetInfo_DefaultRoute(t *testing.T) {
	prefix := netip.MustParsePrefix("0.0.0.0/0")
	SubnetInfo, err := CalcSubnetInfo(prefix)
	assert.NoError(t, err)
	assert.Equal(t, SubnetInfo.NetworkAddress, netip.MustParseAddr("0.0.0.0"))
	assert.Equal(t, SubnetInfo.BroadcastIP, netip.MustParseAddr("255.255.255.255"))
	assert.Equal(t, SubnetInfo.SubnetMask, netip.MustParseAddr("0.0.0.0"))
	assert.Equal(t, SubnetInfo.TotalIP, uint(4294967296))
}

func TestCalcSubnetInfo_LimitedBroadcast(t *testing.T) {
	prefix := netip.MustParsePrefix("255.255.255.255/0")
	SubnetInfo, err := CalcSubnetInfo(prefix)
	assert.NoError(t, err)
	assert.Equal(t, SubnetInfo.NetworkAddress, netip.MustParseAddr("0.0.0.0"))
	assert.Equal(t, SubnetInfo.BroadcastIP, netip.MustParseAddr("255.255.255.255"))
	assert.Equal(t, SubnetInfo.SubnetMask, netip.MustParseAddr("0.0.0.0"))
}

func TestCalcSubnetInfo_Ok(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  SubnetInfo
	}{
		// /0 - Default route
		{
			"0.0.0.0/0",
			"0.0.0.0/0",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("0.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("255.255.255.255"),
				SubnetMask:     netip.MustParseAddr("0.0.0.0"),
				TotalIP:        4294967296,
			},
		},
		{
			"10.0.0.0/0",
			"10.0.0.0/0",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("0.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("255.255.255.255"),
				SubnetMask:     netip.MustParseAddr("0.0.0.0"),
				TotalIP:        4294967296,
			},
		},
		{
			"192.168.1.1/0",
			"192.168.1.1/0",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("0.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("255.255.255.255"),
				SubnetMask:     netip.MustParseAddr("0.0.0.0"),
				TotalIP:        4294967296,
			},
		},
		{
			"8.8.8.8/0",
			"8.8.8.8/0",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("0.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("255.255.255.255"),
				SubnetMask:     netip.MustParseAddr("0.0.0.0"),
				TotalIP:        4294967296,
			},
		},

		// /1
		{
			"0.0.0.0/1",
			"0.0.0.0/1",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("0.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("127.255.255.255"),
				SubnetMask:     netip.MustParseAddr("128.0.0.0"),
				TotalIP:        2147483648,
			},
		},
		{
			"128.0.0.0/1",
			"128.0.0.0/1",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("128.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("255.255.255.255"),
				SubnetMask:     netip.MustParseAddr("128.0.0.0"),
				TotalIP:        2147483648,
			},
		},
		{
			"10.0.0.0/1",
			"10.0.0.0/1",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("0.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("127.255.255.255"),
				SubnetMask:     netip.MustParseAddr("128.0.0.0"),
				TotalIP:        2147483648,
			},
		},
		{
			"192.168.0.0/1",
			"192.168.0.0/1",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("128.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("255.255.255.255"),
				SubnetMask:     netip.MustParseAddr("128.0.0.0"),
				TotalIP:        2147483648,
			},
		},

		// /2
		{
			"0.0.0.0/2",
			"0.0.0.0/2",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("0.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("63.255.255.255"),
				SubnetMask:     netip.MustParseAddr("192.0.0.0"),
				TotalIP:        1073741824,
			},
		},
		{
			"64.0.0.0/2",
			"64.0.0.0/2",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("64.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("127.255.255.255"),
				SubnetMask:     netip.MustParseAddr("192.0.0.0"),
				TotalIP:        1073741824,
			},
		},
		{
			"128.0.0.0/2",
			"128.0.0.0/2",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("128.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("191.255.255.255"),
				SubnetMask:     netip.MustParseAddr("192.0.0.0"),
				TotalIP:        1073741824,
			},
		},
		{
			"192.0.0.0/2",
			"192.0.0.0/2",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("192.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("255.255.255.255"),
				SubnetMask:     netip.MustParseAddr("192.0.0.0"),
				TotalIP:        1073741824,
			},
		},

		// /3
		{
			"0.0.0.0/3",
			"0.0.0.0/3",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("0.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("31.255.255.255"),
				SubnetMask:     netip.MustParseAddr("224.0.0.0"),
				TotalIP:        536870912,
			},
		},
		{
			"32.0.0.0/3",
			"32.0.0.0/3",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("32.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("63.255.255.255"),
				SubnetMask:     netip.MustParseAddr("224.0.0.0"),
				TotalIP:        536870912,
			},
		},
		{
			"10.0.0.0/3",
			"10.0.0.0/3",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("0.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("31.255.255.255"),
				SubnetMask:     netip.MustParseAddr("224.0.0.0"),
				TotalIP:        536870912,
			},
		},
		{
			"192.0.0.0/3",
			"192.0.0.0/3",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("192.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("223.255.255.255"),
				SubnetMask:     netip.MustParseAddr("224.0.0.0"),
				TotalIP:        536870912,
			},
		},

		// /4
		{
			"0.0.0.0/4",
			"0.0.0.0/4",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("0.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("15.255.255.255"),
				SubnetMask:     netip.MustParseAddr("240.0.0.0"),
				TotalIP:        268435456,
			},
		},
		{
			"16.0.0.0/4",
			"16.0.0.0/4",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("16.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("31.255.255.255"),
				SubnetMask:     netip.MustParseAddr("240.0.0.0"),
				TotalIP:        268435456,
			},
		},
		{
			"10.0.0.0/4",
			"10.0.0.0/4",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("0.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("15.255.255.255"),
				SubnetMask:     netip.MustParseAddr("240.0.0.0"),
				TotalIP:        268435456,
			},
		},
		{
			"192.168.0.0/4",
			"192.168.0.0/4",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("192.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("207.255.255.255"),
				SubnetMask:     netip.MustParseAddr("240.0.0.0"),
				TotalIP:        268435456,
			},
		},

		// /5
		{
			"0.0.0.0/5",
			"0.0.0.0/5",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("0.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("7.255.255.255"),
				SubnetMask:     netip.MustParseAddr("248.0.0.0"),
				TotalIP:        134217728,
			},
		},
		{
			"8.0.0.0/5",
			"8.0.0.0/5",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("8.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("15.255.255.255"),
				SubnetMask:     netip.MustParseAddr("248.0.0.0"),
				TotalIP:        134217728,
			},
		},
		{
			"10.0.0.0/5",
			"10.0.0.0/5",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("8.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("15.255.255.255"),
				SubnetMask:     netip.MustParseAddr("248.0.0.0"),
				TotalIP:        134217728,
			},
		},
		{
			"172.16.0.0/5",
			"172.16.0.0/5",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("168.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("175.255.255.255"),
				SubnetMask:     netip.MustParseAddr("248.0.0.0"),
				TotalIP:        134217728,
			},
		},

		// /6
		{
			"0.0.0.0/6",
			"0.0.0.0/6",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("0.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("3.255.255.255"),
				SubnetMask:     netip.MustParseAddr("252.0.0.0"),
				TotalIP:        67108864,
			},
		},
		{
			"4.0.0.0/6",
			"4.0.0.0/6",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("4.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("7.255.255.255"),
				SubnetMask:     netip.MustParseAddr("252.0.0.0"),
				TotalIP:        67108864,
			},
		},
		{
			"8.0.0.0/6",
			"8.0.0.0/6",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("8.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("11.255.255.255"),
				SubnetMask:     netip.MustParseAddr("252.0.0.0"),
				TotalIP:        67108864,
			},
		},
		{
			"10.0.0.0/6",
			"10.0.0.0/6",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("8.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("11.255.255.255"),
				SubnetMask:     netip.MustParseAddr("252.0.0.0"),
				TotalIP:        67108864,
			},
		},

		// /7
		{
			"0.0.0.0/7",
			"0.0.0.0/7",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("0.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("1.255.255.255"),
				SubnetMask:     netip.MustParseAddr("254.0.0.0"),
				TotalIP:        33554432,
			},
		},
		{
			"2.0.0.0/7",
			"2.0.0.0/7",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("2.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("3.255.255.255"),
				SubnetMask:     netip.MustParseAddr("254.0.0.0"),
				TotalIP:        33554432,
			},
		},
		{
			"8.0.0.0/7",
			"8.0.0.0/7",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("8.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("9.255.255.255"),
				SubnetMask:     netip.MustParseAddr("254.0.0.0"),
				TotalIP:        33554432,
			},
		},
		{
			"10.0.0.0/7",
			"10.0.0.0/7",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("11.255.255.255"),
				SubnetMask:     netip.MustParseAddr("254.0.0.0"),
				TotalIP:        33554432,
			},
		},

		// /8 - Class A
		{
			"0.0.0.0/8",
			"0.0.0.0/8",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("0.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("0.255.255.255"),
				SubnetMask:     netip.MustParseAddr("255.0.0.0"),
				TotalIP:        16777216,
			},
		},
		{
			"10.0.0.0/8",
			"10.0.0.0/8",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("10.255.255.255"),
				SubnetMask:     netip.MustParseAddr("255.0.0.0"),
				TotalIP:        16777216,
			},
		},
		{
			"8.8.8.8/8",
			"8.8.8.8/8",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("8.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("8.255.255.255"),
				SubnetMask:     netip.MustParseAddr("255.0.0.0"),
				TotalIP:        16777216,
			},
		},
		{
			"172.16.0.0/8",
			"172.16.0.0/8",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("172.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("172.255.255.255"),
				SubnetMask:     netip.MustParseAddr("255.0.0.0"),
				TotalIP:        16777216,
			},
		},

		// /9
		{
			"0.0.0.0/9",
			"0.0.0.0/9",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("0.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("0.127.255.255"),
				SubnetMask:     netip.MustParseAddr("255.128.0.0"),
				TotalIP:        8388608,
			},
		},
		{
			"10.0.0.0/9",
			"10.0.0.0/9",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("10.127.255.255"),
				SubnetMask:     netip.MustParseAddr("255.128.0.0"),
				TotalIP:        8388608,
			},
		},
		{
			"10.128.0.0/9",
			"10.128.0.0/9",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.128.0.0"),
				BroadcastIP:    netip.MustParseAddr("10.255.255.255"),
				SubnetMask:     netip.MustParseAddr("255.128.0.0"),
				TotalIP:        8388608,
			},
		},
		{
			"192.0.0.0/9",
			"192.0.0.0/9",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("192.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("192.127.255.255"),
				SubnetMask:     netip.MustParseAddr("255.128.0.0"),
				TotalIP:        8388608,
			},
		},

		// /10
		{
			"10.0.0.0/10",
			"10.0.0.0/10",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("10.63.255.255"),
				SubnetMask:     netip.MustParseAddr("255.192.0.0"),
				TotalIP:        4194304,
			},
		},
		{
			"10.64.0.0/10",
			"10.64.0.0/10",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.64.0.0"),
				BroadcastIP:    netip.MustParseAddr("10.127.255.255"),
				SubnetMask:     netip.MustParseAddr("255.192.0.0"),
				TotalIP:        4194304,
			},
		},
		{
			"172.16.0.0/10",
			"172.16.0.0/10",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("172.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("172.63.255.255"),
				SubnetMask:     netip.MustParseAddr("255.192.0.0"),
				TotalIP:        4194304,
			},
		},
		{
			"1.0.0.0/10",
			"1.0.0.0/10",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("1.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("1.63.255.255"),
				SubnetMask:     netip.MustParseAddr("255.192.0.0"),
				TotalIP:        4194304,
			},
		},

		// /11
		{
			"10.0.0.0/11",
			"10.0.0.0/11",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("10.31.255.255"),
				SubnetMask:     netip.MustParseAddr("255.224.0.0"),
				TotalIP:        2097152,
			},
		},
		{
			"10.32.0.0/11",
			"10.32.0.0/11",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.32.0.0"),
				BroadcastIP:    netip.MustParseAddr("10.63.255.255"),
				SubnetMask:     netip.MustParseAddr("255.224.0.0"),
				TotalIP:        2097152,
			},
		},
		{
			"172.16.0.0/11",
			"172.16.0.0/11",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("172.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("172.31.255.255"),
				SubnetMask:     netip.MustParseAddr("255.224.0.0"),
				TotalIP:        2097152,
			},
		},
		{
			"8.0.0.0/11",
			"8.0.0.0/11",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("8.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("8.31.255.255"),
				SubnetMask:     netip.MustParseAddr("255.224.0.0"),
				TotalIP:        2097152,
			},
		},

		// /12
		{
			"10.0.0.0/12",
			"10.0.0.0/12",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("10.15.255.255"),
				SubnetMask:     netip.MustParseAddr("255.240.0.0"),
				TotalIP:        1048576,
			},
		},
		{
			"10.16.0.0/12",
			"10.16.0.0/12",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.16.0.0"),
				BroadcastIP:    netip.MustParseAddr("10.31.255.255"),
				SubnetMask:     netip.MustParseAddr("255.240.0.0"),
				TotalIP:        1048576,
			},
		},
		{
			"172.16.0.0/12",
			"172.16.0.0/12",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("172.16.0.0"),
				BroadcastIP:    netip.MustParseAddr("172.31.255.255"),
				SubnetMask:     netip.MustParseAddr("255.240.0.0"),
				TotalIP:        1048576,
			},
		},
		{
			"1.0.0.0/12",
			"1.0.0.0/12",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("1.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("1.15.255.255"),
				SubnetMask:     netip.MustParseAddr("255.240.0.0"),
				TotalIP:        1048576,
			},
		},

		// /13
		{
			"10.0.0.0/13",
			"10.0.0.0/13",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("10.7.255.255"),
				SubnetMask:     netip.MustParseAddr("255.248.0.0"),
				TotalIP:        524288,
			},
		},
		{
			"10.8.0.0/13",
			"10.8.0.0/13",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.8.0.0"),
				BroadcastIP:    netip.MustParseAddr("10.15.255.255"),
				SubnetMask:     netip.MustParseAddr("255.248.0.0"),
				TotalIP:        524288,
			},
		},
		{
			"172.16.0.0/13",
			"172.16.0.0/13",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("172.16.0.0"),
				BroadcastIP:    netip.MustParseAddr("172.23.255.255"),
				SubnetMask:     netip.MustParseAddr("255.248.0.0"),
				TotalIP:        524288,
			},
		},
		{
			"8.8.0.0/13",
			"8.8.0.0/13",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("8.8.0.0"),
				BroadcastIP:    netip.MustParseAddr("8.15.255.255"),
				SubnetMask:     netip.MustParseAddr("255.248.0.0"),
				TotalIP:        524288,
			},
		},

		// /14
		{
			"10.0.0.0/14",
			"10.0.0.0/14",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("10.3.255.255"),
				SubnetMask:     netip.MustParseAddr("255.252.0.0"),
				TotalIP:        262144,
			},
		},
		{
			"10.4.0.0/14",
			"10.4.0.0/14",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.4.0.0"),
				BroadcastIP:    netip.MustParseAddr("10.7.255.255"),
				SubnetMask:     netip.MustParseAddr("255.252.0.0"),
				TotalIP:        262144,
			},
		},
		{
			"172.16.0.0/14",
			"172.16.0.0/14",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("172.16.0.0"),
				BroadcastIP:    netip.MustParseAddr("172.19.255.255"),
				SubnetMask:     netip.MustParseAddr("255.252.0.0"),
				TotalIP:        262144,
			},
		},
		{
			"1.0.0.0/14",
			"1.0.0.0/14",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("1.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("1.3.255.255"),
				SubnetMask:     netip.MustParseAddr("255.252.0.0"),
				TotalIP:        262144,
			},
		},

		// /15
		{
			"10.0.0.0/15",
			"10.0.0.0/15",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("10.1.255.255"),
				SubnetMask:     netip.MustParseAddr("255.254.0.0"),
				TotalIP:        131072,
			},
		},
		{
			"10.2.0.0/15",
			"10.2.0.0/15",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.2.0.0"),
				BroadcastIP:    netip.MustParseAddr("10.3.255.255"),
				SubnetMask:     netip.MustParseAddr("255.254.0.0"),
				TotalIP:        131072,
			},
		},
		{
			"172.16.0.0/15",
			"172.16.0.0/15",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("172.16.0.0"),
				BroadcastIP:    netip.MustParseAddr("172.17.255.255"),
				SubnetMask:     netip.MustParseAddr("255.254.0.0"),
				TotalIP:        131072,
			},
		},
		{
			"8.8.0.0/15",
			"8.8.0.0/15",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("8.8.0.0"),
				BroadcastIP:    netip.MustParseAddr("8.9.255.255"),
				SubnetMask:     netip.MustParseAddr("255.254.0.0"),
				TotalIP:        131072,
			},
		},

		// /16 - Class B
		{
			"10.0.0.0/16",
			"10.0.0.0/16",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("10.0.255.255"),
				SubnetMask:     netip.MustParseAddr("255.255.0.0"),
				TotalIP:        65536,
			},
		},
		{
			"10.1.0.0/16",
			"10.1.0.0/16",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.1.0.0"),
				BroadcastIP:    netip.MustParseAddr("10.1.255.255"),
				SubnetMask:     netip.MustParseAddr("255.255.0.0"),
				TotalIP:        65536,
			},
		},
		{
			"172.16.0.0/16",
			"172.16.0.0/16",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("172.16.0.0"),
				BroadcastIP:    netip.MustParseAddr("172.16.255.255"),
				SubnetMask:     netip.MustParseAddr("255.255.0.0"),
				TotalIP:        65536,
			},
		},
		{
			"192.168.0.0/16",
			"192.168.0.0/16",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("192.168.0.0"),
				BroadcastIP:    netip.MustParseAddr("192.168.255.255"),
				SubnetMask:     netip.MustParseAddr("255.255.0.0"),
				TotalIP:        65536,
			},
		},
		{
			"8.8.0.0/16",
			"8.8.0.0/16",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("8.8.0.0"),
				BroadcastIP:    netip.MustParseAddr("8.8.255.255"),
				SubnetMask:     netip.MustParseAddr("255.255.0.0"),
				TotalIP:        65536,
			},
		},

		// /17
		{
			"10.0.0.0/17",
			"10.0.0.0/17",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("10.0.127.255"),
				SubnetMask:     netip.MustParseAddr("255.255.128.0"),
				TotalIP:        32768,
			},
		},
		{
			"10.0.128.0/17",
			"10.0.128.0/17",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.0.128.0"),
				BroadcastIP:    netip.MustParseAddr("10.0.255.255"),
				SubnetMask:     netip.MustParseAddr("255.255.128.0"),
				TotalIP:        32768,
			},
		},
		{
			"172.16.0.0/17",
			"172.16.0.0/17",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("172.16.0.0"),
				BroadcastIP:    netip.MustParseAddr("172.16.127.255"),
				SubnetMask:     netip.MustParseAddr("255.255.128.0"),
				TotalIP:        32768,
			},
		},
		{
			"192.168.0.0/17",
			"192.168.0.0/17",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("192.168.0.0"),
				BroadcastIP:    netip.MustParseAddr("192.168.127.255"),
				SubnetMask:     netip.MustParseAddr("255.255.128.0"),
				TotalIP:        32768,
			},
		},

		// /18
		{
			"10.0.0.0/18",
			"10.0.0.0/18",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("10.0.63.255"),
				SubnetMask:     netip.MustParseAddr("255.255.192.0"),
				TotalIP:        16384,
			},
		},
		{
			"10.0.64.0/18",
			"10.0.64.0/18",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.0.64.0"),
				BroadcastIP:    netip.MustParseAddr("10.0.127.255"),
				SubnetMask:     netip.MustParseAddr("255.255.192.0"),
				TotalIP:        16384,
			},
		},
		{
			"172.16.0.0/18",
			"172.16.0.0/18",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("172.16.0.0"),
				BroadcastIP:    netip.MustParseAddr("172.16.63.255"),
				SubnetMask:     netip.MustParseAddr("255.255.192.0"),
				TotalIP:        16384,
			},
		},
		{
			"192.168.0.0/18",
			"192.168.0.0/18",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("192.168.0.0"),
				BroadcastIP:    netip.MustParseAddr("192.168.63.255"),
				SubnetMask:     netip.MustParseAddr("255.255.192.0"),
				TotalIP:        16384,
			},
		},

		// /19
		{
			"10.0.0.0/19",
			"10.0.0.0/19",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("10.0.31.255"),
				SubnetMask:     netip.MustParseAddr("255.255.224.0"),
				TotalIP:        8192,
			},
		},
		{
			"10.0.32.0/19",
			"10.0.32.0/19",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.0.32.0"),
				BroadcastIP:    netip.MustParseAddr("10.0.63.255"),
				SubnetMask:     netip.MustParseAddr("255.255.224.0"),
				TotalIP:        8192,
			},
		},
		{
			"172.16.0.0/19",
			"172.16.0.0/19",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("172.16.0.0"),
				BroadcastIP:    netip.MustParseAddr("172.16.31.255"),
				SubnetMask:     netip.MustParseAddr("255.255.224.0"),
				TotalIP:        8192,
			},
		},
		{
			"192.168.0.0/19",
			"192.168.0.0/19",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("192.168.0.0"),
				BroadcastIP:    netip.MustParseAddr("192.168.31.255"),
				SubnetMask:     netip.MustParseAddr("255.255.224.0"),
				TotalIP:        8192,
			},
		},

		// /20
		{
			"10.0.0.0/20",
			"10.0.0.0/20",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("10.0.15.255"),
				SubnetMask:     netip.MustParseAddr("255.255.240.0"),
				TotalIP:        4096,
			},
		},
		{
			"10.0.16.0/20",
			"10.0.16.0/20",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.0.16.0"),
				BroadcastIP:    netip.MustParseAddr("10.0.31.255"),
				SubnetMask:     netip.MustParseAddr("255.255.240.0"),
				TotalIP:        4096,
			},
		},
		{
			"172.16.0.0/20",
			"172.16.0.0/20",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("172.16.0.0"),
				BroadcastIP:    netip.MustParseAddr("172.16.15.255"),
				SubnetMask:     netip.MustParseAddr("255.255.240.0"),
				TotalIP:        4096,
			},
		},
		{
			"192.168.0.0/20",
			"192.168.0.0/20",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("192.168.0.0"),
				BroadcastIP:    netip.MustParseAddr("192.168.15.255"),
				SubnetMask:     netip.MustParseAddr("255.255.240.0"),
				TotalIP:        4096,
			},
		},

		// /21
		{
			"10.0.0.0/21",
			"10.0.0.0/21",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("10.0.7.255"),
				SubnetMask:     netip.MustParseAddr("255.255.248.0"),
				TotalIP:        2048,
			},
		},
		{
			"10.0.8.0/21",
			"10.0.8.0/21",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.0.8.0"),
				BroadcastIP:    netip.MustParseAddr("10.0.15.255"),
				SubnetMask:     netip.MustParseAddr("255.255.248.0"),
				TotalIP:        2048,
			},
		},
		{
			"172.16.0.0/21",
			"172.16.0.0/21",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("172.16.0.0"),
				BroadcastIP:    netip.MustParseAddr("172.16.7.255"),
				SubnetMask:     netip.MustParseAddr("255.255.248.0"),
				TotalIP:        2048,
			},
		},
		{
			"192.168.0.0/21",
			"192.168.0.0/21",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("192.168.0.0"),
				BroadcastIP:    netip.MustParseAddr("192.168.7.255"),
				SubnetMask:     netip.MustParseAddr("255.255.248.0"),
				TotalIP:        2048,
			},
		},

		// /22
		{
			"10.0.0.0/22",
			"10.0.0.0/22",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("10.0.3.255"),
				SubnetMask:     netip.MustParseAddr("255.255.252.0"),
				TotalIP:        1024,
			},
		},
		{
			"10.0.4.0/22",
			"10.0.4.0/22",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.0.4.0"),
				BroadcastIP:    netip.MustParseAddr("10.0.7.255"),
				SubnetMask:     netip.MustParseAddr("255.255.252.0"),
				TotalIP:        1024,
			},
		},
		{
			"172.16.0.0/22",
			"172.16.0.0/22",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("172.16.0.0"),
				BroadcastIP:    netip.MustParseAddr("172.16.3.255"),
				SubnetMask:     netip.MustParseAddr("255.255.252.0"),
				TotalIP:        1024,
			},
		},
		{
			"192.168.0.0/22",
			"192.168.0.0/22",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("192.168.0.0"),
				BroadcastIP:    netip.MustParseAddr("192.168.3.255"),
				SubnetMask:     netip.MustParseAddr("255.255.252.0"),
				TotalIP:        1024,
			},
		},

		// /23
		{
			"10.0.0.0/23",
			"10.0.0.0/23",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("10.0.1.255"),
				SubnetMask:     netip.MustParseAddr("255.255.254.0"),
				TotalIP:        512,
			},
		},
		{
			"10.0.2.0/23",
			"10.0.2.0/23",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.0.2.0"),
				BroadcastIP:    netip.MustParseAddr("10.0.3.255"),
				SubnetMask:     netip.MustParseAddr("255.255.254.0"),
				TotalIP:        512,
			},
		},
		{
			"172.16.0.0/23",
			"172.16.0.0/23",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("172.16.0.0"),
				BroadcastIP:    netip.MustParseAddr("172.16.1.255"),
				SubnetMask:     netip.MustParseAddr("255.255.254.0"),
				TotalIP:        512,
			},
		},
		{
			"192.168.0.0/23",
			"192.168.0.0/23",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("192.168.0.0"),
				BroadcastIP:    netip.MustParseAddr("192.168.1.255"),
				SubnetMask:     netip.MustParseAddr("255.255.254.0"),
				TotalIP:        512,
			},
		},

		// /24 - Class C
		{
			"10.0.0.0/24",
			"10.0.0.0/24",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("10.0.0.255"),
				SubnetMask:     netip.MustParseAddr("255.255.255.0"),
				TotalIP:        256,
			},
		},
		{
			"10.1.1.0/24",
			"10.1.1.0/24",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.1.1.0"),
				BroadcastIP:    netip.MustParseAddr("10.1.1.255"),
				SubnetMask:     netip.MustParseAddr("255.255.255.0"),
				TotalIP:        256,
			},
		},
		{
			"172.16.0.0/24",
			"172.16.0.0/24",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("172.16.0.0"),
				BroadcastIP:    netip.MustParseAddr("172.16.0.255"),
				SubnetMask:     netip.MustParseAddr("255.255.255.0"),
				TotalIP:        256,
			},
		},
		{
			"192.168.1.0/24",
			"192.168.1.0/24",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("192.168.1.0"),
				BroadcastIP:    netip.MustParseAddr("192.168.1.255"),
				SubnetMask:     netip.MustParseAddr("255.255.255.0"),
				TotalIP:        256,
			},
		},
		{
			"8.8.8.0/24",
			"8.8.8.0/24",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("8.8.8.0"),
				BroadcastIP:    netip.MustParseAddr("8.8.8.255"),
				SubnetMask:     netip.MustParseAddr("255.255.255.0"),
				TotalIP:        256,
			},
		},

		// /25
		{
			"10.0.0.0/25",
			"10.0.0.0/25",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("10.0.0.127"),
				SubnetMask:     netip.MustParseAddr("255.255.255.128"),
				TotalIP:        128,
			},
		},
		{
			"10.0.0.128/25",
			"10.0.0.128/25",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.0.0.128"),
				BroadcastIP:    netip.MustParseAddr("10.0.0.255"),
				SubnetMask:     netip.MustParseAddr("255.255.255.128"),
				TotalIP:        128,
			},
		},
		{
			"172.16.0.0/25",
			"172.16.0.0/25",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("172.16.0.0"),
				BroadcastIP:    netip.MustParseAddr("172.16.0.127"),
				SubnetMask:     netip.MustParseAddr("255.255.255.128"),
				TotalIP:        128,
			},
		},
		{
			"192.168.1.0/25",
			"192.168.1.0/25",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("192.168.1.0"),
				BroadcastIP:    netip.MustParseAddr("192.168.1.127"),
				SubnetMask:     netip.MustParseAddr("255.255.255.128"),
				TotalIP:        128,
			},
		},

		// /26
		{
			"10.0.0.0/26",
			"10.0.0.0/26",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("10.0.0.63"),
				SubnetMask:     netip.MustParseAddr("255.255.255.192"),
				TotalIP:        64,
			},
		},
		{
			"10.0.0.64/26",
			"10.0.0.64/26",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.0.0.64"),
				BroadcastIP:    netip.MustParseAddr("10.0.0.127"),
				SubnetMask:     netip.MustParseAddr("255.255.255.192"),
				TotalIP:        64,
			},
		},
		{
			"172.16.0.0/26",
			"172.16.0.0/26",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("172.16.0.0"),
				BroadcastIP:    netip.MustParseAddr("172.16.0.63"),
				SubnetMask:     netip.MustParseAddr("255.255.255.192"),
				TotalIP:        64,
			},
		},
		{
			"192.168.1.0/26",
			"192.168.1.0/26",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("192.168.1.0"),
				BroadcastIP:    netip.MustParseAddr("192.168.1.63"),
				SubnetMask:     netip.MustParseAddr("255.255.255.192"),
				TotalIP:        64,
			},
		},

		// /27
		{
			"10.0.0.0/27",
			"10.0.0.0/27",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("10.0.0.31"),
				SubnetMask:     netip.MustParseAddr("255.255.255.224"),
				TotalIP:        32,
			},
		},
		{
			"10.0.0.32/27",
			"10.0.0.32/27",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.0.0.32"),
				BroadcastIP:    netip.MustParseAddr("10.0.0.63"),
				SubnetMask:     netip.MustParseAddr("255.255.255.224"),
				TotalIP:        32,
			},
		},
		{
			"172.16.0.0/27",
			"172.16.0.0/27",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("172.16.0.0"),
				BroadcastIP:    netip.MustParseAddr("172.16.0.31"),
				SubnetMask:     netip.MustParseAddr("255.255.255.224"),
				TotalIP:        32,
			},
		},
		{
			"192.168.1.0/27",
			"192.168.1.0/27",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("192.168.1.0"),
				BroadcastIP:    netip.MustParseAddr("192.168.1.31"),
				SubnetMask:     netip.MustParseAddr("255.255.255.224"),
				TotalIP:        32,
			},
		},

		// /28
		{
			"10.0.0.0/28",
			"10.0.0.0/28",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("10.0.0.15"),
				SubnetMask:     netip.MustParseAddr("255.255.255.240"),
				TotalIP:        16,
			},
		},
		{
			"10.0.0.16/28",
			"10.0.0.16/28",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.0.0.16"),
				BroadcastIP:    netip.MustParseAddr("10.0.0.31"),
				SubnetMask:     netip.MustParseAddr("255.255.255.240"),
				TotalIP:        16,
			},
		},
		{
			"172.16.0.0/28",
			"172.16.0.0/28",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("172.16.0.0"),
				BroadcastIP:    netip.MustParseAddr("172.16.0.15"),
				SubnetMask:     netip.MustParseAddr("255.255.255.240"),
				TotalIP:        16,
			},
		},
		{
			"192.168.1.0/28",
			"192.168.1.0/28",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("192.168.1.0"),
				BroadcastIP:    netip.MustParseAddr("192.168.1.15"),
				SubnetMask:     netip.MustParseAddr("255.255.255.240"),
				TotalIP:        16,
			},
		},

		// /29
		{
			"10.0.0.0/29",
			"10.0.0.0/29",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("10.0.0.7"),
				SubnetMask:     netip.MustParseAddr("255.255.255.248"),
				TotalIP:        8,
			},
		},
		{
			"10.0.0.8/29",
			"10.0.0.8/29",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.0.0.8"),
				BroadcastIP:    netip.MustParseAddr("10.0.0.15"),
				SubnetMask:     netip.MustParseAddr("255.255.255.248"),
				TotalIP:        8,
			},
		},
		{
			"172.16.0.0/29",
			"172.16.0.0/29",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("172.16.0.0"),
				BroadcastIP:    netip.MustParseAddr("172.16.0.7"),
				SubnetMask:     netip.MustParseAddr("255.255.255.248"),
				TotalIP:        8,
			},
		},
		{
			"192.168.1.0/29",
			"192.168.1.0/29",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("192.168.1.0"),
				BroadcastIP:    netip.MustParseAddr("192.168.1.7"),
				SubnetMask:     netip.MustParseAddr("255.255.255.248"),
				TotalIP:        8,
			},
		},

		// /30 - Point-to-point
		{
			"10.0.0.0/30",
			"10.0.0.0/30",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("10.0.0.3"),
				SubnetMask:     netip.MustParseAddr("255.255.255.252"),
				TotalIP:        4,
			},
		},
		{
			"10.0.0.4/30",
			"10.0.0.4/30",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.0.0.4"),
				BroadcastIP:    netip.MustParseAddr("10.0.0.7"),
				SubnetMask:     netip.MustParseAddr("255.255.255.252"),
				TotalIP:        4,
			},
		},
		{
			"172.16.0.0/30",
			"172.16.0.0/30",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("172.16.0.0"),
				BroadcastIP:    netip.MustParseAddr("172.16.0.3"),
				SubnetMask:     netip.MustParseAddr("255.255.255.252"),
				TotalIP:        4,
			},
		},
		{
			"192.168.1.0/30",
			"192.168.1.0/30",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("192.168.1.0"),
				BroadcastIP:    netip.MustParseAddr("192.168.1.3"),
				SubnetMask:     netip.MustParseAddr("255.255.255.252"),
				TotalIP:        4,
			},
		},

		// /31 - Point-to-point (RFC 3021)
		{
			"10.0.0.0/31",
			"10.0.0.0/31",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("10.0.0.1"),
				SubnetMask:     netip.MustParseAddr("255.255.255.254"),
				TotalIP:        2,
			},
		},
		{
			"10.0.0.2/31",
			"10.0.0.2/31",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.0.0.2"),
				BroadcastIP:    netip.MustParseAddr("10.0.0.3"),
				SubnetMask:     netip.MustParseAddr("255.255.255.254"),
				TotalIP:        2,
			},
		},
		{
			"172.16.0.0/31",
			"172.16.0.0/31",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("172.16.0.0"),
				BroadcastIP:    netip.MustParseAddr("172.16.0.1"),
				SubnetMask:     netip.MustParseAddr("255.255.255.254"),
				TotalIP:        2,
			},
		},
		{
			"192.168.1.0/31",
			"192.168.1.0/31",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("192.168.1.0"),
				BroadcastIP:    netip.MustParseAddr("192.168.1.1"),
				SubnetMask:     netip.MustParseAddr("255.255.255.254"),
				TotalIP:        2,
			},
		},

		// /32 - Single host
		{
			"10.0.0.1/32",
			"10.0.0.1/32",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.0.0.1"),
				BroadcastIP:    netip.MustParseAddr("10.0.0.1"),
				SubnetMask:     netip.MustParseAddr("255.255.255.255"),
				TotalIP:        1,
			},
		},
		{
			"10.0.0.2/32",
			"10.0.0.2/32",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.0.0.2"),
				BroadcastIP:    netip.MustParseAddr("10.0.0.2"),
				SubnetMask:     netip.MustParseAddr("255.255.255.255"),
				TotalIP:        1,
			},
		},
		{
			"172.16.0.1/32",
			"172.16.0.1/32",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("172.16.0.1"),
				BroadcastIP:    netip.MustParseAddr("172.16.0.1"),
				SubnetMask:     netip.MustParseAddr("255.255.255.255"),
				TotalIP:        1,
			},
		},
		{
			"192.168.1.1/32",
			"192.168.1.1/32",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("192.168.1.1"),
				BroadcastIP:    netip.MustParseAddr("192.168.1.1"),
				SubnetMask:     netip.MustParseAddr("255.255.255.255"),
				TotalIP:        1,
			},
		},
		{
			"8.8.8.8/32",
			"8.8.8.8/32",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("8.8.8.8"),
				BroadcastIP:    netip.MustParseAddr("8.8.8.8"),
				SubnetMask:     netip.MustParseAddr("255.255.255.255"),
				TotalIP:        1,
			},
		},
		{
			"1.1.1.1/32",
			"1.1.1.1/32",
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("1.1.1.1"),
				BroadcastIP:    netip.MustParseAddr("1.1.1.1"),
				SubnetMask:     netip.MustParseAddr("255.255.255.255"),
				TotalIP:        1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := netip.MustParsePrefix(tt.input)
			got, err := CalcSubnetInfo(input)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestCalcSubnetInfo_AnyAddressInSubnet tests that CalcSubnetInfo returns the same subnet info
// regardless of which IP address within a subnet is provided as input
func TestCalcSubnetInfo_AnyAddressInSubnet(t *testing.T) {
	tests := []struct {
		name           string
		inputs         []string // Different addresses in the same subnet
		expectedSubnet SubnetInfo
	}{
		{
			"/8 - 10.0.0.0/8 subnet",
			[]string{
				"10.0.0.0/8",
				"10.0.0.1/8",
				"10.128.64.32/8",
				"10.255.255.254/8",
				"10.255.255.255/8",
			},
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("10.255.255.255"),
				SubnetMask:     netip.MustParseAddr("255.0.0.0"),
				TotalIP:        16777216},
		},
		{
			"/16 - 192.168.0.0/16 subnet",
			[]string{
				"192.168.0.0/16",
				"192.168.0.1/16",
				"192.168.128.100/16",
				"192.168.255.254/16",
				"192.168.255.255/16",
			},
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("192.168.0.0"),
				BroadcastIP:    netip.MustParseAddr("192.168.255.255"),
				SubnetMask:     netip.MustParseAddr("255.255.0.0"),
				TotalIP:        65536},
		},
		{
			"/24 - 172.16.1.0/24 subnet",
			[]string{
				"172.16.1.0/24",
				"172.16.1.1/24",
				"172.16.1.128/24",
				"172.16.1.254/24",
				"172.16.1.255/24",
			},
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("172.16.1.0"),
				BroadcastIP:    netip.MustParseAddr("172.16.1.255"),
				SubnetMask:     netip.MustParseAddr("255.255.255.0"),
				TotalIP:        256},
		},
		{
			"/28 - 10.0.0.16/28 subnet",
			[]string{
				"10.0.0.16/28",
				"10.0.0.17/28",
				"10.0.0.24/28",
				"10.0.0.30/28",
				"10.0.0.31/28",
			},
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.0.0.16"),
				BroadcastIP:    netip.MustParseAddr("10.0.0.31"),
				SubnetMask:     netip.MustParseAddr("255.255.255.240"),
				TotalIP:        16},
		},
		{
			"/30 - 192.168.1.8/30 point-to-point subnet",
			[]string{
				"192.168.1.8/30",
				"192.168.1.9/30",
				"192.168.1.10/30",
				"192.168.1.11/30",
			},
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("192.168.1.8"),
				BroadcastIP:    netip.MustParseAddr("192.168.1.11"),
				SubnetMask:     netip.MustParseAddr("255.255.255.252"),
				TotalIP:        4},
		},
		{
			"/31 - 10.0.0.4/31 point-to-point subnet (RFC 3021)",
			[]string{
				"10.0.0.4/31",
				"10.0.0.5/31",
			},
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.0.0.4"),
				BroadcastIP:    netip.MustParseAddr("10.0.0.5"),
				SubnetMask:     netip.MustParseAddr("255.255.255.254"),
				TotalIP:        2},
		},
		{
			"/12 - 172.16.0.0/12 subnet",
			[]string{
				"172.16.0.0/12",
				"172.16.0.1/12",
				"172.24.128.64/12",
				"172.31.255.254/12",
				"172.31.255.255/12",
			},
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("172.16.0.0"),
				BroadcastIP:    netip.MustParseAddr("172.31.255.255"),
				SubnetMask:     netip.MustParseAddr("255.240.0.0"),
				TotalIP:        1048576},
		},
		{
			"/20 - 10.64.0.0/20 subnet",
			[]string{
				"10.64.0.0/20",
				"10.64.0.1/20",
				"10.64.8.128/20",
				"10.64.15.254/20",
				"10.64.15.255/20",
			},
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("10.64.0.0"),
				BroadcastIP:    netip.MustParseAddr("10.64.15.255"),
				SubnetMask:     netip.MustParseAddr("255.255.240.0"),
				TotalIP:        4096},
		},
		{
			"/4 - 16.0.0.0/4 subnet",
			[]string{
				"16.0.0.0/4",
				"16.0.0.1/4",
				"24.128.64.32/4",
				"31.255.255.254/4",
				"31.255.255.255/4",
			},
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("16.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("31.255.255.255"),
				SubnetMask:     netip.MustParseAddr("240.0.0.0"),
				TotalIP:        268435456},
		},
		{
			"/0 - default route",
			[]string{
				"0.0.0.0/0",
				"0.0.0.1/0",
				"128.128.128.128/0",
				"255.255.255.254/0",
				"255.255.255.255/0",
			},
			SubnetInfo{
				NetworkAddress: netip.MustParseAddr("0.0.0.0"),
				BroadcastIP:    netip.MustParseAddr("255.255.255.255"),
				SubnetMask:     netip.MustParseAddr("0.0.0.0"),
				TotalIP:        4294967296},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i, inputStr := range tt.inputs {
				input := netip.MustParsePrefix(inputStr)
				got, err := CalcSubnetInfo(input)
				assert.NoError(t, err, "input %d: %s", i, inputStr)
				assert.Equal(t, tt.expectedSubnet, got, "input %d: %s should produce same subnet info", i, inputStr)
			}
		})
	}
}
