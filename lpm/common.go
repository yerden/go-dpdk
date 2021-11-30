package lpm

import (
	"encoding/binary"
	"net"
)

func cvtIPv4(addr net.IP) (ip uint32) {
	if addr = addr.To4(); addr == nil {
		panic("not an IPv4 address")
	}
	return binary.BigEndian.Uint32(addr[:])
}

func cvtIPv4Net(ipp net.IPNet) (ip uint32, prefix uint8) {
	addr := ipp.IP.Mask(ipp.Mask)
	ones, _ := ipp.Mask.Size()
	return cvtIPv4(addr), uint8(ones)
}

func cvtIPv6Net(ipp net.IPNet) (ip [16]byte, prefix uint8) {
	addr := ipp.IP.Mask(ipp.Mask)
	ones, _ := ipp.Mask.Size()
	if len(addr) <= 4 {
		panic("not an IPv6 address")
	}
	copy(ip[:], addr.To16())
	return ip, uint8(ones)
}
