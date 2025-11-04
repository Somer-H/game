package game

import (
	"concurrencyAppGo/logic"
	"concurrencyAppGo/logic/types"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	State            *GameState
	FoodChan         chan types.Food
	SerpentSprite    *ebiten.Image
	AppleSprite      *ebiten.Image
	BananaSprite     *ebiten.Image
	StrawberrySprite *ebiten.Image
	mousePressed     bool 
}

func NewGame() *Game {
	state := &GameState{
		Snakes: map[int]types.Snake{},
		Food:   []types.Food{},
	}

	positions := [][]types.Position{
		{{X: 5, Y: 5}},   
		{{X: 40, Y: 5}},   
		{{X: 5, Y: 32}},  
		{{X: 40, Y: 32}},  
	}
	
	for i := 0; i < 4; i++ {
		state.Snakes[i] = types.Snake{
			ID:        i,
			Body:      positions[i],
			Score:     0,
			Direction: randomDirection(),
		}
	}

	g := &Game{
		State:            state,
		FoodChan:         make(chan types.Food),
		SerpentSprite:    mustLoad("assets/snake.png"),
		AppleSprite:      mustLoad("assets/apple.png"),
		BananaSprite:     mustLoad("assets/banana.png"),
		StrawberrySprite: mustLoad("assets/strawberry.png"),
		mousePressed:     false,
	}

	go g.foodConsumer()

	for id := range state.Snakes {
		go g.snakeWorker(id)
	}

	return g
}

func (g *Game) Update() error {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if !g.mousePressed {
			g.mousePressed = true
			
			x, y := ebiten.CursorPosition()
			gridX := x / 16
			gridY := y / 16

			kind := randomFoodType()
			food := types.Food{
				Pos:    types.Position{X: gridX, Y: gridY},
				Type:   kind,
				Value:  foodValue(kind),
				Sprite: g.spriteForType(kind),
			}
			g.FoodChan <- food
		}
	} else {
		g.mousePressed = false
	}
	
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 220, G: 220, B: 220, A: 255})

	g.State.Lock()
	snakes := make(map[int]types.Snake)
	for k, v := range g.State.Snakes {
		snakes[k] = v
	}
	food := make([]types.Food, len(g.State.Food))
	copy(food, g.State.Food)
	g.State.Unlock()

	drawSnakes(screen, snakes, g.SerpentSprite)
	drawFood(screen, food)
	drawUI(screen, snakes)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 800, 600
}

func mustLoad(path string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatalf("Error cargando %s: %v", path, err)
	}
	return img
}

func (g *Game) foodConsumer() {
	for f := range g.FoodChan {
		g.State.Lock()
		g.State.Food = append(g.State.Food, f)
		g.State.Unlock()
	}
}

func (g *Game) snakeWorker(id int) {
	ticker := time.NewTicker(120 * time.Millisecond)
	for range ticker.C {
		g.State.Lock()
		s := g.State.Snakes[id]

		if len(g.State.Food) > 0 {
			target := g.State.Food[0].Pos
			head := s.Head()

			if head.X < target.X {
				s.Direction = types.Position{X: 1, Y: 0}
			} else if head.X > target.X {
				s.Direction = types.Position{X: -1, Y: 0}
			} else if head.Y < target.Y {
				s.Direction = types.Position{X: 0, Y: 1}
			} else if head.Y > target.Y {
				s.Direction = types.Position{X: 0, Y: -1}
			}
		} else {
			head := s.Head()
			maxX := 49 - 1
			maxY := 36 - 1
			
			if head.X <= 1 && s.Direction.X == -1 {
				s.Direction = randomDirection()
			} else if head.X >= maxX-1 && s.Direction.X == 1 {
				s.Direction = randomDirection()
			} else if head.Y <= 1 && s.Direction.Y == -1 {
				s.Direction = randomDirection()
			} else if head.Y >= maxY-1 && s.Direction.Y == 1 {
				s.Direction = randomDirection()
			}
		}

		actualX := s.Body[0].X + s.Direction.X
		actualY := s.Body[0].Y + s.Direction.Y

		maxX := 49
		maxY := 33
		
		if actualX < 0 || actualX > maxX || actualY < 0 ||actualY > maxY {
			s.Direction.X = -s.Direction.X
			s.Direction.Y = -s.Direction.Y
		} else {
			logic.MoveSnake(&s)
		}
		
		logic.CheckFoodCollision(&s, &g.State.Food)
		g.State.Snakes[id] = s
		g.State.Unlock()
	}
}


func randomFoodType() string {
	types := []string{"apple", "banana", "strawberry"}
	return types[rand.Intn(len(types))]
}

func foodValue(t string) int {
	switch t {
	case "apple":
		return 1
	case "banana":
		return 2
	case "strawberry":
		return 3
	default:
		return 1
	}
}

func (g *Game) spriteForType(t string) *ebiten.Image {
	switch t {
	case "apple":
		return g.AppleSprite
	case "banana":
		return g.BananaSprite
	case "strawberry":
		return g.StrawberrySprite
	default:
		return g.AppleSprite
	}
}

func randomDirection() types.Position {
	dirs := []types.Position{
		{X: 1, Y: 0},
		{X: -1, Y: 0},
		{X: 0, Y: 1},
		{X: 0, Y: -1},
	}
	return dirs[rand.Intn(len(dirs))]
}