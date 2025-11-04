package types

import "github.com/hajimehoshi/ebiten/v2"

type Food struct {
	Pos    Position
	Type   string
	Value  int
	Sprite *ebiten.Image
}