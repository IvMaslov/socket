package socket

import (
	"fmt"
	"log"
	"net"
	"os/exec"
	"syscall"
)

// htons converts a short (16-bit) integer from host byte order to network byte order.
func htons(i uint16) uint16 {
	return (i<<8)&0xff00 | i>>8
}

// getInterfaceIndex returns the index of the network interface with the given name.
func getInterfaceIndex(iface string) int {
	ifaceObj, err := net.InterfaceByName(iface)
	if err != nil {
		log.Fatalf("Failed to get interface index for %s: %v", iface, err)
	}

	return ifaceObj.Index
}

func createInterface(cidr string) error {
	_, err := exec.Command("/sbin/ip", "tuntap", "add", "dev", defaultName, "mode", "tap").Output()
	if err != nil {
		return fmt.Errorf("failed to create device: %w", err)
	}

	_, err = exec.Command("/sbin/ip", "addr", "add", cidr, "dev", defaultName).Output()
	if err != nil {
		return fmt.Errorf("failed to add address to device: %w", err)
	}

	_, err = exec.Command("/sbin/ip", "link", "set", defaultName, "up").Output()
	if err != nil {
		return fmt.Errorf("failed to up device: %w", err)
	}

	return nil
}

func stopInterface() error {
	_, err := exec.Command("/sbin/ip", "tuntap", "del", "dev", defaultName, "mode", "tap").Output()
	if err != nil {
		return fmt.Errorf("failed to stop device: %w", err)
	}

	return nil
}

func open(name string) (int, error) {
	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, int(htons(syscall.ETH_P_ALL)))
	if err != nil {
		return -1, fmt.Errorf("failed to open raw socket: %w", err)
	}

	ifaceName := [16]byte{}
	copy(ifaceName[:], name)

	sll := syscall.SockaddrLinklayer{
		Protocol: htons(syscall.ETH_P_ALL),
		Ifindex:  getInterfaceIndex(name),
	}

	if err := syscall.Bind(fd, &sll); err != nil {
		return -1, fmt.Errorf("failed to bind socket to interface %s: %w", name, err)
	}

	return fd, nil
}
