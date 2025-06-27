package blackjack

import (
	"fmt"
)

type Action rune

const (
	NA        Action = 'N' // Default action
	Blackjack Action = 'B'
	Hit       Action = 'H'
	Stand     Action = 'S'
	Double    Action = 'D'
	Split     Action = 'P'
	Surrender Action = 'U'
)

func (a Action) String() string {
	return string(a)
}

func StringToActions(actionString string) ([]Action, error) {
	var actions []Action
	actionMap := map[rune]Action{
		rune(Hit):       Hit,
		rune(Stand):     Stand,
		rune(Double):    Double,
		rune(Split):     Split,
		rune(Surrender): Surrender,
	}

	for _, action := range actionString {
		mappedAction, exists := actionMap[action]
		if !exists {
			return nil, fmt.Errorf("invalid action: %c", action)
		}

		actions = append(actions, mappedAction)
	}
	return actions, nil
}
