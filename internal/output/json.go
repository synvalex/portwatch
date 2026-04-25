package output

import (
	"encoding/json"
	"io"

	"github.com/user/portwatch/internal/ports"
)

type jsonListener struct {
	Proto   string      `json:"proto"`
	Address string      `json:"address"`
	Port    uint16      `json:"port"`
	Process *jsonProcess `json:"process,omitempty"`
}

type jsonProcess struct {
	PID  int    `json:"pid"`
	Name string `json:"name,omitempty"`
	Exe  string `json:"exe,omitempty"`
}

func writeJSON(w io.Writer, listeners []ports.Listener) error {
	out := make([]jsonListener, 0, len(listeners))
	for _, l := range listeners {
		jl := jsonListener{
			Proto:   l.Proto,
			Address: l.Addr.String(),
			Port:    l.Port,
		}
		if l.Process != nil {
			jl.Process = &jsonProcess{
				PID:  l.Process.PID,
				Name: l.Process.Name,
				Exe:  l.Process.Exe,
			}
		}
		out = append(out, jl)
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
