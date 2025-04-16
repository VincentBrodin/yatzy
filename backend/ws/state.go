package ws

import "github.com/google/uuid"

type State struct {
	Dice    []*Die               `json:"dice"`
	Players map[uuid.UUID]string `json:"players"`
}

type Die struct {
	Index    int  `json:"index"`
	Value    int  `json:"value"`
	Selected bool `json:"selected"`
}

func NewState() *State {
	state := &State{
		make([]*Die, 5),
		make(map[uuid.UUID]string),
	}

	for i := range state.Dice {
		state.Dice[i] = &Die{
			i,
			1,
			false,
		}
	}

	return state
}
