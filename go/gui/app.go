package classicui

import (
	"runtime"
	"sync"

	"classicui/backend/sdl2"
	"classicui/desktop"
	"classicui/event"
	"classicui/geom"
	"classicui/paint"
	uitext "classicui/text"
	"classicui/theme"
)

type Config struct {
	Title        string
	LogicalSize  geom.Size
	PresentScale int
	Theme        *theme.Theme
}

type App struct {
	config  Config
	desktop *desktop.Desktop
	canvas  *paint.Canvas
	text    *uitext.Renderer

	mu    sync.Mutex
	tasks []func()

	backend *sdl2.Backend
	running bool
}

func NewApp(cfg Config) *App {
	cfg = normalizeConfig(cfg)
	activeTheme := cfg.Theme
	if activeTheme == nil {
		activeTheme = theme.DefaultClassic()
	} else {
		activeTheme = activeTheme.Clone()
	}

	return &App{
		config:  cfg,
		desktop: desktop.New(cfg.LogicalSize, activeTheme),
		canvas:  paint.NewCanvas(cfg.LogicalSize.W, cfg.LogicalSize.H),
	}
}

func (a *App) Desktop() *desktop.Desktop {
	return a.desktop
}

func (a *App) Post(fn func()) {
	if fn == nil {
		return
	}
	a.mu.Lock()
	a.tasks = append(a.tasks, fn)
	a.mu.Unlock()
}

func (a *App) Quit() {
	a.running = false
}

func (a *App) Run() error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	backend, err := sdl2.New(sdl2.Config{
		Title:        a.config.Title,
		LogicalSize:  a.config.LogicalSize,
		PresentScale: a.config.PresentScale,
	})
	if err != nil {
		return err
	}
	defer backend.Close()
	a.backend = backend

	textRenderer, err := uitext.NewRenderer(a.desktop.Theme().Fonts)
	if err != nil {
		return err
	}
	defer textRenderer.Close()
	a.text = textRenderer

	a.running = true
	a.desktop.InvalidateAll()

	for a.running {
		a.runPostedTasks()

		evt, err := a.backend.WaitEventTimeout(16)
		if err != nil {
			return err
		}
		if evt != nil {
			if evt.Type() == event.TypeQuit {
				a.running = false
				continue
			}
			a.desktop.HandleEvent(evt)
		}

		a.runPostedTasks()

		if a.desktop.HasDirtyRegion() {
			if err := a.desktop.Paint(a.canvas, a.text); err != nil {
				return err
			}
			if err := a.backend.Present(a.canvas); err != nil {
				return err
			}
			a.desktop.ClearDirty()
		}
	}

	return nil
}

func normalizeConfig(cfg Config) Config {
	if cfg.Title == "" {
		cfg.Title = "Classic UI"
	}
	if cfg.LogicalSize.W <= 0 || cfg.LogicalSize.H <= 0 {
		cfg.LogicalSize = geom.Size{W: 640, H: 480}
	}
	if cfg.PresentScale <= 0 {
		cfg.PresentScale = 2
	}
	return cfg
}

func (a *App) runPostedTasks() {
	a.mu.Lock()
	tasks := a.tasks
	a.tasks = nil
	a.mu.Unlock()

	for _, task := range tasks {
		task()
	}
}
