package game

import (
	"concurrencyAppGo/logic/types"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func drawSnakes(screen *ebiten.Image, snakes map[int]types.Snake, sprite *ebiten.Image) {
	for _, s := range snakes {
		for _, pos := range s.Body {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(0.30, 0.30)

			offsetX := -90.0
			
			op.GeoM.Translate(float64(pos.X*16)+offsetX, float64(pos.Y*16))
			screen.DrawImage(sprite, op)
		}
	}
}

func drawFood(screen *ebiten.Image, food []types.Food) {
	for _, f := range food {
		op := &ebiten.DrawImageOptions{}
		
		var scale float64
		if f.Type == "apple" {
			scale = 0.20
		} else {
			scale = 0.25
		}
		
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(float64(f.Pos.X*16), float64(f.Pos.Y*16))
		screen.DrawImage(f.Sprite, op)
	}
}

func drawUI(screen *ebiten.Image, snakes map[int]types.Snake) {
	y := 10
	for id := 0; id < 5; id++ {
		if s, ok := snakes[id]; ok {
			text := fmt.Sprintf("Snake %d : %d puntos", (id+1), s.Score)
			ebitenutil.DebugPrintAt(screen, text, 10, y)
			y += 20
		}
	}
}