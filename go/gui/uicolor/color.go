package uicolor

import (
	"fmt"
	"strconv"
	"strings"
)

type RGBA struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

func ParseHex(s string) (RGBA, error) {
	raw := strings.TrimPrefix(strings.TrimSpace(s), "#")
	switch len(raw) {
	case 6:
		raw += "FF"
	case 8:
	default:
		return RGBA{}, fmt.Errorf("invalid hex color %q", s)
	}

	value, err := strconv.ParseUint(raw, 16, 32)
	if err != nil {
		return RGBA{}, fmt.Errorf("parse color %q: %w", s, err)
	}

	return RGBA{
		R: uint8(value >> 24),
		G: uint8(value >> 16),
		B: uint8(value >> 8),
		A: uint8(value),
	}, nil
}

func MustParseHex(s string) RGBA {
	c, err := ParseHex(s)
	if err != nil {
		panic(err)
	}
	return c
}
