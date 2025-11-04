package logic

import (
	"concurrencyAppGo/logic/types"
)

func MoveSnake(s *types.Snake) {
	s.Body[0].X += s.Direction.X
	s.Body[0].Y += s.Direction.Y
}

func CheckFoodCollision(s *types.Snake, food *[]types.Food) {
	head := s.Head()
	for i := 0; i < len(*food); i++ {
		if head.Equals((*food)[i].Pos) {
			s.Score += (*food)[i].Value
			*food = append((*food)[:i], (*food)[i+1:]...)
			return
		}
	}
}

func CheckSelfCollision(s *types.Snake) bool {
	head := s.Head()
	for _, p := range s.Body[1:] {
		if head.Equals(p) {
			return true
		}
	}
	return false
}