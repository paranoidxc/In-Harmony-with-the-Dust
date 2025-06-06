package ui2d

import (
	"github.com/veandco/go-sdl2/sdl"
	"learn-with-game/rpg/game"
)

func (ui *ui) DrawInventory(level *game.Level) {
	playerSrcRect := ui.textureIndex[level.Player.Rune][0]
	invRect := ui.getInventoryRect()

	ui.renderer.Copy(ui.groundInventoryBackground, nil, invRect)
	ui.renderer.Copy(ui.slotBackground, nil, ui.getHelmetSlotRect())
	ui.renderer.Copy(ui.slotBackground, nil, ui.getWeaponSlotRect())

	if level.Player.Helmet != nil {
		ui.renderer.Copy(ui.textureAtlas, &ui.textureIndex[level.Player.Helmet.Rune][0], ui.getHelmetSlotRect())
	}

	if level.Player.Weapon != nil {
		ui.renderer.Copy(ui.textureAtlas, &ui.textureIndex[level.Player.Weapon.Rune][0], ui.getWeaponSlotRect())
	}
	offset := int32(float64(invRect.H) * 0.05)

	ui.renderer.Copy(ui.textureAtlas, &playerSrcRect,
		&sdl.Rect{invRect.X + invRect.X/4, invRect.Y + offset,
			invRect.W / 2, invRect.H / 2})

	for i, item := range level.Player.Items {
		itemSrcRect := ui.textureIndex[item.Rune][0]
		if item == ui.draggedItem {
			itemSize := int32(ItemSizeRatio * float32(ui.winWidth))
			ui.renderer.Copy(ui.textureAtlas, &itemSrcRect,
				&sdl.Rect{int32(ui.currentMouseState.pos.X), int32(ui.currentMouseState.pos.Y), itemSize, itemSize})
		} else {
			ui.renderer.Copy(ui.textureAtlas, &itemSrcRect,
				ui.getInventoryItemRect(i))
		}
	}
}

func (ui *ui) getHelmetSlotRect() *sdl.Rect {
	invRect := ui.getInventoryRect()
	soltSize := int32(ItemSizeRatio * float32(ui.winWidth) * 1.05)
	return &sdl.Rect{invRect.X + invRect.W/2 - soltSize/2, invRect.Y, soltSize, soltSize}
}

func (ui *ui) getWeaponSlotRect() *sdl.Rect {
	invRect := ui.getInventoryRect()
	soltSize := int32(ItemSizeRatio * float32(ui.winWidth) * 1.05)
	yoffset := int32(float32(invRect.H) * 0.18)
	xoffset := int32(float32(invRect.W) * 0.18)
	return &sdl.Rect{invRect.X + xoffset, invRect.Y + yoffset, soltSize, soltSize}
}

func (ui *ui) getInventoryRect() *sdl.Rect {
	invWidth := int32(float32(ui.winWidth) * 0.40)
	invHeight := int32(float32(ui.winHeight) * 0.75)
	offsetX := (int32(ui.winWidth) - invWidth) / 2
	offsetY := (int32(ui.winHeight) - invHeight) / 2
	return &sdl.Rect{offsetX, offsetY, invWidth, invHeight}
}

func (ui *ui) getInventoryItemRect(i int) *sdl.Rect {
	invRect := ui.getInventoryRect()
	itemSize := int32(ItemSizeRatio * float32(ui.winWidth))
	return &sdl.Rect{invRect.X + int32(i)*itemSize, invRect.Y + invRect.H - itemSize,
		itemSize, itemSize}
}

func (ui *ui) CheckEquippedItem(level *game.Level) *game.Item {
	mousePos := ui.currentMouseState.pos

	if ui.draggedItem.Typ == game.Weapon {
		r := ui.getWeaponSlotRect()
		if r.HasIntersection(&sdl.Rect{int32(mousePos.X), int32(mousePos.Y), 1, 1}) {
			return ui.draggedItem
		}
	} else if ui.draggedItem.Typ == game.Helmet {
		r := ui.getHelmetSlotRect()
		if r.HasIntersection(&sdl.Rect{int32(mousePos.X), int32(mousePos.Y), 1, 1}) {
			return ui.draggedItem
		}
	}

	return nil
}

func (ui *ui) CheckDroppedItems(level *game.Level) *game.Item {
	invRect := ui.getInventoryRect()
	mousePos := ui.currentMouseState.pos
	if invRect.HasIntersection(&sdl.Rect{int32(mousePos.X), int32(mousePos.Y), 1, 1}) {
		return nil
	}
	return ui.draggedItem
}

func (ui *ui) CheckInventoryItems(level *game.Level) *game.Item {
	if ui.currentMouseState.leftButton {
		mousePos := ui.currentMouseState.pos
		for i, item := range level.Player.Items {
			itemRect := ui.getInventoryItemRect(i)
			if itemRect.HasIntersection(&sdl.Rect{int32(mousePos.X), int32(mousePos.Y), 1, 1}) {
				return item
			}
		}
	}
	return nil
}

func (ui *ui) CheckGroundItems(level *game.Level) *game.Item {
	if !ui.currentMouseState.leftButton && ui.prevMouseState.leftButton {
		// got a left click
		items := level.Items[level.Player.Pos]
		mousePos := ui.currentMouseState.pos
		for i, item := range items {
			itemRect := ui.getGroundItemRect(i)
			if itemRect.HasIntersection(&sdl.Rect{int32(mousePos.X), int32(mousePos.Y), 1, 1}) {
				return item
			}
		}
	}

	return nil
}
