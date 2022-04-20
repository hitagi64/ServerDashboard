package main

import (
	"io/ioutil"
	"net"
)

type interfaceInfo struct {
	Name string
	Addr4 string
	Addr6 string
	Duplex string
	Operstate string
	TXSpeed uint64
	RXSpeed uint64
	TXSpeedShortened string
	RXSpeedShortened string
	MaxSpeedShortened string
}

func getIpv46Ip(addrs []net.Addr) (ipv4Addr net.IP, ipv6Addr net.IP) {
	for _, addr := range addrs { // get ipv4 address
        if ipv4Addr = addr.(*net.IPNet).IP.To4(); ipv4Addr != nil {
            break
        }
    }
	for _, addr := range addrs { // get ipv6 address
		ipv6Addr = addr.(*net.IPNet).IP.To16()
		ipv4Addr2 := addr.(*net.IPNet).IP.To4()
        if ipv4Addr2 == nil {
            break
        }
    }
	return
}

func getNetworkInterfaceInfo() (interfaces []interfaceInfo, err error) {
	infs, err := net.Interfaces()
	if err != nil {
		return interfaces, err
	}
	for _, inf := range infs {
		var iface interfaceInfo
		var val []byte
		infpath := "/sys/class/net/" + inf.Name + "/"

		iface.Name = inf.Name

		ipslist, err := inf.Addrs()
		if err == nil {
			ip4, ip6 := getIpv46Ip(ipslist)
			if (ip4 != nil) {
				iface.Addr4 = ip4.String()
			}
			if (ip6 != nil) {
				iface.Addr6 = ip6.String()
			}
		}

		val, err = ioutil.ReadFile(infpath + "duplex")
		iface.Duplex = stringStripNewline(string(val))
		if (err != nil) {
			iface.Duplex = "Unknown"
		}

		val, err = ioutil.ReadFile(infpath + "operstate")
		iface.Operstate = stringStripNewline(string(val))
		if (err != nil) {
			iface.Operstate = "Unknown"
		}

		iface.TXSpeed = 0
		iface.RXSpeed = 0
		iface.TXSpeedShortened = "1G"
		iface.RXSpeedShortened = "1G"
		
		val, err = ioutil.ReadFile(infpath + "speed")
		iface.MaxSpeedShortened = stringStripNewline(string(val)) + " Mbit"
		if (err != nil) {
			iface.MaxSpeedShortened = "Unknown"
		}

		interfaces = append(interfaces, iface)
	}

	return
}