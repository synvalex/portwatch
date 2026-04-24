//go:build linux

package ports

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ProcScanner reads open TCP/UDP listeners from /proc/net.
type ProcScanner struct{}

// NewScanner returns a platform-specific Scanner implementation.
func NewScanner() Scanner {
	return &ProcScanner{}
}

// Scan reads /proc/net/tcp and /proc/net/udp to find listening ports.
func (s *ProcScanner) Scan() ([]Listener, error) {
	var listeners []Listener

	for _, entry := range []struct {
		file  string
		proto string
	}{
		{"/proc/net/tcp", "tcp"},
		{"/proc/net/tcp6", "tcp6"},
		{"/proc/net/udp", "udp"},
		{"/proc/net/udp6", "udp6"},
	} {
		ls, err := parseProcNet(entry.file, entry.proto)
		if err != nil {
			continue // file may not exist on all kernels
		}
		listeners = append(listeners, ls...)
	}

	return DeduplicateListeners(listeners), nil
}

// parseProcNet parses a single /proc/net/{tcp,udp} file.
func parseProcNet(path, proto string) ([]Listener, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	var listeners []Listener
	scanner := bufio.NewScanner(f)
	scanner.Scan() // skip header line

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 4 {
			continue
		}
		// state 0A = TCP_LISTEN, for UDP we include all (stateless)
		state := fields[3]
		if strings.HasPrefix(proto, "tcp") && state != "0A" {
			continue
		}

		addr, port, err := parseHexAddr(fields[1])
		if err != nil {
			continue
		}

		listeners = append(listeners, Listener{
			Protocol: proto,
			Address:  addr,
			Port:     port,
		})
	}
	return listeners, scanner.Err()
}

// parseHexAddr converts a /proc/net hex-encoded "addr:port" to dotted-decimal and int.
func parseHexAddr(hexAddr string) (string, int, error) {
	parts := strings.Split(hexAddr, ":")
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("unexpected format: %s", hexAddr)
	}
	portVal, err := strconv.ParseInt(parts[1], 16, 32)
	if err != nil {
		return "", 0, err
	}
	addrVal, err := strconv.ParseUint(parts[0], 16, 32)
	if err != nil {
		return "", 0, err
	}
	ip := fmt.Sprintf("%d.%d.%d.%d",
		byte(addrVal), byte(addrVal>>8), byte(addrVal>>16), byte(addrVal>>24))
	return ip, int(portVal), nil
}
