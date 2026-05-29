package theme

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"classicui/uicolor"
)

type ColorTable struct {
	Desktop         uicolor.RGBA
	Face            uicolor.RGBA
	Light           uicolor.RGBA
	Lightest        uicolor.RGBA
	Shadow          uicolor.RGBA
	DarkShadow      uicolor.RGBA
	Window          uicolor.RGBA
	WindowText      uicolor.RGBA
	Highlight       uicolor.RGBA
	HighlightText   uicolor.RGBA
	GrayText        uicolor.RGBA
	ActiveCaption   uicolor.RGBA
	InactiveCaption uicolor.RGBA
	CaptionText     uicolor.RGBA
}

type Metrics struct {
	BorderWidth      int `json:"border_width"`
	WindowFrameInner int `json:"window_frame_inner"`
	CaptionHeight    int `json:"caption_height"`
	MenuHeight       int `json:"menu_height"`
	IconSizeSmall    int `json:"icon_size_small"`
	ScrollbarSize    int `json:"scrollbar_size"`
	ButtonMinHeight  int `json:"button_min_height"`
	CheckboxGlyph    int `json:"checkbox_glyph_size"`
	RadioGlyph       int `json:"radio_glyph_size"`
	EditPaddingX     int `json:"edit_padding_x"`
	EditPaddingY     int `json:"edit_padding_y"`
	GroupBoxTitleX   int `json:"groupbox_title_offset_x"`
	FocusRectInset   int `json:"focus_rect_inset"`
}

type FontSet struct {
	Size       int      `json:"size"`
	Candidates []string `json:"candidates"`
}

type Theme struct {
	Name    string
	Colors  ColorTable
	Metrics Metrics
	Fonts   FontSet
}

type rawTheme struct {
	Name    string            `json:"name"`
	Colors  map[string]string `json:"colors"`
	Metrics Metrics           `json:"metrics"`
	Fonts   FontSet           `json:"fonts"`
}

//go:embed default_classic.json
var defaultClassicJSON []byte

func Load(data []byte) (*Theme, error) {
	var raw rawTheme
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("decode theme: %w", err)
	}

	colors := ColorTable{}
	loaders := []struct {
		key string
		dst *uicolor.RGBA
	}{
		{"desktop", &colors.Desktop},
		{"face", &colors.Face},
		{"light", &colors.Light},
		{"lightest", &colors.Lightest},
		{"shadow", &colors.Shadow},
		{"dark_shadow", &colors.DarkShadow},
		{"window", &colors.Window},
		{"window_text", &colors.WindowText},
		{"highlight", &colors.Highlight},
		{"highlight_text", &colors.HighlightText},
		{"gray_text", &colors.GrayText},
		{"active_caption", &colors.ActiveCaption},
		{"inactive_caption", &colors.InactiveCaption},
		{"caption_text", &colors.CaptionText},
	}

	for _, loader := range loaders {
		value, ok := raw.Colors[loader.key]
		if !ok {
			return nil, fmt.Errorf("theme color %q is required", loader.key)
		}
		color, err := uicolor.ParseHex(value)
		if err != nil {
			return nil, fmt.Errorf("theme color %q: %w", loader.key, err)
		}
		*loader.dst = color
	}

	return &Theme{
		Name:    raw.Name,
		Colors:  colors,
		Metrics: raw.Metrics,
		Fonts:   raw.Fonts,
	}, nil
}

func DefaultClassic() *Theme {
	theme, err := Load(defaultClassicJSON)
	if err != nil {
		panic(err)
	}
	return theme
}

func (t *Theme) Clone() *Theme {
	if t == nil {
		return nil
	}
	copy := *t
	if len(t.Fonts.Candidates) > 0 {
		copy.Fonts.Candidates = append([]string(nil), t.Fonts.Candidates...)
	}
	return &copy
}
