package alert

import (
	"fmt"
	"log"

	"github.com/user/portwatch/internal/ports"
	"github.com/user/portwatch/internal/rules"
)

// Dispatcher evaluates listeners against rules and dispatches alerts.
type Dispatcher struct {
	Rules     []rules.Rule
	Notifiers []Notifier
}

// NewDispatcher creates a Dispatcher with the given rules and notifiers.
func NewDispatcher(r []rules.Rule, notifiers ...Notifier) *Dispatcher {
	return &Dispatcher{
		Rules:     r,
		Notifiers: notifiers,
	}
}

// Evaluate checks each listener against configured rules and fires alerts
// for any listener that matches an ALERT-level rule or has no matching rule.
func (d *Dispatcher) Evaluate(listeners []ports.Listener) {
	for _, l := range listeners {
		d.evaluateOne(l)
	}
}

func (d *Dispatcher) evaluateOne(l ports.Listener) {
	for _, r := range d.Rules {
		if r.Matches(l) {
			if r.Action == rules.ActionAllow {
				return // explicitly allowed, no alert
			}
			// ActionDeny or ActionAlert
			a := New(LevelAlert, l, fmt.Sprintf("denied by rule %q", r.Name))
			d.dispatch(a)
			return
		}
	}
	// No rule matched — unexpected listener
	a := New(LevelWarn, l, "unexpected listener: no matching rule")
	d.dispatch(a)
}

func (d *Dispatcher) dispatch(a Alert) {
	for _, n := range d.Notifiers {
		if err := n.Notify(a); err != nil {
			log.Printf("portwatch: alert dispatch error: %v", err)
		}
	}
}
