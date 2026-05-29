package sdl2

import (
	"fmt"
	"unsafe"

	"classicui/event"
	"classicui/geom"
	"classicui/paint"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Config struct {
	Title        string
	LogicalSize  geom.Size
	PresentScale int
}

type Backend struct {
	window       *sdl.Window
	renderer     *sdl.Renderer
	texture      *sdl.Texture
	logicalSize  geom.Size
	presentScale int
	textInputOn  bool
}

func New(cfg Config) (*Backend, error) {
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		return nil, fmt.Errorf("init sdl video: %w", err)
	}
	if err := ttf.Init(); err != nil {
		sdl.Quit()
		return nil, fmt.Errorf("init sdl_ttf: %w", err)
	}

	if cfg.PresentScale <= 0 {
		cfg.PresentScale = 1
	}

	if !sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "0") {
		ttf.Quit()
		sdl.Quit()
		return nil, fmt.Errorf("set render scale quality hint")
	}
	sdl.SetHint(sdl.HINT_IME_INTERNAL_EDITING, "0")

	window, err := sdl.CreateWindow(
		cfg.Title,
		sdl.WINDOWPOS_CENTERED,
		sdl.WINDOWPOS_CENTERED,
		int32(cfg.LogicalSize.W*cfg.PresentScale),
		int32(cfg.LogicalSize.H*cfg.PresentScale),
		sdl.WINDOW_SHOWN,
	)
	if err != nil {
		ttf.Quit()
		sdl.Quit()
		return nil, fmt.Errorf("create window: %w", err)
	}

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_SOFTWARE)
		if err != nil {
			window.Destroy()
			ttf.Quit()
			sdl.Quit()
			return nil, fmt.Errorf("create renderer: %w", err)
		}
	}

	texture, err := renderer.CreateTexture(uint32(sdl.PIXELFORMAT_RGBA32), sdl.TEXTUREACCESS_STREAMING, int32(cfg.LogicalSize.W), int32(cfg.LogicalSize.H))
	if err != nil {
		renderer.Destroy()
		window.Destroy()
		ttf.Quit()
		sdl.Quit()
		return nil, fmt.Errorf("create streaming texture: %w", err)
	}

	return &Backend{
		window:       window,
		renderer:     renderer,
		texture:      texture,
		logicalSize:  cfg.LogicalSize,
		presentScale: cfg.PresentScale,
	}, nil
}

func (b *Backend) Close() {
	if b.textInputOn {
		sdl.StopTextInput()
		b.textInputOn = false
	}
	if b.texture != nil {
		b.texture.Destroy()
		b.texture = nil
	}
	if b.renderer != nil {
		b.renderer.Destroy()
		b.renderer = nil
	}
	if b.window != nil {
		b.window.Destroy()
		b.window = nil
	}
	ttf.Quit()
	sdl.Quit()
}

func (b *Backend) WaitEventTimeout(timeoutMS int) (event.Event, error) {
	raw := sdl.WaitEventTimeout(timeoutMS)
	if raw == nil {
		return nil, nil
	}
	return b.translate(raw), nil
}

func (b *Backend) ClipboardText() string {
	text, err := sdl.GetClipboardText()
	if err != nil {
		return ""
	}
	return text
}

func (b *Backend) SetClipboardText(text string) {
	_ = sdl.SetClipboardText(text)
}

func (b *Backend) SetTextInput(active bool, rect geom.Rect) {
	if active {
		inputRect := sdl.Rect{
			X: int32(rect.X * b.presentScale),
			Y: int32(rect.Y * b.presentScale),
			W: int32(max(rect.W, 1) * b.presentScale),
			H: int32(max(rect.H, 1) * b.presentScale),
		}
		sdl.SetTextInputRect(&inputRect)
		if !b.textInputOn {
			sdl.StartTextInput()
			b.textInputOn = true
		}
		return
	}

	if b.textInputOn {
		sdl.StopTextInput()
		b.textInputOn = false
	}
}

func (b *Backend) Present(canvas *paint.Canvas) error {
	if len(canvas.Pix) == 0 {
		return nil
	}
	if err := b.texture.Update(nil, unsafe.Pointer(&canvas.Pix[0]), canvas.Width*4); err != nil {
		return fmt.Errorf("update texture: %w", err)
	}
	if err := b.renderer.SetDrawColor(0, 0, 0, 0xFF); err != nil {
		return fmt.Errorf("set draw color: %w", err)
	}
	if err := b.renderer.Clear(); err != nil {
		return fmt.Errorf("clear renderer: %w", err)
	}
	dst := sdl.Rect{
		X: 0,
		Y: 0,
		W: int32(b.logicalSize.W * b.presentScale),
		H: int32(b.logicalSize.H * b.presentScale),
	}
	if err := b.renderer.Copy(b.texture, nil, &dst); err != nil {
		return fmt.Errorf("copy texture: %w", err)
	}
	b.renderer.Present()
	return nil
}

