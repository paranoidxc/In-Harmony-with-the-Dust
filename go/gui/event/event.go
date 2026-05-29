package event

import "classicui/geom"

type Type int

const (
	TypeQuit Type = iota
	TypeWindowExposed
	TypeWindowResized
	TypeMouseMove
	TypeMouseDown
	TypeMouseUp
	TypeMouseWheel
	TypeKeyDown
	TypeKeyUp
	TypeTextInput
	TypeTextEditing
)

type Event interface {
	Type() Type
}

type Quit struct{}

func (Quit) Type() Type { return TypeQuit }

type WindowExposed struct{}

func (WindowExposed) Type() Type { return TypeWindowExposed }

type WindowResized struct {
	Size geom.Size
}

func (WindowResized) Type() Type { return TypeWindowResized }

type MouseMove struct {
	Position geom.Point
}

func (MouseMove) Type() Type { return TypeMouseMove }

type MouseButton int

const (
	MouseButtonUnknown MouseButton = iota
	MouseButtonLeft
	MouseButtonMiddle
	MouseButtonRight
)

type MouseButtonEvent struct {
	Down      bool
	Button    MouseButton
	Position  geom.Point
	Modifiers Modifiers
}

func (e MouseButtonEvent) Type() Type {
	if e.Down {
		return TypeMouseDown
	}
	return TypeMouseUp
}

type MouseWheel struct {
	Position geom.Point
	Delta    int
}

func (MouseWheel) Type() Type { return TypeMouseWheel }

type Key int

const (
	KeyUnknown Key = iota
	KeyEscape
	KeyEnter
	KeySpace
	KeyTab
	KeyBackspace
	KeyDelete
	KeyLeft
	KeyRight
	KeyUp
	KeyDown
	KeyHome
	KeyEnd
	KeyPageUp
	KeyPageDown
	KeyF2
	KeyA
	KeyB
	KeyC
	KeyD
	KeyE
	KeyF
	KeyG
	KeyH
	KeyI
	KeyJ
	KeyK
	KeyL
	KeyM
	KeyN
	KeyO
	KeyP
	KeyQ
	KeyR
	KeyS
	KeyT
	KeyU
	KeyV
	KeyW
	KeyX
	KeyY
	KeyZ
	KeyLeftAlt
	KeyRightAlt
)

type Modifiers uint16

const (
	ModShift Modifiers = 1 << iota
	ModCtrl
	ModAlt
)

type KeyEvent struct {
	Down      bool
	Key       Key
	Modifiers Modifiers
	Repeat    bool
}

func (e KeyEvent) Type() Type {
	if e.Down {
		return TypeKeyDown
	}
	return TypeKeyUp
}

type TextInput struct {
	Text string
}

func (TextInput) Type() Type { return TypeTextInput }

type TextEditing struct {
	Text   string
	Start  int
	Length int
}

func (TextEditing) Type() Type { return TypeTextEditing }
