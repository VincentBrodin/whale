package valj

import "github.com/VincentBrodin/valj/stage"

type Valj struct {
	stage *stage.Stage
}

func New() (*Valj, error) {
	s, err := stage.New()
	if err != nil {
		return nil, err
	}
	return &Valj{
		stage: s,
	}, nil
}

func (v *Valj)NewList(items[]string) *List {
	return &List{
		Stage: v.stage,
		Items: items,
		Size: 3,
	}
}
