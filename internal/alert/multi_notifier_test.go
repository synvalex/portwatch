package alert

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// recordingNotifier captures every event it receives.
type recordingNotifier struct {
	events []Event
	err    error
}

func (r *recordingNotifier) Notify(e Event) error {
	r.events = append(r.events, e)
	return r.err
}

func TestMultiNotifier_Empty(t *testing.T) {
	mn := NewMultiNotifier()
	err := mn.Notify(Event{})
	require.NoError(t, err)
	assert.Equal(t, 0, mn.Len())
}

func TestMultiNotifier_AllReceiveEvent(t *testing.T) {
	a := &recordingNotifier{}
	b := &recordingNotifier{}
	mn := NewMultiNotifier(a, b)

	ev := Event{Type: EventAppeared, Listener: makeListener("tcp", "0.0.0.0", 8080)}
	err := mn.Notify(ev)

	require.NoError(t, err)
	assert.Len(t, a.events, 1)
	assert.Len(t, b.events, 1)
	assert.Equal(t, ev, a.events[0])
	assert.Equal(t, ev, b.events[0])
}

func TestMultiNotifier_PartialError_ContinuesAndReturnsError(t *testing.T) {
	good := &recordingNotifier{}
	bad := &recordingNotifier{err: errors.New("boom")}
	mn := NewMultiNotifier(good, bad)

	err := mn.Notify(Event{Type: EventAppeared, Listener: makeListener("tcp", "0.0.0.0", 9090)})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "boom")
	// good notifier must still have received the event
	assert.Len(t, good.events, 1)
}

func TestMultiNotifier_Add(t *testing.T) {
	mn := NewMultiNotifier()
	assert.Equal(t, 0, mn.Len())

	mn.Add(&recordingNotifier{})
	assert.Equal(t, 1, mn.Len())

	mn.Add(&recordingNotifier{})
	assert.Equal(t, 2, mn.Len())
}

func TestMultiNotifier_AllErrors_CombinedMessage(t *testing.T) {
	a := &recordingNotifier{err: errors.New("err-a")}
	b := &recordingNotifier{err: errors.New("err-b")}
	mn := NewMultiNotifier(a, b)

	err := mn.Notify(Event{})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "2 error(s)")
	assert.Contains(t, err.Error(), "err-a")
	assert.Contains(t, err.Error(), "err-b")
}
