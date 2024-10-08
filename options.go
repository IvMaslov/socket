package socket

import (
	"time"
)

type InterfaceOption func(*Interface)

// Option for network device, `tap_netstack` by deafult
func WithDevice(dev string) InterfaceOption {
	return func(i *Interface) {
		i.name = dev
	}
}

// Option for IP address with mask, `10.58.0.1/24` by default
func WithCIDR(cidr string) InterfaceOption {
	return func(i *Interface) {
		i.cidr = cidr
	}
}

// Option to set timeout for read and write for interface
func WithTimeout(t time.Duration) InterfaceOption {
	return func(i *Interface) {
		i.timeout = t
	}
}

// Option for setting default gateway
func WithDefaultGateway() InterfaceOption {
	return func(i *Interface) {
		i.name = getDefaultGateway()
	}
}
