package output

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/user/portwatch/internal/ports"
)

const (
	TableFormatText = "text"
	TableFormatJSON = "json"
)

// TableWriter renders a list of listeners as a formatted table.
type TableWriter struct {
	w      io.Writer
	format string
}

// NewTableWriter creates a TableWriter that writes to w in the given format.
func NewTableWriter(w io.Writer, format string) *TableWriter {
	return &TableWriter{w: w, format: strings.ToLower(format)}
}

// Write renders listeners to the underlying writer.
func (t *TableWriter) Write(listeners []ports.Listener) error {
	switch t.format {
	case TableFormatJSON:
		return writeJSON(t.w, listeners)
	default:
		return writeText(t.w, listeners)
	}
}

func writeText(w io.Writer, listeners []ports.Listener) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "PROTO\tADDRESS\tPORT\tPID\tPROCESS")
	for _, l := range listeners {
		pid := "-"
		name := "-"
		if l.Process != nil {
			pid = fmt.Sprintf("%d", l.Process.PID)
			if l.Process.Name != "" {
				name = l.Process.Name
			}
		}
		fmt.Fprintf(tw, "%s\t%s\t%d\t%s\t%s\n",
			l.Proto, l.Addr.String(), l.Port, pid, name)
	}
	return tw.Flush()
}
