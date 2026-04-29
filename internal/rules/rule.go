package rules

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// Action defines what to do when a rule matches.
type Action string

const (
	ActionAllow Action = "allow"
	ActionAlert Action = "alert"
	ActionDeny  Action = "deny"
)

// Rule represents a single port monitoring rule.
type Rule struct {
	Name     string `yaml:"name"`
	Port     uint16 `yaml:"port"`
	Protocol string `yaml:"protocol"` // tcp, udp, or "" for any
	Address  string `yaml:"address"`  // CIDR or IP, empty means any
	Action   Action `yaml:"action"`
}

// Matches returns true if the rule applies to the given port, protocol, and address.
func (r *Rule) Matches(port uint16, protocol, address string) bool {
	if r.Port != 0 && r.Port != port {
		return false
	}
	if r.Protocol != "" && !strings.EqualFold(r.Protocol, protocol) {
		return false
	}
	if r.Address != "" {
		if !r.matchesAddress(address) {
			return false
		}
	}
	return true
}

func (r *Rule) matchesAddress(address string) bool {
	ip := net.ParseIP(address)
	if ip == nil {
		return false
	}
	// Try CIDR first
	if strings.Contains(r.Address, "/") {
		_, network, err := net.ParseCIDR(r.Address)
		if err != nil {
			return false
		}
		return network.Contains(ip)
	}
	// Exact IP match
	ruleIP := net.ParseIP(r.Address)
	return ruleIP != nil && ruleIP.Equal(ip)
}

// Validate checks that the rule is well-formed.
func (r *Rule) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("rule name must not be empty")
	}
	switch r.Action {
	case ActionAllow, ActionAlert, ActionDeny:
	default:
		return fmt.Errorf("rule %q: unknown action %q", r.Name, r.Action)
	}
	if r.Protocol != "" {
		p := strings.ToLower(r.Protocol)
		if p != "tcp" && p != "udp" {
			return fmt.Errorf("rule %q: protocol must be tcp, udp, or empty", r.Name)
		}
	}
	if r.Address != "" {
		if strings.Contains(r.Address, "/") {
			if _, _, err := net.ParseCIDR(r.Address); err != nil {
				return fmt.Errorf("rule %q: invalid CIDR address %q: %w", r.Name, r.Address, err)
			}
		} else if net.ParseIP(r.Address) == nil {
			return fmt.Errorf("rule %q: invalid address %q", r.Name, r.Address)
		}
	}
	_ = strconv.Itoa(int(r.Port)) // port is uint16, always valid
	return nil
}
