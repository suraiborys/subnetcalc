// Package subnetcalc calculates IPv4 subnet information from CIDR notation.
// It computes network addresses, broadcast IPs, subnet masks, and address counts
// for any valid IPv4 prefix.
package subnetcalc

import (
	"encoding/binary"
	"errors"
	"math"
	"net/netip"
)

// SubnetInfo represents calculated information about an IPv4 subnet.
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

func calcNetworkAddress(prefix netip.Prefix) netip.Addr {
	return prefix.Masked().Addr()
}

func calcMasks(prefix netip.Prefix) masks {
	networkBits := prefix.Bits()
	subnetMask := uint32(0xffffffff << (32 - networkBits))
	wildcardMask := ^subnetMask
	return masks{SubnetMask: subnetMask, WildcardMask: wildcardMask}
}

func calcBroadcastIPAddress(networkAddress netip.Addr, wildcardMask uint32) netip.Addr {
	networkAddressUint32 := binary.BigEndian.Uint32(networkAddress.AsSlice())
	lastIPUint32 := networkAddressUint32 | wildcardMask
	lastIPBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lastIPBytes, lastIPUint32)
	lastIPBytesArray := [4]byte{lastIPBytes[0], lastIPBytes[1], lastIPBytes[2], lastIPBytes[3]}
	return netip.AddrFrom4(lastIPBytesArray)
}

func calcTotalIP(prefix netip.Prefix) uint {
	return uint(math.Pow(2, float64(32-prefix.Bits())))
}

func getSingleIPSubnetInfo(ip netip.Addr) SubnetInfo {
	return SubnetInfo{
		NetworkAddress: ip,
		BroadcastIP:    ip,
		SubnetMask:     netip.AddrFrom4([4]byte{255, 255, 255, 255}),
		TotalIP:        1,
	}
}

func uint32ToAddr(ipInt uint32) netip.Addr {
	subnetMaskBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(subnetMaskBytes, ipInt)
	subnetMaskBytesArray := [4]byte{subnetMaskBytes[0], subnetMaskBytes[1], subnetMaskBytes[2], subnetMaskBytes[3]}
	return netip.AddrFrom4(subnetMaskBytesArray)
}

// CalcSubnetInfo calculates subnet information for the given IPv4 prefix.
// It returns the network address, broadcast IP, subnet mask, and total IP count.
//
// Special cases:
//   - /32 prefix: single host, NetworkAddress equals BroadcastIP
//   - IPv6 prefix: returns error (not yet supported)
//   - Invalid prefix: returns error
//
// Example:
//
//	prefix := netip.MustParsePrefix("192.168.1.0/24")
//	info, err := CalcSubnetInfo(prefix)
//	if err != nil {
//	    return err
//	}
//	fmt.Printf("Network: %s\n", info.NetworkAddress)
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
	subnetMask := uint32ToAddr(masks.SubnetMask)
	broadcastIP := calcBroadcastIPAddress(networkAddress, masks.WildcardMask)
	totalIP := calcTotalIP(prefix)

	return SubnetInfo{
		NetworkAddress: networkAddress,
		BroadcastIP:    broadcastIP,
		SubnetMask:     subnetMask,
		TotalIP:        totalIP,
	}, nil
}
