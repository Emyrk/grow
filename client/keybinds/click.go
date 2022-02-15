package keybinds

import (
	"github.com/emyrk/grow/game/events"
	world2 "github.com/emyrk/grow/game/world"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// KeyWatcher watches all user keystrokes
type KeyWatcher struct {
	// You have to know who you are to create events on your own behalf
	me *world2.Player

	LClick ebiten.MouseButton
	RClick ebiten.MouseButton
}

func NewKeybinds(me *world2.Player) *KeyWatcher {
	return &KeyWatcher{
		me:     me,
		LClick: ebiten.MouseButtonLeft,
		RClick: ebiten.MouseButtonRight,
	}
}

func (k *KeyWatcher) Update() []events.Event {
	var actions []events.Event

	if inpututil.IsMouseButtonJustPressed(k.LClick) {
		x, y := ebiten.CursorPosition()
		actions = append(actions, events.NewClickEvent(k.me, true, x, y))
	}

	if inpututil.IsMouseButtonJustPressed(k.RClick) {
		x, y := ebiten.CursorPosition()
		actions = append(actions, events.NewClickEvent(k.me, false, x, y))
	}

	return actions
}
