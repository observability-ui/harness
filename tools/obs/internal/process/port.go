package process

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func CheckPort(port int) error {
	addr := fmt.Sprintf(":%d", port)
	ln4, err4 := net.Listen("tcp4", addr)
	if err4 != nil {
		return fmt.Errorf("port %d is already in use", port)
	}
	ln4.Close()
	ln6, err6 := net.Listen("tcp6", addr)
	if err6 != nil {
		if !isAddrFamilyUnsupported(err6) {
			return fmt.Errorf("port %d is already in use", port)
		}
	} else {
		ln6.Close()
	}
	return nil
}

func isAddrFamilyUnsupported(err error) bool {
	var opErr *net.OpError
	if errors.As(err, &opErr) {
		var sysErr *os.SyscallError
		if errors.As(opErr.Err, &sysErr) {
			return errors.Is(sysErr.Err, syscall.EAFNOSUPPORT) || errors.Is(sysErr.Err, syscall.ENOPROTOOPT)
		}
	}
	return false
}

func FreePorts(ports []int) error {
	for _, p := range ports {
		if err := CheckPort(p); err != nil {
			if freeErr := freePort(p); freeErr != nil {
				return fmt.Errorf("port %d in use and could not free it: %w", p, freeErr)
			}
		}
	}
	return nil
}

func freePort(port int) error {
	pids, err := findPIDsOnPort(port)
	if err != nil {
		return err
	}
	for _, pid := range pids {
		// SIGTERM the process holding the port
		syscall.Kill(pid, syscall.SIGTERM)
	}
	// Wait for port to become available
	for range 20 {
		time.Sleep(250 * time.Millisecond)
		if err := CheckPort(port); err == nil {
			return nil
		}
	}
	// Force kill if still occupied
	for _, pid := range pids {
		if IsAlive(pid) {
			syscall.Kill(pid, syscall.SIGKILL)
		}
	}
	time.Sleep(500 * time.Millisecond)
	if err := CheckPort(port); err != nil {
		return fmt.Errorf("port %d still in use after killing processes", port)
	}
	return nil
}

func findPIDsOnPort(port int) ([]int, error) {
	out, err := exec.Command("lsof", "-ti", fmt.Sprintf(":%d", port)).Output()
	if err != nil {
		return nil, fmt.Errorf("no process found on port %d", port)
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var pids []int
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		pid, err := strconv.Atoi(line)
		if err != nil {
			continue
		}
		pids = append(pids, pid)
	}
	if len(pids) == 0 {
		return nil, fmt.Errorf("no process found on port %d", port)
	}
	return pids, nil
}

func ProbePort(port int) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", port), 500*time.Millisecond)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func ProbePorts(ports []int) bool {
	for _, p := range ports {
		if !ProbePort(p) {
			return false
		}
	}
	return len(ports) > 0
}

func IsAlive(pid int) bool {
	return syscall.Kill(pid, 0) == nil
}

func KillGroup(pid int) {
	syscall.Kill(-pid, syscall.SIGINT)
}
