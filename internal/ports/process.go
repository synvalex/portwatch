package ports

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// ProcessInfo holds metadata about the process owning a socket.
type ProcessInfo struct {
	PID  int
	Name string
	Exe  string
}

func (p ProcessInfo) String() string {
	if p.Name != "" {
		return fmt.Sprintf("%s(%d)", p.Name, p.PID)
	}
	return fmt.Sprintf("pid(%d)", p.PID)
}

// LookupInode searches /proc for the process that owns the given socket inode.
func LookupInode(inode uint64) (*ProcessInfo, error) {
	target := fmt.Sprintf("socket:[%d]", inode)

	entries, err := filepath.Glob("/proc/[0-9]*/fd/*")
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		link, err := os.Readlink(entry)
		if err != nil {
			continue
		}
		if link != target {
			continue
		}

		// Extract PID from path /proc/<pid>/fd/<n>
		parts := strings.Split(entry, string(os.PathSeparator))
		if len(parts) < 4 {
			continue
		}
		pid, err := strconv.Atoi(parts[2])
		if err != nil {
			continue
		}

		name := readProcName(pid)
		exe := readProcExe(pid)
		return &ProcessInfo{PID: pid, Name: name, Exe: exe}, nil
	}
	return nil, nil
}

func readProcName(pid int) string {
	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/comm", pid))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

func readProcExe(pid int) string {
	link, err := os.Readlink(fmt.Sprintf("/proc/%d/exe", pid))
	if err != nil {
		return ""
	}
	return link
}
