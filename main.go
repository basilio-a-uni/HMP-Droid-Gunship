package main

import (
	"fmt"
	"strings"
	"math/rand"
	"sync"
)

type CaseState int

const (
	EMPTY 	CaseState = iota
	CROSS
	CIRCLE
)

type GameState int

const (
	NOT_FINISHED GameState = iota
	CROSS_WON
	CIRCLE_WON
	DRAW
)

const (
	horizontalSide = "───────"
	verticalSize = "│"
	leftUpCorner = "┌"
	leftDownCorner = "└"
	rightUpCorner = "┐"
	rightDownCorner = "┘"
	threeWayLeft = "├"
	threeWayRight = "┤"
	threeWayUp = "┬"
	threeWayDown = "┴"
	fourWay = "┼"
)

type LittleGame [3][3]CaseState

type GamePos struct {
	x, y int
}

type BigGame struct {
	board 			[3][3]LittleGame
	nextPlayer 		int
	nextCase 		GamePos
	statusCache 	[3][3]Status
}

type Status struct {
	state 					GameState
	winnable1, winnable2 	bool 
}

func contains(check [3]CaseState, element CaseState) bool {
	return check[0] == element || check[1] == element || check[2] == element
}

func (game LittleGame) GetStatus() Status {
	var state, winnable1, winnable2 = DRAW, false, false
	checks := [8][3]CaseState{
		{game[0][0], game[0][1], game[0][2]},
		{game[1][0], game[1][1], game[1][2]},
		{game[2][0], game[2][1], game[2][2]},
		{game[0][0], game[1][0], game[2][0]},
		{game[0][1], game[1][1], game[2][1]},
		{game[0][2], game[1][2], game[2][2]},
		{game[0][0], game[1][1], game[2][2]},
		{game[0][2], game[1][1], game[2][0]},
	}
	for i := 0; i < 3; i++ {
		for e := 0; e < 3; e++ {
			if game[i][e] == EMPTY {
				state = NOT_FINISHED
			}
		}
	}
	for i := 0; i < 8; i++ {
		check := checks[i]
		crosses, circles := 0, 0
		for e := 0; e < 3; e++ {
			if check[e] == CROSS {
				crosses++
			} else if check[e] == CIRCLE {
				circles++
			}
		}
		if crosses == 3 {
			return Status{CROSS_WON, true, false}
		}
		if circles == 3 {
			return Status{CIRCLE_WON, false, true}
		}
		if circles == 0 {
			winnable1 = true
		}
		if crosses == 0 {
			winnable2 = true
		}
	}
	return Status{state, winnable1, winnable2}
}

func (game BigGame) GetStatus() Status {
	var state, winnable1, winnable2 = DRAW, false, false
	var results [3][3]CaseState

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			result := CaseState(game.statusCache[i][j].state)
			results[i][j] = result
			if result == EMPTY {
				state = NOT_FINISHED
			}
		}
	}
	checks := [8][3]CaseState{
		{results[0][0], results[0][1], results[0][2]},
		{results[1][0], results[1][1], results[1][2]},
		{results[2][0], results[2][1], results[2][2]},
		{results[0][0], results[1][0], results[2][0]},
		{results[0][1], results[1][1], results[2][1]},
		{results[0][2], results[1][2], results[2][2]},
		{results[0][0], results[1][1], results[2][2]},
		{results[0][2], results[1][1], results[2][0]},
	}
	for i := 0; i < 8; i++ {
		check := checks[i]
		crosses, circles := 0, 0
		for e := 0; e < 3; e++ {
			if check[e] == CROSS {
				crosses++
			} else if check[e] == CIRCLE {
				circles++
			}
		}
		if crosses == 3 {
			return Status{CROSS_WON, true, false}
		}
		if circles == 3 {
			return Status{CIRCLE_WON, false, true}
		}
		if circles == 0 {
			winnable1 = true
		}
		if crosses == 0 {
			winnable2 = true
		}

	}
	return Status{state, winnable1, winnable2}
}

func replaceSymbols(s string) string {
	result := strings.ReplaceAll(s, "0", "-")
	result = strings.ReplaceAll(result, "1", "X")
	result = strings.ReplaceAll(result, "2", "O")
	return result
}

