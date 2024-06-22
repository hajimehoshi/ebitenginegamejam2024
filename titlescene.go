// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 Hajime Hoshi

package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type TitleScene struct {
}

func (t *TitleScene) Update(game GameContext) error {
	if len(inpututil.AppendJustPressedKeys(nil)) > 0 {
		game.GoToGame()
	}
	return nil
}

func (t *TitleScene) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Title Scene")
}
