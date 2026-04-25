package output

import (
	"fmt"
	"io"
	"os"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/ports"
)

// Renderer writes a snapshot of listeners to an output destination.
type Renderer struct {
	w      io.Writer
	writer *TableWriter
}

// NewRenderer creates a Renderer configured from cfg, writing to stdout.
func NewRenderer(cfg config.OutputConfig) (*Renderer, error) {
	return NewRendererWithWriter(os.Stdout, cfg)
}

// NewRendererWithWriter creates a Renderer writing to the given writer.
func NewRendererWithWriter(w io.Writer, cfg config.OutputConfig) (*Renderer, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("output config: %w", err)
	}
	return &Renderer{
		w:      w,
		writer: NewTableWriter(w, cfg.Format),
	}, nil
}

// Render writes the listeners to the configured output.
func (r *Renderer) Render(listeners []ports.Listener) error {
	return r.writer.Write(listeners)
}
