package socket

import (
	"fmt"
	"net"
	"syscall"
	"time"
)

const (
	defaultName = "tap_netstack"
	defaultCIDR = "10.58.0.1/24"
)

type Interface struct {
	fd      int
	name    string
	cidr    string
	timeout time.Duration
}

func New(opts ...InterfaceOption) (*Interface, error) {
	i := &Interface{
		name: defaultName,
		cidr: defaultCIDR,
	}

	for _, opt := range opts {
		opt(i)
	}

	if i.name == defaultName {
		err := createInterface(i.cidr)
		if err != nil {
			return nil, err
		}
	}

	fd, err := open(i.name)
	if err != nil {
		return nil, err
	}

	i.fd = fd

	err = i.setTimeout()
	if err != nil {
		return nil, err
	}

	return i, nil
}

func (i *Interface) Read(buf []byte) (int, error) {
	return syscall.Read(i.fd, buf)
}

func (i *Interface) Write(buf []byte) (int, error) {
	return syscall.Write(i.fd, buf)
}

func (i *Interface) Close() error {
	if i.name == defaultName {
		err := stopInterface()
		if err != nil {
			return fmt.Errorf("failed to stop default interface: %w", err)
		}
	}

	return syscall.Close(i.fd)
}

func (i *Interface) GetHardwareAddr() net.HardwareAddr {
	iface, err := net.InterfaceByName(i.name)
	if err != nil {
		return net.HardwareAddr{}
	}

	return iface.HardwareAddr
}

func (i *Interface) setTimeout() error {
	if i.timeout != 0 {
		t := &syscall.Timeval{
			Sec:  int64(i.timeout.Seconds()),
			Usec: 0,
		}

		err := syscall.SetsockoptTimeval(i.fd, syscall.SOL_SOCKET, syscall.SO_RCVTIMEO, t)
		if err != nil {
			return fmt.Errorf("failed to set timeout: %w", err)
		}
	}

	return nil
}
