// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 Hajime Hoshi

package game

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Field struct {
	difficulty    Difficulty
	data          *FieldData
	playerX       int
	playerY       int
	dx            int
	dy            int
	currentDepth0 int
	currentDepth1 int
	goalReached   bool

	playerImage *ebiten.Image
}

func NewField(difficulty Difficulty) *Field {
	f := &Field{
		difficulty: difficulty,
		data:       NewFieldData(difficulty),
		playerX:    1,
		playerY:    1,
	}

	f.playerImage = f.data.tilesImage.SubImage(image.Rect(1*GridSize, 0*GridSize, 2*GridSize, 1*GridSize)).(*ebiten.Image)

	return f
}

func (f *Field) IsGoalReached() bool {
	return f.goalReached
}

func (f *Field) Update() {
	if f.goalReached {
		return
	}

	const v = 3

	if f.dx != 0 || f.dy != 0 {
		if f.dx > 0 {
			f.dx += v
		} else if f.dx < 0 {
			f.dx -= v
		}
		if f.dy > 0 {
			f.dy += v
		} else if f.dy < 0 {
			f.dy -= v
		}
		if f.dx >= GridSize {
			f.playerX++
			f.dx = 0
		}
		if f.dx <= -GridSize {
			f.playerX--
			f.dx = 0
		}
		if f.dy >= GridSize {
			f.playerY++
			f.dy = 0
		}
		if f.dy <= -GridSize {
			f.playerY--
			f.dy = 0
		}
		if f.data.isGoal(f.playerX, f.playerY) {
			f.goalReached = true
		}
		return
	}

	prevX, prevY := f.playerX, f.playerY
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		if f.data.hasSwitch(prevX, prevY, f.currentDepth1) {
			f.currentDepth0++
			f.currentDepth0 %= f.data.depth0
		}
		if f.data.hasDoor(prevX, prevY, f.currentDepth0) {
			f.currentDepth1++
			f.currentDepth1 %= f.data.depth1
		}
	}

	nextX, nextY := prevX, prevY
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		nextY++
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		nextY--
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		nextX--
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		nextX++
	}
	if !f.data.passable(nextX, nextY, prevY, f.currentDepth0, f.currentDepth1) {
		return
	}
	if nextX > f.playerX {
		f.dx = v
	}
	if nextX < f.playerX {
		f.dx = -v
	}
	if nextY > f.playerY {
		f.dy = v
	}
	if nextY < f.playerY {
		f.dy = -v
	}
}

func (f *Field) Draw(screen *ebiten.Image) {
	cx := screen.Bounds().Dx() / 2
	cy := screen.Bounds().Dy() / 3 * 2
	offsetX := cx - (f.playerX*GridSize + f.dx)
	offsetY := cy + (f.playerY*GridSize + f.dy)
	f.data.Draw(screen, offsetX, offsetY, f.currentDepth0, f.currentDepth1)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(f.playerX*GridSize+f.dx), float64(-((f.playerY+1)*GridSize + f.dy)))
	op.GeoM.Translate(float64(offsetX), float64(offsetY))
	screen.DrawImage(f.playerImage, op)

	msg := "Difficulty: " + f.difficulty.String()
	msg += "\n" + fmt.Sprintf("%dF / %dF", f.data.floorNumber(f.playerY), f.data.floorCount())
	ebitenutil.DebugPrint(screen, msg)
}
