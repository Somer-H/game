package game

import (
	"concurrencyAppGo/logic/types"
	"sync"
)

type GameState struct {
	sync.Mutex
	Snakes map[int]types.Snake
	Food   []types.Food
}