func (b *Backend) translate(raw sdl.Event) event.Event {
	switch e := raw.(type) {
	case *sdl.QuitEvent:
		return event.Quit{}
	case *sdl.WindowEvent:
		switch e.Event {
		case sdl.WINDOWEVENT_EXPOSED:
			return event.WindowExposed{}
		case sdl.WINDOWEVENT_SIZE_CHANGED:
			return event.WindowResized{
				Size: geom.Size{
					W: int(e.Data1) / b.presentScale,
					H: int(e.Data2) / b.presentScale,
				},
			}
		}
	case *sdl.MouseMotionEvent:
		return event.MouseMove{Position: b.toLogical(int(e.X), int(e.Y))}
	case *sdl.MouseButtonEvent:
		return event.MouseButtonEvent{
			Down:     e.State == sdl.PRESSED,
			Button:   translateMouseButton(e.Button),
			Position: b.toLogical(int(e.X), int(e.Y)),
		}
	case *sdl.MouseWheelEvent:
		x, y, _ := sdl.GetMouseState()
		delta := int(e.Y)
		if e.Direction == sdl.MOUSEWHEEL_FLIPPED {
			delta = -delta
		}
		return event.MouseWheel{
			Position: b.toLogical(int(x), int(y)),
			Delta:    delta,
		}
	case *sdl.KeyboardEvent:
		return event.KeyEvent{
			Down:      e.State == sdl.PRESSED,
			Key:       translateKey(e.Keysym.Sym),
			Modifiers: translateModifiers(sdl.Keymod(e.Keysym.Mod)),
			Repeat:    e.Repeat != 0,
		}
	case *sdl.TextInputEvent:
		return event.TextInput{Text: e.GetText()}
	case *sdl.TextEditingEvent:
		return event.TextEditing{
			Text:   e.GetText(),
			Start:  int(e.Start),
			Length: int(e.Length),
		}
	}
	return nil
}

func (b *Backend) toLogical(x, y int) geom.Point {
	return geom.Point{
		X: x / b.presentScale,
		Y: y / b.presentScale,
	}
}

func translateMouseButton(button uint8) event.MouseButton {
	switch button {
	case sdl.BUTTON_LEFT:
		return event.MouseButtonLeft
	case sdl.BUTTON_MIDDLE:
		return event.MouseButtonMiddle
	case sdl.BUTTON_RIGHT:
		return event.MouseButtonRight
	default:
		return event.MouseButtonUnknown
	}
}

func translateKey(key sdl.Keycode) event.Key {
	switch key {
	case sdl.K_ESCAPE:
		return event.KeyEscape
	case sdl.K_RETURN:
		return event.KeyEnter
	case sdl.K_SPACE:
		return event.KeySpace
	case sdl.K_TAB:
		return event.KeyTab
	case sdl.K_BACKSPACE:
		return event.KeyBackspace
	case sdl.K_DELETE:
		return event.KeyDelete
	case sdl.K_LEFT:
		return event.KeyLeft
	case sdl.K_RIGHT:
		return event.KeyRight
	case sdl.K_UP:
		return event.KeyUp
	case sdl.K_DOWN:
		return event.KeyDown
	case sdl.K_HOME:
		return event.KeyHome
	case sdl.K_END:
		return event.KeyEnd
	case sdl.K_PAGEUP:
		return event.KeyPageUp
	case sdl.K_PAGEDOWN:
		return event.KeyPageDown
	case sdl.K_F2:
		return event.KeyF2
	case sdl.K_a:
		return event.KeyA
	case sdl.K_b:
		return event.KeyB
	case sdl.K_c:
		return event.KeyC
	case sdl.K_d:
		return event.KeyD
	case sdl.K_e:
		return event.KeyE
	case sdl.K_f:
		return event.KeyF
	case sdl.K_g:
		return event.KeyG
	case sdl.K_h:
		return event.KeyH
	case sdl.K_i:
		return event.KeyI
	case sdl.K_j:
		return event.KeyJ
	case sdl.K_k:
		return event.KeyK
	case sdl.K_l:
		return event.KeyL
	case sdl.K_m:
		return event.KeyM
	case sdl.K_n:
		return event.KeyN
	case sdl.K_o:
		return event.KeyO
	case sdl.K_p:
		return event.KeyP
	case sdl.K_q:
		return event.KeyQ
	case sdl.K_r:
		return event.KeyR
	case sdl.K_s:
		return event.KeyS
	case sdl.K_t:
		return event.KeyT
	case sdl.K_u:
		return event.KeyU
	case sdl.K_v:
		return event.KeyV
	case sdl.K_w:
		return event.KeyW
	case sdl.K_x:
		return event.KeyX
	case sdl.K_y:
		return event.KeyY
	case sdl.K_z:
		return event.KeyZ
	case sdl.K_LALT:
		return event.KeyLeftAlt
	case sdl.K_RALT:
		return event.KeyRightAlt
	default:
		return event.KeyUnknown
	}
}

func translateModifiers(mod sdl.Keymod) event.Modifiers {
	var out event.Modifiers
	if mod&sdl.KMOD_SHIFT != 0 {
		out |= event.ModShift
	}
	if mod&sdl.KMOD_CTRL != 0 {
		out |= event.ModCtrl
	}
	if mod&sdl.KMOD_GUI != 0 {
		out |= event.ModCtrl
	}
	if mod&sdl.KMOD_ALT != 0 {
		out |= event.ModAlt
	}
	return out
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
