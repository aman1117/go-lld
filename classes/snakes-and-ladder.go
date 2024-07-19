package classes

import (
	"fmt"
	"math/rand"
	"sync"
)

type Snake struct {
	start, end int
}

func NewSnake(start, end int) *Snake {
	return &Snake{start: start, end: end}
}

type Ladder struct {
	start, end int
}

func NewLadder(start, end int) *Ladder {
	return &Ladder{start: start, end: end}
}

type Player struct {
	id              int
	name            string
	currentPosition int
}

var playerIDCounter int
var playerIDMutex sync.Mutex

func NewPlayer(name string) *Player {
	playerIDMutex.Lock()
	defer playerIDMutex.Unlock()
	playerIDCounter++
	return &Player{id: playerIDCounter, name: name, currentPosition: 0}
}

func (p *Player) GetCurrentPosition() int {
	return p.currentPosition
}

func (p *Player) SetCurrentPosition(position int) {
	p.currentPosition = position
}

func (p *Player) GetName() string {
	return p.name
}

func (p *Player) GetId() int {
	return p.id
}

type Game struct {
	players          []*Player
	currentTurn      int
	winner           *Player
	snakesAndLadders map[int]int
	mu               sync.Mutex
}

func NewGame(snakes []*Snake, ladders []*Ladder, players []*Player) *Game {
	snakesAndLadders := make(map[int]int)
	for _, snake := range snakes {
		snakesAndLadders[snake.start] = snake.end
	}
	for _, ladder := range ladders {
		snakesAndLadders[ladder.start] = ladder.end
	}
	return &Game{players: players, snakesAndLadders: snakesAndLadders}
}

func (g *Game) Roll(player *Player, diceValue int) bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.winner != nil || diceValue > 6 || diceValue < 1 || g.players[g.currentTurn].GetId() != player.GetId() {
		return false
	}

	destination := g.players[g.currentTurn].GetCurrentPosition() + diceValue
	if destination <= 100 {
		if end, exists := g.snakesAndLadders[destination]; exists {
			g.players[g.currentTurn].SetCurrentPosition(end)
		} else {
			g.players[g.currentTurn].SetCurrentPosition(destination)
		}
	}

	if destination == 100 {
		g.winner = g.players[g.currentTurn]
	}

	g.nextPlayer()
	return true
}

func (g *Game) nextPlayer() {
	g.currentTurn = (g.currentTurn + 1) % len(g.players)
}

func (g *Game) GetPlayers() []*Player {
	return g.players
}

func (g *Game) GetWinner() *Player {
	return g.winner
}

func SnakesAndLadder() {

	p1 := NewPlayer("Robert")
	p2 := NewPlayer("Stannis")
	p3 := NewPlayer("Renly")

	snakes := []*Snake{
		NewSnake(17, 7),
		NewSnake(54, 34),
		NewSnake(62, 19),
		NewSnake(64, 60),
		NewSnake(87, 36),
		NewSnake(92, 73),
		NewSnake(95, 75),
		NewSnake(98, 79),
	}

	ladders := []*Ladder{
		NewLadder(1, 38),
		NewLadder(4, 14),
		NewLadder(9, 31),
		NewLadder(21, 42),
		NewLadder(28, 84),
		NewLadder(51, 67),
		NewLadder(72, 91),
		NewLadder(80, 99),
	}

	players := []*Player{p1, p2, p3}

	g := NewGame(snakes, ladders, players)

	for g.GetWinner() == nil {
		diceVal := rand.Intn(6) + 1
		g.Roll(p1, diceVal)
		diceVal = rand.Intn(6) + 1
		g.Roll(p2, diceVal)
		diceVal = rand.Intn(6) + 1
		g.Roll(p3, diceVal)
	}

	fmt.Printf("The winner is: %s\n", g.GetWinner().GetName())
	fmt.Print("All Scores: ")
	for _, player := range g.GetPlayers() {
		fmt.Printf("%d ", player.GetCurrentPosition())
	}
	fmt.Println()
}
