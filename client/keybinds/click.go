package keybinds

import (
	"github.com/emyrk/grow/game/events"
	"github.com/emyrk/grow/world"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// KeyWatcher watches all user keystrokes
type KeyWatcher struct {
	// You have to know who you are to create events on your own behalf
	me *world.Player

	Click ebiten.MouseButton
}

func NewKeybinds(me *world.Player) *KeyWatcher {
	return &KeyWatcher{
		me:    me,
		Click: ebiten.MouseButtonLeft,
	}
}

func (k *KeyWatcher) Update() []events.Event {
	var actions []events.Event

	if inpututil.IsMouseButtonJustPressed(k.Click) {
		x, y := ebiten.CursorPosition()
		actions = append(actions, events.NewClickEvent(k.me, x, y))
	}

	return actions
}
