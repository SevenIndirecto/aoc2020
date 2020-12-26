package main

import (
	"testing"
)

const SETUP = `Player 1:
9
2
6
3
1

Player 2:
5
8
4
7
10
`

func TestGame_Play(t *testing.T) {
	expectedScore := 306
	expectedWinner := 2

	game := MakeGame(SETUP)
	game.Play(false, REGULAR_COMBAT)
	gotScore := game.WinnerScore()
	gotWinner := game.Winner

	if gotScore != expectedScore || gotWinner != expectedWinner {
		t.Errorf(
			"Failed gameplay got winner %d with score %d, expected winner %d with score %d",
			gotWinner, gotScore, expectedWinner, expectedScore,
		)
	}
}

func TestCreateSubGame(t *testing.T) {
	gameInput := `Player 1:
4
9
8
5
2

Player 2:
3
10
1
7
6
`
	expectedSig := "1:9,8,5,2-2:10,1,7"
	game := MakeGame(gameInput)
	p1cards := game.Decks[1].Draw()
	p2cards := game.Decks[2].Draw()

	subGame := MakeSubGame(&game, p1cards, p2cards)
	got := subGame.StateSig()

	if got != expectedSig {
		t.Errorf("Got %s expected %s", got, expectedSig)
	}
}

func TestPlayRecursive(t *testing.T) {
	expected := 291
	game := MakeGame(SETUP)
	game.Play(false, RECURSIVE_COMBAT)

	got := game.WinnerScore()

	if got != expected {
		t.Errorf("Failed to play recursive combat got %d, expected %d", got, expected)
	}
}
