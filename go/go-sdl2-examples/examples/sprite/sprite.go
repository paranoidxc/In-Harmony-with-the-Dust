// author: Jacky Boen

package main

import (
	"fmt"
	"image/png"

	"os"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	pngImagePath = "../../assets/nijia.png"
)

var winTitle string = "Go-SDL2 Render"
var winWidth, winHeight int32 = 800, 600

func imgFileToTexture(renderer *sdl.Renderer, filename string) *sdl.Texture {
	infile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer infile.Close()

	img, err := png.Decode(infile)
	if err != nil {
		panic(err)
	}

	w := img.Bounds().Max.X
	h := img.Bounds().Max.Y

	pixels := make([]byte, w*h*4)
	bIndex := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			pixels[bIndex] = byte(r / 256)
			bIndex++
			pixels[bIndex] = byte(g / 256)
			bIndex++
			pixels[bIndex] = byte(b / 256)
			bIndex++
			pixels[bIndex] = byte(a / 256)
			bIndex++
		}
	}

	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(w), int32(h))
	if err != nil {
		panic(err)
	}
	tex.Update(nil, unsafe.Pointer(&pixels[0]), w*4)

	err = tex.SetBlendMode(sdl.BLENDMODE_BLEND)
	if err != nil {
		panic(err)
	}

	return tex
}

func run() int {
	var window *sdl.Window
	var renderer *sdl.Renderer

	sdl.Init(sdl.INIT_VIDEO)

	window, err := sdl.CreateWindow(winTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		winWidth, winHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		return 1
	}
	defer window.Destroy()

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
		return 2
	}
	defer renderer.Destroy()

	tex := imgFileToTexture(renderer, pngImagePath)
	srcX := 0
	srcY := 0
	srcIdx := 0

	//tex, _ := img.LoadTexture(renderer, pngImagePath)

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			}
		}

		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()

		renderer.SetDrawColor(255, 255, 255, 255)
		renderer.DrawPoint(150, 300)

		srcY = srcIdx * 16
		//fmt.Println("srcIdx", srcIdx, " srcY", srcY)
		renderer.Copy(tex, &sdl.Rect{X: int32(srcX), Y: int32(srcY), W: 16, H: 16}, &sdl.Rect{X: 100, Y: 100, W: 32, H: 32})
		srcIdx++
		if srcIdx == 7 {
			srcIdx = 0
		}

		renderer.Present()

		sdl.Delay(1000 / 10)
	}

	return 0
}

func main() {
	os.Exit(run())
}
