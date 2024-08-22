package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

type bingoBoard struct {
	board  [25]int8
	played [25]bool
}

func makeTestBoard() *bingoBoard {
	board := [25]int8{1, 2, 3, 4, 5, 1, 2, 3, 4, 5, 1, 2, 3, 4, 5, 1, 2, 3, 4, 5, 1, 2, 3, 4, 5}
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			board[i*5+j] += int8(i * 15)
		}
	}
	return &bingoBoard{board, [25]bool{}}
}

func makeBoard() *bingoBoard {
	var board [25]int8
	for col_idx := 0; col_idx < 5; col_idx++ {
		perm := rand.Perm(15)[:5]
		sort.Ints(perm)
		for row_idx, val := range perm {
			board[col_idx*5+row_idx] = (int8)((15 * col_idx) + val)
		}
	}
	return &bingoBoard{board, [25]bool{}}
}

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"

func (board *bingoBoard) String() string {
	s := ""
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			if board.played[j*5+i] {
				s += Green + fmt.Sprintf("%-2d", board.board[j*5+i])
			} else {
				s += Reset + fmt.Sprintf("%-2d", board.board[j*5+i])
			}
			if j%5 == 4 {
				s += "\n"
			} else {
				s += " "
			}
		}
	}
	return s + Reset
}

func (board *bingoBoard) reset() {
	board.played = [25]bool{}
}

func (board *bingoBoard) play(num int) int {
	col := num / 15
	for idx, val := range board.board[5*col : 5*(col+1)] {
		if val > int8(num) {
			return -1
		} else if val == int8(num) {
			board.played[5*col+idx] = true
			return 5*col + idx
		}
	}
	return -1
}

func (board *bingoBoard) is_win(index int) int {
	won := true
	ret_val := 0
	// col win
	col := index / 5
	for i := 0; i < 5; i++ {
		if !board.played[col*5+i] {
			won = false
			break
		}
	}
	if won {
		ret_val += 1
	}

	won = true

	row := index % 5
	for i := 0; i < 5; i++ {
		if !board.played[i*5+row] {
			won = false
			break
		}
	}
	if won {
		ret_val += 2
	}

	return ret_val
}

type bingoSim struct {
	boards     []*bingoBoard
	sequences  [][75]int
	wins       []int
	gameCount  int
	boardCount int
}

func makeTestSim(boardCount int, gameCount int) *bingoSim {
	boards := make([]*bingoBoard, boardCount)
	for i := 0; i < boardCount; i++ {
		boards[i] = makeTestBoard()
	}
	sequences := make([][75]int, gameCount)
	var perm [75]int
	copy(perm[:], []int{1, 2, 3, 4, 20, 35, 50, 65, 5})
	for i := 0; i < gameCount; i++ {
		sequences[i] = perm
	}
	return &bingoSim{boards, sequences, make([]int, gameCount), gameCount, boardCount}
}

func makeSim(boardCount int, gameCount int) *bingoSim {
	boards := make([]*bingoBoard, boardCount)
	for i := 0; i < boardCount; i++ {
		boards[i] = makeBoard()
	}
	sequences := make([][75]int, gameCount)
	for i := 0; i < gameCount; i++ {
		sequences[i] = [75]int(rand.Perm(75))
	}
	return &bingoSim{boards, sequences, make([]int, gameCount), gameCount, boardCount}
}

func (sim *bingoSim) simGames() {
	var time_passed time.Duration
	for game_idx := 0; game_idx < sim.gameCount; game_idx++ {
		start := time.Now()
		if game_idx%100 == 0 {
			fmt.Println(game_idx, time_passed/time.Duration(game_idx+1))
		}
		win_status := 0

		win_steps := 76
		for _, board := range sim.boards {
			for step, num := range sim.sequences[game_idx][:win_steps-1] {
				play_idx := board.play(num)
				if play_idx >= 0 {
					cur_win_status := board.is_win(play_idx)
					if cur_win_status != 0 {
						win_status = cur_win_status
						win_steps = step
						break
					}
				}
			}
		}

		if win_status == 0 {
			fmt.Println(sim.sequences[game_idx])
			fmt.Println(sim.boards[0])
			panic("No winner")
		}

		sim.wins[game_idx] = win_status
		for _, board := range sim.boards {
			board.reset()
		}

		time_passed += time.Since(start)
	}
}

func main() {
	fmt.Println("hello world")

	sim := makeSim(100000, 1000)

	sim.simGames()

	dict := make(map[int]int)
	for _, num := range sim.wins {
		dict[num] = dict[num] + 1
	}
	fmt.Println(dict)

}
