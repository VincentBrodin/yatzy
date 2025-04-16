package ws

type State struct {
	Dice []*Die `json:"dice"`
}

type Die struct {
	Index    int  `json:"index"`
	Value    int  `json:"value"`
	Selected bool `json:"selected"`
}

func NewState() *State {
	state := &State{
		make([]*Die, 5),
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
