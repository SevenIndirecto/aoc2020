package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

type Card struct {
	Value int
	Next *Card
}

type Deck struct {
	Size int
	TopCard *Card
	BottomCard *Card
}

func (deck *Deck) AddTop(value int) {
	card := &Card{Value: value}
	if deck.TopCard == nil {
		deck.TopCard = card
		deck.BottomCard = card
	} else {
		deck.TopCard.Next = card
		deck.TopCard = card
	}
	deck.Size++
}

func (deck *Deck) AddBottom(value int) {
	card := &Card{Value: value}
	if deck.BottomCard == nil {
		deck.BottomCard = card
	} else {
		deck.BottomCard.Next = card
		deck.BottomCard = card
	}

	if deck.TopCard == nil {
		deck.TopCard = card
	}
	deck.Size++
}

func (deck *Deck) Draw() int {
	if deck.Size < 1 {
		panic("Trying to draw from empty deck")
	}

	card := deck.TopCard
	deck.TopCard = card.Next
	deck.Size--

	if deck.Size < 2 {
		deck.BottomCard = deck.TopCard
	}

	return card.Value
}

func (deck *Deck) String() string {
	if deck.Size < 1 {
		return ""
	}

	str := ""
	for card := deck.TopCard; card != nil; card = card.Next {
		str += fmt.Sprintf("%d, ", card.Value)
	}
	return str[:len(str)-2]
}

const (
	REGULAR_COMBAT = iota
	RECURSIVE_COMBAT
)

type Game struct {
	Players []int
	Decks map[int]*Deck
	Round int
	Winner int
	UsedConfigurations map[string]bool
	GameId int
	GamesPlayed int
}

func (game *Game) WinnerScore() int {
	if game.Winner < 1 {
		panic("Winner not determined")
	}

	deck := game.Decks[game.Winner]
	score := 0
	for i := deck.Size; i >= 1; i-- {
		score += i * deck.Draw()
	}
	return score
}

func (game *Game) Play(output bool, mode int) {
	game.GamesPlayed = game.GameId

	for ; game.Winner < 1; {
		if mode == RECURSIVE_COMBAT {
			game.ExecuteRecursiveRound(output)
		} else {
			game.ExecuteRound(output)
		}
	}

	if output {
		fmt.Printf("\n== Post-game %d results ==\n", game.GameId)
		for p, d := range game.Decks {
			fmt.Printf("Player %d's deck: %v\n", p, d)
		}
	}
}

func (game *Game) MarkCurrentConfig() {
	sig := game.StateSig()
	game.UsedConfigurations[sig] = true
}

func (game *Game) IsCurrentConfigNew() bool {
	sig := game.StateSig()
	_, exists := game.UsedConfigurations[sig]
	return !exists
}

func (game *Game) StateSig() string {
	str := "1:" + game.Decks[1].String() + "-2:" + game.Decks[2].String()
	return strings.ReplaceAll(str, " ", "")
}

func (game *Game) ExecuteRound(output bool) {
	for loser, deck := range game.Decks {
		if deck.Size < 1 {
			// Found winner
			if loser == 1 {
				game.Winner = 2
			} else {
				game.Winner = 1
			}
			return
		}
	}

	game.Round++
	if output {
		fmt.Printf("\nRound %d\n", game.Round)
		fmt.Printf("Player %d's deck: %v\n", 1, game.Decks[1])
		fmt.Printf("Player %d's deck: %v\n", 2, game.Decks[2])
	}

	values := map[int]int{
		1: game.Decks[1].Draw(),
		2: game.Decks[2].Draw(),
	}

	var winner, loser int
	if values[1] > values[2] {
		winner = 1
		loser = 2
	} else {
		winner = 2
		loser = 1
	}

	if output {
		fmt.Printf("Player %d plays: %d\n", 1, values[1])
		fmt.Printf("Player %d plays: %d\n", 2, values[2])
		fmt.Printf("Player %d wins the round!\n", winner)
	}

	game.Decks[winner].AddBottom(values[winner])
	game.Decks[winner].AddBottom(values[loser])
}

