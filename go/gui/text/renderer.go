package text

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"classicui/geom"
	"classicui/paint"
	"classicui/theme"
	"classicui/uicolor"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Renderer struct {
	font     *ttf.Font
	fontPath string
	lineSkip int
	cache    map[string]bitmap
}

type bitmap struct {
	Width  int
	Height int
	Pitch  int
	Pix    []byte
}

func NewRenderer(fonts theme.FontSet) (*Renderer, error) {
	fontPath, err := resolveFontPath(fonts.Candidates)
	if err != nil {
		return nil, err
	}

	size := fonts.Size
	if size <= 0 {
		size = 12
	}

	font, err := ttf.OpenFont(fontPath, size)
	if err != nil {
		return nil, fmt.Errorf("open font %q: %w", fontPath, err)
	}
	font.SetHinting(ttf.HINTING_NORMAL)

	return &Renderer{
		font:     font,
		fontPath: fontPath,
		lineSkip: font.LineSkip(),
		cache:    make(map[string]bitmap),
	}, nil
}

func (r *Renderer) Close() {
	if r == nil || r.font == nil {
		return
	}
	r.font.Close()
	r.font = nil
}

func (r *Renderer) FontPath() string {
	return r.fontPath
}

func (r *Renderer) LineHeight() int {
	if r == nil {
		return 0
	}
	return r.lineSkip
}

func (r *Renderer) MeasureString(text string) geom.Size {
	if r == nil || r.font == nil || text == "" {
		return geom.Size{}
	}
	w, h, err := r.font.SizeUTF8(text)
	if err != nil {
		return geom.Size{}
	}
	return geom.Size{W: w, H: h}
}

func (r *Renderer) DrawString(canvas *paint.Canvas, pos geom.Point, text string, color uicolor.RGBA) error {
	if r == nil || r.font == nil || text == "" {
		return nil
	}
	key := cacheKey(text, color)
	bm, ok := r.cache[key]
	if !ok {
		var err error
		bm, err = r.renderBitmap(text, color)
		if err != nil {
			return err
		}
		r.cache[key] = bm
	}
	canvas.BlitRGBA(geom.Rect{X: pos.X, Y: pos.Y, W: bm.Width, H: bm.Height}, bm.Pix, bm.Pitch)
	return nil
}

func (r *Renderer) renderBitmap(text string, color uicolor.RGBA) (bitmap, error) {
	surface, err := r.font.RenderUTF8Blended(text, sdl.Color{
		R: color.R,
		G: color.G,
		B: color.B,
		A: color.A,
	})
	if err != nil {
		return bitmap{}, fmt.Errorf("render text %q: %w", text, err)
	}
	defer surface.Free()

	converted, err := surface.ConvertFormat(uint32(sdl.PIXELFORMAT_RGBA32), 0)
	if err != nil {
		return bitmap{}, fmt.Errorf("convert text surface: %w", err)
	}
	defer converted.Free()

	pixels := converted.Pixels()
	copied := append([]byte(nil), pixels...)
	return bitmap{
		Width:  int(converted.W),
		Height: int(converted.H),
		Pitch:  int(converted.Pitch),
		Pix:    copied,
	}, nil
}

func cacheKey(text string, color uicolor.RGBA) string {
	return fmt.Sprintf("%02x%02x%02x%02x:%s", color.R, color.G, color.B, color.A, text)
}

func resolveFontPath(candidates []string) (string, error) {
	searchDirs := fontSearchDirs()
	for _, candidate := range candidates {
		if candidate == "" {
			continue
		}
		if filepath.IsAbs(candidate) {
			if fileExists(candidate) {
				return candidate, nil
			}
			continue
		}

		for _, dir := range searchDirs {
			for _, variant := range candidateVariants(candidate) {
				fullPath := filepath.Join(dir, variant)
				if fileExists(fullPath) {
					return fullPath, nil
				}
			}
		}
	}

	return "", fmt.Errorf("no usable font found from candidates %v", candidates)
}

func candidateVariants(name string) []string {
	variants := []string{name}
	if filepath.Ext(name) == "" {
		variants = append(variants,
			name+".ttf",
			name+".ttc",
			name+".otf",
		)
	}
	lower := strings.ToLower(name)
	if lower != name {
		variants = append(variants, lower)
		if filepath.Ext(lower) == "" {
			variants = append(variants, lower+".ttf", lower+".ttc", lower+".otf")
		}
	}
	return variants
}

func fontSearchDirs() []string {
	switch runtime.GOOS {
	case "darwin":
		return []string{
			"/System/Library/Fonts",
			"/System/Library/Fonts/Supplemental",
			"/Library/Fonts",
		}
	case "windows":
		if windir := os.Getenv("WINDIR"); windir != "" {
			return []string{filepath.Join(windir, "Fonts")}
		}
		return []string{`C:\Windows\Fonts`}
	default:
		return []string{
			"/usr/share/fonts",
			"/usr/local/share/fonts",
			filepath.Join(os.Getenv("HOME"), ".fonts"),
			filepath.Join(os.Getenv("HOME"), ".local", "share", "fonts"),
		}
	}
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
