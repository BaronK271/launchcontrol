package kiiroo

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/funjack/launchcontrol/protocol"
)

var (
	// ErrEventFormat is the error returned the event could not be parsed.
	ErrEventFormat = errors.New("invalid event format")
	// ErrNoEvents is the error returned when no events could be detected.
	ErrNoEvents = errors.New("no events found")
)

// Algorithm interface converts Kiiroo events into TimedActions.
type Algorithm interface {
	Actions(es Events) []protocol.TimedAction
}

// Event contains the values of a single Kiiroo event.
type Event struct {
	Time  time.Duration
	Value int
}

// MarshalText implements the encoding.TextMarshaler interface.
func (e Event) MarshalText() (text []byte, err error) {
	text = []byte(fmt.Sprintf("%.2f:%d", e.Time.Seconds(), e.Value))
	return text, err
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (e *Event) UnmarshalText(text []byte) error {
	s := strings.Split(string(text), ":")
	if len(s) != 2 {
		return ErrEventFormat
	}
	t, err := strconv.ParseFloat(s[0], 64)
	if err != nil {
		return ErrEventFormat
	}
	timestamp := time.Duration(int64(t*1000)) * time.Millisecond
	value, err := strconv.ParseInt(s[1], 10, 16)
	if err != nil || value < 0 || value > 4 {
		return ErrEventFormat
	}
	*e = Event{
		Time:  timestamp,
		Value: int(value),
	}
	return nil
}

// Events is an ordered series of Event objects.
type Events []Event

// Len is the number of elements in the collection.
func (es Events) Len() int {
	return len(es)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (es Events) Less(i, j int) bool {
	return es[i].Time < es[j].Time
}

// Swap swaps the elements with indexes i and j.
func (es Events) Swap(i, j int) {
	es[i], es[j] = es[j], es[i]
}

// MarshalText implements the encoding.TextMarshaler interface.
func (es Events) MarshalText() (text []byte, err error) {
	var values = make([]string, len(es))
	for i, e := range es {
		v, err := e.MarshalText()
		if err != nil {
			return []byte{}, err
		}
		values[i] = string(v)
	}
	return []byte("{" + strings.Join(values, ",") + "}"), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (es *Events) UnmarshalText(text []byte) error {
	t := strings.TrimSpace(string(text))
	size := strings.Count(t, ",")
	var events = make([]Event, size+1)

	for i, s := range strings.Split(t[1:len(t)-1], ",") {
		e := new(Event)
		err := e.UnmarshalText([]byte(s))
		if err != nil {
			return err
		}
		events[i] = *e
	}
	*es = events
	return nil
}