func (game *Game) ExecuteRecursiveRound(output bool) {
	for loser, deck := range game.Decks {
		// Found a winner
		if deck.Size < 1 {
			if loser == 1 {
				game.Winner = 2
			} else {
				game.Winner = 1
			}
			return
		}
	}

	if !game.IsCurrentConfigNew() {
		game.Winner = 1
		return
	}

	game.MarkCurrentConfig()
	game.Round++
	if output {
		fmt.Printf("\nRound %d - Game %d\n", game.Round, game.GameId)
		fmt.Printf("Player %d's deck: %v\n", 1, game.Decks[1])
		fmt.Printf("Player %d's deck: %v\n", 2, game.Decks[2])
	}

	// Draw Cards
	values := map[int]int{
		1: game.Decks[1].Draw(),
		2: game.Decks[2].Draw(),
	}
	if output {
		fmt.Printf("Player %d plays: %d\n", 1, values[1])
		fmt.Printf("Player %d plays: %d\n", 2, values[2])
	}

	var roundWinner, roundLoser int

	if game.Decks[1].Size >= values[1] && game.Decks[2].Size >= values[2] {
		if output {
			fmt.Println("Playing a sub-game to determine a winner...", game.Decks[1],
				game.Decks[1].Size, "Deck2:", game.Decks[2].Size, game.Decks[2])
		}
		game.GamesPlayed++
		subGame := MakeSubGame(game, values[1], values[2])
		subGame.Play(output, RECURSIVE_COMBAT)
		if output {
			fmt.Printf("...anyway, back to game %d\n", game.GameId)
		}
		roundWinner = subGame.Winner
		if roundWinner == 1 {
			roundLoser = 2
		} else {
			roundLoser = 1
		}
	} else {
		// One player does not have enough cards, win based on card value
		if values[1] > values[2] {
			roundWinner = 1
			roundLoser = 2
		} else {
			roundWinner = 2
			roundLoser = 1
		}
	}


	if output {
		fmt.Printf("Player %d wins round %d of game %d!\n", roundWinner, game.Round, game.GameId)
	}

	game.Decks[roundWinner].AddBottom(values[roundWinner])
	game.Decks[roundWinner].AddBottom(values[roundLoser])
}

func MakeSubGame(game *Game, cardsToCopyP1, cardsToCopyP2 int) Game {
	subGame := Game{GameId: game.GamesPlayed, Decks: make(map[int]*Deck), UsedConfigurations: make(map[string]bool)}
	cardsToCopyPerPlayer := map[int]int{1: cardsToCopyP1, 2: cardsToCopyP2}

	for player, cardsToCopy := range cardsToCopyPerPlayer {
		currentCard := game.Decks[player].TopCard

		for i := 0; i < cardsToCopy; i++ {
			if _, exists := subGame.Decks[player]; !exists {
				deck := Deck{}
				deck.AddTop(currentCard.Value)
				subGame.Decks[player] = &deck
			} else {
				subGame.Decks[player].AddBottom(currentCard.Value)
			}
			currentCard = currentCard.Next
		}
	}
	return subGame
}

func MakeGame(txt string) Game {
	lines := strings.Split(txt, "\n")
	game := Game{GameId: 1, Decks: make(map[int]*Deck), UsedConfigurations: make(map[string]bool)}

	re := regexp.MustCompile(`Player (\d):`)
	var currentPlayer int

	for _, line := range lines {
		if line == "" {
			continue
		}

		match := re.FindStringSubmatch(line)
		if match != nil {
			// New player
			id, _ := strconv.Atoi(match[1])
			currentPlayer = id
			game.Players = append(game.Players, currentPlayer)
			continue
		}

		val, _ := strconv.Atoi(line)
		if _, exists := game.Decks[currentPlayer]; !exists {
			deck := Deck{}
			deck.AddTop(val)
			game.Decks[currentPlayer] = &deck
		} else {
			game.Decks[currentPlayer].AddBottom(val)
		}
	}
	return game
}

func main() {
	dat, err := ioutil.ReadFile("aoc22.txt")
	if err != nil {
		panic(err)
	}

	txt := string(dat)
	game := MakeGame(txt)
	game.Play(false, REGULAR_COMBAT)
	fmt.Println("Part one:", game.WinnerScore())

	game = MakeGame(txt)
	game.Play(false, RECURSIVE_COMBAT)
	fmt.Println("Part two:", game.WinnerScore())
}