func (game BigGame) ToString() string {
	var result string
	result += leftUpCorner + horizontalSide + threeWayUp + horizontalSide + threeWayUp + horizontalSide + rightUpCorner + "\n"
	pattern := "| %v %v %v | %v %v %v | %v %v %v |\n"
	for i := 0; i < 9; i++ {
		result += fmt.Sprintf(pattern, 
			game.board[i/3][0][i%3][0], 
			game.board[i/3][0][i%3][1],
			game.board[i/3][0][i%3][2],
			game.board[i/3][1][i%3][0], 
			game.board[i/3][1][i%3][1],
			game.board[i/3][1][i%3][2],
			game.board[i/3][2][i%3][0], 
			game.board[i/3][2][i%3][1],
			game.board[i/3][2][i%3][2])
		if i == 2 || i == 5 {
			result += threeWayLeft + horizontalSide + fourWay + horizontalSide + fourWay + horizontalSide + threeWayRight + "\n"
		}
	}
	result += leftDownCorner + horizontalSide + threeWayDown + horizontalSide + threeWayDown + horizontalSide + rightDownCorner
	return replaceSymbols(result)
}

func makeMove(game *BigGame, move [2]GamePos) {
	game.board[move[0].y][move[0].x][move[1].y][move[1].x] = CaseState(game.nextPlayer)
	game.statusCache[move[0].y][move[0].x] = game.board[move[0].y][move[0].x].GetStatus()
	game.nextPlayer = 3-game.nextPlayer
	game.nextCase = move[1]
}

func (game BigGame) GetMoves() [][2]GamePos {
	var moves [][2]GamePos
	currGamePos := game.nextCase
	status := game.statusCache[currGamePos.y][currGamePos.x]
	if status.state != NOT_FINISHED {
		for y := 0; y < 3; y++ {
			for x := 0; x < 3; x++ {
				if game.statusCache[y][x].state == CROSS_WON || game.statusCache[y][x].state == CIRCLE_WON {
					continue
				}
				for i := 0; i < 3; i++ {
					for j := 0; j < 3; j++ {
						if game.board[y][x][i][j] == EMPTY {
							moves = append(moves, [2]GamePos{GamePos{x, y}, GamePos{j, i}})
						}
					}
				}
			}
		}
	} else {
		var currGame = game.board[currGamePos.y][currGamePos.x]
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if currGame[i][j] == EMPTY {
					moves = append(moves, [2]GamePos{currGamePos, GamePos{j, i}})
				}
			}
		}
	}
	return moves
}

func Simulate(game BigGame, n int) map[GameState]int {
	var mapResults = map[GameState]int{
		NOT_FINISHED: 0,
		CROSS_WON: 0,
		CIRCLE_WON: 0,
		DRAW: 0,
	}
	for i := 0; i < n; i++ {
		curr := game
		for curr.GetStatus().state == NOT_FINISHED {
			moves := curr.GetMoves()
			move := moves[rand.Intn(len(moves))]
			makeMove(&curr, move)
		}
		mapResults[curr.GetStatus().state]++
	}
	return mapResults
}

func (game *BigGame) Explore(n int) [2]GamePos {
	simulations := n
	moves := game.GetMoves()
	lenMoves := len(moves)
	simPerMove := simulations/lenMoves
	movesScores := make(map[[2]GamePos]int)
	for _, m := range moves {
		gameToSim := *game
		player := gameToSim.nextPlayer
		makeMove(&gameToSim, m)

		mapResults := Simulate(gameToSim, simPerMove)

		if player == 1 { 
			movesScores[m] = mapResults[CROSS_WON] - mapResults[CIRCLE_WON]
		} else {
			movesScores[m] = mapResults[CIRCLE_WON] - mapResults[CROSS_WON]
		}
	}
	
	var bestMove [2]GamePos
	var bestScore = -simulations

	for key, value := range movesScores {
		if value > bestScore {
			bestScore = value
			bestMove = key
		}
	}
	return bestMove
}
func main() {
	mainGame := new(BigGame)
	mainGame.nextCase = GamePos{1, 1}
	mainGame.nextPlayer = 1
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			mainGame.statusCache[i][j] = Status{NOT_FINISHED, true, true}
		}
	}

	player := CROSS

	var mapResults = map[GameState]int{
		NOT_FINISHED: 0,
		CROSS_WON: 0,
		CIRCLE_WON: 0,
		DRAW: 0,
	}	
	resultsChan := make(chan GameState)
	wg := new(sync.WaitGroup)

	go func() {
		for result := range resultsChan {
			mapResults[result]++
		}
	}()
	
	threads := 25
	games := 25
	gamesPerThread := games/threads

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() { 
			defer wg.Done()
			for e := 0; e < gamesPerThread; e++ {
				currGame := *mainGame

				for currGame.GetStatus().state == NOT_FINISHED {
					if CaseState(currGame.nextPlayer) == player {
						move := currGame.Explore(100000)
						makeMove(&currGame, move)
					} else {
						move := currGame.Explore(100000)
						makeMove(&currGame, move)
					}
				}
				resultsChan <- currGame.GetStatus().state
			}
		}()
	}
	wg.Wait()
	close(resultsChan)

	fmt.Println(mapResults)
}
