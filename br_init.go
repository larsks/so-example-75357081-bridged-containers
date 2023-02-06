package main

import (
	"flag"
	"log"
	"net"

	"github.com/vishvananda/netlink"
)

var (
	bridgeIp   string
	bridgeName string
	bridgeMac  string
	setHwAddr  bool
)

func init() {
	flag.StringVar(&bridgeIp, "ip", "10.10.10.1/24", "Set bridge ip address")
	flag.StringVar(&bridgeName, "name", "br1", "Set bridge name")
	flag.StringVar(&bridgeMac, "mac", "", "Set bridge MAC address")
	flag.BoolVar(&setHwAddr, "sethwaddr", true, "Set to false to skip setting bridge MAC")
}

func main() {
	flag.Parse()

	br := &netlink.Bridge{
		LinkAttrs: netlink.LinkAttrs{
			Name:   bridgeName,
			MTU:    1500,
			TxQLen: -1,
		},
	}

	log.Printf("add bridge device %s", bridgeName)
	if err := netlink.LinkAdd(br); err != nil {
		panic(err)
	}

	bridgeIf, err := netlink.LinkByName(bridgeName)
	if err != nil {
		panic(err)
	}

	log.Printf("set bridge link up")
	err = netlink.LinkSetUp(bridgeIf)
	if err != nil {
		panic(err)
	}

	log.Printf("set bridge promiscuous on")
	if err = netlink.SetPromiscOn(bridgeIf); err != nil {
		panic(err)
	}

	bridgeIp, bridgeNet, err := net.ParseCIDR(bridgeIp)
	if err != nil {
		panic(err)
	}
	bridgeNet.IP = bridgeIp

	log.Printf("set bridge ip addr %s", bridgeIp)
	addr := &netlink.Addr{IPNet: bridgeNet}
	if err = netlink.AddrAdd(bridgeIf, addr); err != nil {
		panic(err)
	}

	if setHwAddr {
		var hwaddr net.HardwareAddr

		if bridgeMac != "" {
			hwaddr, err = net.ParseMAC(bridgeMac)
			if err != nil {
				panic(err)
			}
		} else {
			hwaddr = bridgeIf.Attrs().HardwareAddr
		}

		log.Printf("set bridge hardware address %v", hwaddr)

		if err = netlink.LinkSetHardwareAddr(bridgeIf, hwaddr); err != nil {
			panic(err)
		}
	}

	log.Printf("bridge created successfully")
}
