package main

import (
	"errors"
	"io/ioutil"
	"net"
	"strconv"
	"sync"
	"time"
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

		iface.TXSpeed, iface.RXSpeed = nwl1.getNetworkRateForInf(inf.Name)
		iface.TXSpeedShortened = normalizeBValue(iface.TXSpeed) + "/s"
		iface.RXSpeedShortened = normalizeBValue(iface.RXSpeed) + "/s"
		
		val, err = ioutil.ReadFile(infpath + "speed")
		iface.MaxSpeedShortened = stringStripNewline(string(val)) + " Mbit"
		if (err != nil) {
			iface.MaxSpeedShortened = "Unknown"
		}

		interfaces = append(interfaces, iface)
	}

	return
}

func getTotalUpAndDownForInterfaceName(name string) (up uint64, down uint64, err error) {
	up = 0
	down = 0

	infpath := "/sys/class/net/" + name + "/"
	var val []byte

	val, err = ioutil.ReadFile(infpath + "statistics/tx_bytes")
	if (err != nil) {
		return up, down, errors.New("can't read tx bytes")
	}
	up, err = strconv.ParseUint(stringStripNewline(string(val)), 10, 64)
	if (err != nil) {
		return up, down, errors.New("can't convert tx bytes to a number")
	}

	val, err = ioutil.ReadFile(infpath + "statistics/rx_bytes")
	if (err != nil) {
		return up, down, errors.New("can't read rx bytes")
	}
	down, err = strconv.ParseUint(stringStripNewline(string(val)), 10, 64)
	if (err != nil) {
		return up, down, errors.New("can't convert rx bytes to a number")
	}

	return
}

type networkingLog struct {
	mu sync.Mutex
	maxlen int
	waittime time.Duration
	// History of up and down rate per interface
	up map[string][]uint64
	down map[string][]uint64

	// Last up or down value per interface
	uplast map[string]uint64
	downlast map[string]uint64
}
var nwl1 networkingLog

func (nwl *networkingLog) getNetworkLog() (up map[string][]uint64, down map[string][]uint64) {
	nwl.mu.Lock()
	defer nwl.mu.Unlock()
	return nwl.up, nwl.down
}

func (nwl *networkingLog) getNetworkRateForInf(inf string) (up uint64, down uint64) {
	nwl.mu.Lock()
	up = 0
	down = 0
	if val, ok := nwl.up[inf]; ok {
		if (len(val) > 0) {
			up = val[len(val)-1]
		}
	}
	if val, ok := nwl.down[inf]; ok {
		if (len(val) > 0) {
			down = val[len(val)-1]
		}
	}
	nwl.mu.Unlock()
	return
}

func (nwl *networkingLog) getNetworkLogMaxLength() (max int) {
	return nwl.maxlen
}

func (nwl *networkingLog) networkLoggingTask() {
	nwl.up = make(map[string][]uint64)
	nwl.down = make(map[string][]uint64)
	nwl.uplast = make(map[string]uint64)
	nwl.downlast = make(map[string]uint64)
	infs, _ := net.Interfaces()
	for _, inf := range infs {
		up, down, _ := getTotalUpAndDownForInterfaceName(inf.Name)
		nwl.uplast[inf.Name] = up
		nwl.downlast[inf.Name] = down
	}
	for {
		infs, err := net.Interfaces()
		if err != nil {
			continue
		}
		for _, inf := range infs {
			up, down, err := getTotalUpAndDownForInterfaceName(inf.Name)
			if (err == nil){
				nwl.mu.Lock()

				nwl.up[inf.Name] = append(nwl.up[inf.Name], up - nwl.uplast[inf.Name])
				if (len(nwl.up[inf.Name]) > 0) {
					for l := range nwl.up {
						if (len(nwl.up[l]) > nwl.maxlen) {
							nwl.up[l] = nwl.up[l][1:]
						}
					}
				}

				nwl.down[inf.Name] = append(nwl.down[inf.Name], down - nwl.downlast[inf.Name])
				if (len(nwl.down[inf.Name]) > 0) {
					for l := range nwl.down {
						if (len(nwl.down[l]) > nwl.maxlen) {
							nwl.down[l] = nwl.down[l][1:]
						}
					}
				}

				nwl.mu.Unlock()
			}
			nwl.uplast[inf.Name] = up
			nwl.downlast[inf.Name] = down
		}
		time.Sleep(nwl.waittime)
	}
}