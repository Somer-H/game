package main

import (
	"concurrencyAppGo/game"
	"log"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
    g := game.NewGame()

    ebiten.SetWindowSize(800, 600) 
    ebiten.SetWindowTitle("Serpientes Paralelas")

    if err := ebiten.RunGame(g); err != nil {
        log.Fatal(err)
    }
}