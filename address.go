package fins

import "net"

// DeviceAddress A FINS device address
type DeviceAddress struct {
	network byte
	node    byte
	unit    byte
}

// UDPAddress A full device address
type UDPAddress struct {
	deviceAddress DeviceAddress
	udpAddress    *net.UDPAddr
}

func NewUDPAddress(ip string, port int, network, node, unit byte) UDPAddress {
	return UDPAddress{
		udpAddress: &net.UDPAddr{
			IP:   net.ParseIP(ip),
			Port: port,
		},
		deviceAddress: DeviceAddress{
			network: network,
			node:    node,
			unit:    unit,
		},
	}
}
