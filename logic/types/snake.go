package types

type Snake struct {
	ID        int
	Body      []Position
	Score     int
	Direction Position
}

func (s *Snake) Head() Position {
	return s.Body[0]
}