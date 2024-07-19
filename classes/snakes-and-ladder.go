// package classes

// import (
// 	"fmt"
// 	"sync"
// )

// var (
// 	playerId int = 1
// 	mu       sync.Mutex
// )

// type Snake struct {
// 	start, end int
// }

// func NewSnake(start, end int) *Snake {
// 	return &Snake{start: start, end: end}
// }

// func (snake *Snake) GetStart() int {
// 	return snake.start
// }

// func (snake *Snake) GetEnd() int {
// 	return snake.end
// }

// type Ladder struct {
// 	start, end int
// }

// func NewLadder(start, end int) *Ladder {
// 	return &Ladder{start: start, end: end}
// }

// func (ladder *Ladder) GetStart() int {
// 	return ladder.start
// }

// func (ladder *Ladder) GetEnd() int {
// 	return ladder.end
// }

// /*
// class Player {
//  public:
//   Player(string);
//   int getCurrentPosition() const;
//   void setCurrentPosition(int currentPosition);
//   int getId() const;
//   const string& getName() const;

//  private:
//   int id;
//   string name;
//   int currentPosition;
//   int getUniqueId();
// };
// */

// type Player struct {
// 	id, currentPosition int
// 	name                string
// }

// func NewPlayer(name string) *Player {
// 	return &Player{name: name}
// }

// func (player *Player) GetCurrentPosition() int {
// 	return player.currentPosition
// }

// func (player *Player) SetCurrentPosition(currentPosition int) {
// 	player.currentPosition = currentPosition
// }

// func (player *Player) GetId() int {
// 	return player.id
// }

// func (player *Player) GetName() string {
// 	return player.name
// }

// func (player *Player) getUniqueId() int {
// 	mu.Lock()
// 	defer mu.Unlock()

// 	id := playerId
// 	playerId++
// 	return id
// }

// func SnakeLadderGame() {
// 	fmt.Printf("Hello, world.\n")

// }
