package subnetcalc

import (
	"encoding/binary"
	"errors"
	"math"
	"net/netip"
)

type SubnetInfo struct {
	NetworkAddress netip.Addr
	BroadcastIP    netip.Addr
	SubnetMask     netip.Addr
	TotalIP        uint
}

type masks struct {
	SubnetMask   uint32
	WildcardMask uint32
}

// calcNetworkAddress parses the network address from a prefix.
func calcNetworkAddress(prefix netip.Prefix) netip.Addr {
	return prefix.Masked().Addr()
}

// calcMasks parses the subnet mask and wildcard mask from a prefix.
func calcMasks(prefix netip.Prefix) masks {
	networkBits := prefix.Bits()
	subnetMask := uint32(0xffffffff << (32 - networkBits))
	wildcardMask := ^subnetMask
	return masks{SubnetMask: subnetMask, WildcardMask: wildcardMask}
}

// calcBroadcastIPAddress parses the broadcast IP address from a network address and wildcard mask.
func calcBroadcastIPAddress(networkAddress netip.Addr, wildcardMask uint32) netip.Addr {
	networkAddressUint32 := binary.BigEndian.Uint32(networkAddress.AsSlice())
	lastIPUint32 := networkAddressUint32 | wildcardMask
	lastIPBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lastIPBytes, lastIPUint32)
	lastIPBytesArray := [4]byte{lastIPBytes[0], lastIPBytes[1], lastIPBytes[2], lastIPBytes[3]}
	return netip.AddrFrom4(lastIPBytesArray)
}

// calcTotalIP calculates the total number of IPs in a subnet.
func calcTotalIP(prefix netip.Prefix) uint {
	return uint(math.Pow(2, float64(32-prefix.Bits())))
}

// getSingleIPSubnetInfo returns subnet info for a single IP.
func getSingleIPSubnetInfo(ip netip.Addr) SubnetInfo {
	return SubnetInfo{
		NetworkAddress: ip,
		BroadcastIP:    ip,
		SubnetMask:     netip.AddrFrom4([4]byte{255, 255, 255, 255}),
		TotalIP:        1,
	}
}

// convSubnetMaskToIPAddr converts a subnet mask to an IP address.
func convSubnetMaskToIPAddr(subnetMask uint32) netip.Addr {
	subnetMaskBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(subnetMaskBytes, subnetMask)
	subnetMaskBytesArray := [4]byte{subnetMaskBytes[0], subnetMaskBytes[1], subnetMaskBytes[2], subnetMaskBytes[3]}
	return netip.AddrFrom4(subnetMaskBytesArray)
}

// CalcSubnetInfo calculates subnet info for a given prefix.
func CalcSubnetInfo(prefix netip.Prefix) (SubnetInfo, error) {
	if !prefix.IsValid() {
		return SubnetInfo{}, errors.New("invalid prefix")
	}

	if prefix.Addr().Is6() {
		return SubnetInfo{}, errors.New("IPv6 not supported yet")
	}

	if prefix.IsSingleIP() {
		return getSingleIPSubnetInfo(prefix.Addr()), nil
	}

	masks := calcMasks(prefix)
	networkAddress := calcNetworkAddress(prefix)
	subnetMask := convSubnetMaskToIPAddr(masks.SubnetMask)
	broadcastIP := calcBroadcastIPAddress(networkAddress, masks.WildcardMask)
	totalIP := calcTotalIP(prefix)

	return SubnetInfo{
		NetworkAddress: networkAddress,
		BroadcastIP:    broadcastIP,
		SubnetMask:     subnetMask,
		TotalIP:        totalIP,
	}, nil
}
