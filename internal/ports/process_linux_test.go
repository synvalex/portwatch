package ports

import (
	"fmt"
	"os"
	"testing"
)

func TestProcessInfo_String_WithName(t *testing.T) {
	p := ProcessInfo{PID: 42, Name: "nginx"}
	got := p.String()
	want := "nginx(42)"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestProcessInfo_String_NoName(t *testing.T) {
	p := ProcessInfo{PID: 99}
	got := p.String()
	want := "pid(99)"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestReadProcName_Self(t *testing.T) {
	pid := os.Getpid()
	name := readProcName(pid)
	if name == "" {
		t.Skip("cannot read /proc/self/comm (not Linux or no permission)")
	}
	if name == "" {
		t.Error("expected non-empty process name for self")
	}
}

func TestReadProcExe_Self(t *testing.T) {
	pid := os.Getpid()
	exe := readProcExe(pid)
	if exe == "" {
		t.Skip("cannot read /proc/self/exe")
	}
}

func TestLookupInode_NotFound(t *testing.T) {
	// inode 0 should never match a real socket
	info, err := LookupInode(0)
	if err != nil {
		// On non-Linux systems the glob may fail — skip gracefully
		t.Skipf("LookupInode not supported: %v", err)
	}
	if info != nil {
		t.Errorf("expected nil info for inode 0, got %v", info)
	}
}

func TestProcessInfo_String_Format(t *testing.T) {
	cases := []struct {
		info ProcessInfo
		want string
	}{
		{ProcessInfo{PID: 1, Name: "init"}, "init(1)"},
		{ProcessInfo{PID: 1234, Name: ""}, "pid(1234)"},
		{ProcessInfo{PID: 7, Name: "sshd", Exe: "/usr/sbin/sshd"}, "sshd(7)"},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf("pid=%d", tc.info.PID), func(t *testing.T) {
			if got := tc.info.String(); got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}
