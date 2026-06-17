package process

import (
	"fmt"
	"net"
)

func CheckPort(port int) error {
	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return fmt.Errorf("port %d is already in use", port)
	}
	ln.Close()
	return nil
}

func CheckPorts(ports []int) error {
	for _, p := range ports {
		if err := CheckPort(p); err != nil {
			return err
		}
	}
	return nil
}
