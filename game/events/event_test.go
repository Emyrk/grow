package events_test

import (
	world2 "github.com/emyrk/grow/game/world"
	"math/rand"
	"testing"

	"github.com/emyrk/grow/game/events"
	"github.com/stretchr/testify/require"
)

func TestMarshal(t *testing.T) {
	ps := world2.NewPlayerSet()
	p := ps.AddRandomPlayer()
	testCases := []struct {
		Evts []events.Event
	}{
		{
			Evts: []events.Event{
				events.NewClickEvent(p, rand.Int(), rand.Int()),
				events.NewClickEvent(p, rand.Int(), rand.Int()),
				events.NewClickEvent(p, rand.Int(), rand.Int()),
				events.NewClickEvent(p, rand.Int(), rand.Int()),
			},
		},
	}

	for _, c := range testCases {
		t.Run("Marshaling", func(t *testing.T) {
			data, err := events.MarshalJsonEvents(c.Evts)
			require.NoError(t, err, "marshal evts")

			newEvts, err := events.UnmarshalJsonEvents(data)
			require.NoError(t, err, "unmarshal evts")

			require.Equal(t, len(c.Evts), len(newEvts), "same amt")
			require.ElementsMatch(t, c.Evts, newEvts, "events match")
		})
	}
}
