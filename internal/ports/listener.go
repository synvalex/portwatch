package ports

// Listener represents a single network listener discovered on the host.
// It is the central data structure passed through the pipeline.
type Listener struct {
	// Address holds the parsed IP and port.
	Address ParsedAddress

	// Protocol is "tcp", "tcp6", "udp", or "udp6".
	Protocol string

	// Inode is the socket inode number from /proc/net.
	Inode uint64

	// Process holds enriched process information (may be nil).
	Process *ProcessInfo

	// Hostname is the result of a reverse-DNS lookup (may be empty).
	Hostname string

	// Geo holds GeoIP enrichment data (may be nil).
	Geo *GeoInfo

	// Fingerprint holds binary-fingerprint data (may be nil).
	Fingerprint *Fingerprint
}

// String returns a human-readable representation of the listener.
func (l Listener) String() string {
	proto := l.Protocol
	if proto == "" {
		proto = "unknown"
	}
	addr := l.Address.String()
	if l.Hostname != "" {
		addr = l.Hostname + " (" + addr + ")"
	}
	return proto + " " + addr
}
