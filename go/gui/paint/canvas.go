package paint

import (
	"classicui/geom"
	"classicui/uicolor"
)

type Canvas struct {
	Pix       []byte
	Width     int
	Height    int
	clip      geom.Rect
	clipStack []geom.Rect
}

func NewCanvas(width, height int) *Canvas {
	c := &Canvas{
		Pix:    make([]byte, width*height*4),
		Width:  width,
		Height: height,
	}
	c.ResetClip()
	return c
}

func (c *Canvas) ResetClip() {
	c.clip = geom.Rect{X: 0, Y: 0, W: c.Width, H: c.Height}
	c.clipStack = c.clipStack[:0]
}

func (c *Canvas) PushClip(rect geom.Rect) {
	c.clipStack = append(c.clipStack, c.clip)
	if clipped, ok := geom.Intersect(c.clip, rect); ok {
		c.clip = clipped
		return
	}
	c.clip = geom.Rect{}
}

func (c *Canvas) PopClip() {
	if len(c.clipStack) == 0 {
		return
	}
	last := len(c.clipStack) - 1
	c.clip = c.clipStack[last]
	c.clipStack = c.clipStack[:last]
}

func (c *Canvas) Clear(color uicolor.RGBA) {
	c.FillRect(geom.Rect{X: 0, Y: 0, W: c.Width, H: c.Height}, color)
}

func (c *Canvas) DrawPixel(x, y int, color uicolor.RGBA) {
	if !c.clip.Contains(geom.Point{X: x, Y: y}) {
		return
	}
	if x < 0 || y < 0 || x >= c.Width || y >= c.Height {
		return
	}
	idx := (y*c.Width + x) * 4
	c.Pix[idx] = color.R
	c.Pix[idx+1] = color.G
	c.Pix[idx+2] = color.B
	c.Pix[idx+3] = color.A
}

func (c *Canvas) BlendPixel(x, y int, color uicolor.RGBA) {
	if !c.clip.Contains(geom.Point{X: x, Y: y}) {
		return
	}
	if x < 0 || y < 0 || x >= c.Width || y >= c.Height {
		return
	}
	if color.A == 0 {
		return
	}
	if color.A == 0xFF {
		c.DrawPixel(x, y, color)
		return
	}

	idx := (y*c.Width + x) * 4
	dstR := c.Pix[idx]
	dstG := c.Pix[idx+1]
	dstB := c.Pix[idx+2]
	a := uint16(color.A)
	invA := uint16(0xFF - color.A)
	c.Pix[idx] = uint8((uint16(color.R)*a + uint16(dstR)*invA) / 0xFF)
	c.Pix[idx+1] = uint8((uint16(color.G)*a + uint16(dstG)*invA) / 0xFF)
	c.Pix[idx+2] = uint8((uint16(color.B)*a + uint16(dstB)*invA) / 0xFF)
	c.Pix[idx+3] = 0xFF
}

func (c *Canvas) FillRect(rect geom.Rect, color uicolor.RGBA) {
	clipped, ok := geom.Intersect(rect, c.clip)
	if !ok {
		return
	}
	for y := clipped.Y; y < clipped.Bottom(); y++ {
		for x := clipped.X; x < clipped.Right(); x++ {
			idx := (y*c.Width + x) * 4
			c.Pix[idx] = color.R
			c.Pix[idx+1] = color.G
			c.Pix[idx+2] = color.B
			c.Pix[idx+3] = color.A
		}
	}
}

func (c *Canvas) DrawHLine(x, y, width int, color uicolor.RGBA) {
	if width <= 0 {
		return
	}
	c.FillRect(geom.Rect{X: x, Y: y, W: width, H: 1}, color)
}

func (c *Canvas) DrawVLine(x, y, height int, color uicolor.RGBA) {
	if height <= 0 {
		return
	}
	c.FillRect(geom.Rect{X: x, Y: y, W: 1, H: height}, color)
}

func (c *Canvas) FrameRect(rect geom.Rect, color uicolor.RGBA) {
	if rect.W <= 0 || rect.H <= 0 {
		return
	}
	c.DrawHLine(rect.X, rect.Y, rect.W, color)
	c.DrawHLine(rect.X, rect.Bottom()-1, rect.W, color)
	c.DrawVLine(rect.X, rect.Y, rect.H, color)
	c.DrawVLine(rect.Right()-1, rect.Y, rect.H, color)
}

func (c *Canvas) DrawBevel(rect geom.Rect, topLeft, bottomRight uicolor.RGBA) {
	if rect.W < 2 || rect.H < 2 {
		return
	}
	c.DrawHLine(rect.X, rect.Y, rect.W-1, topLeft)
	c.DrawVLine(rect.X, rect.Y, rect.H-1, topLeft)
	c.DrawHLine(rect.X+1, rect.Bottom()-1, rect.W-1, bottomRight)
	c.DrawVLine(rect.Right()-1, rect.Y+1, rect.H-1, bottomRight)
}

func (c *Canvas) DrawDoubleBevel(rect geom.Rect, outerTopLeft, outerBottomRight, innerTopLeft, innerBottomRight uicolor.RGBA) {
	c.DrawBevel(rect, outerTopLeft, outerBottomRight)
	c.DrawBevel(rect.Inset(1), innerTopLeft, innerBottomRight)
}

func (c *Canvas) DrawFocusRect(rect geom.Rect, color uicolor.RGBA) {
	if rect.W <= 1 || rect.H <= 1 {
		return
	}
	for x := rect.X; x < rect.Right(); x++ {
		if (x-rect.X)%2 == 0 {
			c.DrawPixel(x, rect.Y, color)
			c.DrawPixel(x, rect.Bottom()-1, color)
		}
	}
	for y := rect.Y + 1; y < rect.Bottom()-1; y++ {
		if (y-rect.Y)%2 == 0 {
			c.DrawPixel(rect.X, y, color)
			c.DrawPixel(rect.Right()-1, y, color)
		}
	}
}

func (c *Canvas) BlitRGBA(rect geom.Rect, src []byte, pitch int) {
	if rect.Empty() || pitch <= 0 {
		return
	}
	clipped, ok := geom.Intersect(rect, c.clip)
	if !ok {
		return
	}

	srcOffsetX := clipped.X - rect.X
	srcOffsetY := clipped.Y - rect.Y
	for y := 0; y < clipped.H; y++ {
		srcRow := (srcOffsetY+y)*pitch + srcOffsetX*4
		dstRow := ((clipped.Y+y)*c.Width + clipped.X) * 4
		for x := 0; x < clipped.W; x++ {
			si := srcRow + x*4
			di := dstRow + x*4
			alpha := src[si+3]
			switch alpha {
			case 0:
				continue
			case 0xFF:
				c.Pix[di] = src[si]
				c.Pix[di+1] = src[si+1]
				c.Pix[di+2] = src[si+2]
				c.Pix[di+3] = 0xFF
			default:
				dstR := c.Pix[di]
				dstG := c.Pix[di+1]
				dstB := c.Pix[di+2]
				a := uint16(alpha)
				invA := uint16(0xFF - alpha)
				c.Pix[di] = uint8((uint16(src[si])*a + uint16(dstR)*invA) / 0xFF)
				c.Pix[di+1] = uint8((uint16(src[si+1])*a + uint16(dstG)*invA) / 0xFF)
				c.Pix[di+2] = uint8((uint16(src[si+2])*a + uint16(dstB)*invA) / 0xFF)
				c.Pix[di+3] = 0xFF
			}
		}
	}
}
