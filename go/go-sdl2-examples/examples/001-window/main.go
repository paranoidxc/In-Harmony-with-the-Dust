// author: Jacky Boen

package main

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

var winTitle string = "Go-SDL2 Render"
var winWidth, winHeight int32 = 800, 600

func run() int {
	var window *sdl.Window
	var renderer *sdl.Renderer
	var renderer2 *sdl.Renderer
	var points []sdl.Point
	var rect sdl.Rect
	var rects []sdl.Rect

	window, err := sdl.CreateWindow(winTitle, 100, 100, winWidth, winHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		return 1
	}
	defer window.Destroy()

	window2, err := sdl.CreateWindow(winTitle, 300, 300, winWidth, winHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		return 1
	}
	defer window2.Destroy()

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
		return 2
	}

	renderer2, err = sdl.CreateRenderer(window2, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
		return 2
	}
	defer renderer2.Destroy()

	running := true

	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				fmt.Println("QUIT")
				running = false
			case *sdl.WindowEvent:
				if e.Event == sdl.WINDOWEVENT_CLOSE {
					win, _ := sdl.GetWindowFromID(e.WindowID)
					win.Destroy()
				}
			}
		}

		if renderer != nil {
			renderer.SetDrawColor(0, 0, 0, 255)
			renderer.Clear()

			renderer.SetDrawColor(255, 255, 255, 255)
			renderer.DrawPoint(150, 300)

			renderer.SetDrawColor(0, 0, 255, 255)
			renderer.DrawLine(0, 0, 200, 200)

			points = []sdl.Point{{0, 0}, {100, 300}, {100, 300}, {200, 0}}
			renderer.SetDrawColor(255, 255, 0, 255)
			renderer.DrawLines(points)

			rect = sdl.Rect{300, 0, 200, 200}
			renderer.SetDrawColor(255, 0, 0, 255)
			renderer.DrawRect(&rect)

			rects = []sdl.Rect{{400, 400, 100, 100}, {550, 350, 200, 200}}
			renderer.SetDrawColor(0, 255, 255, 255)
			renderer.DrawRects(rects)

			rect = sdl.Rect{250, 250, 200, 200}
			renderer.SetDrawColor(0, 255, 0, 255)
			renderer.FillRect(&rect)
			renderer.Present()
		}

		if renderer2 != nil {
			rects = []sdl.Rect{{500, 300, 100, 100}, {200, 300, 200, 200}}
			renderer2.SetDrawColor(255, 0, 255, 255)
			renderer2.FillRects(rects)
			renderer2.Present()
		}

		sdl.Delay(16)
	}

	return 0
}

func main() {
	os.Exit(run())
}
