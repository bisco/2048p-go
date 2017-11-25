package main

import (
	"fmt"
	"testing"
)

func TestMakingBoard(t *testing.T) {
	size := 4
	g := NewGameBoard(size)
	if g.size != size {
		t.Error("GameBoard size is should be 4")
	}
	for _, row := range g.board {
		for _, v := range row {
			if !(v == 4 || v == 2 || v == -1 || v == 0) {
				t.Error("Initial value should be 2 or 0 or -1 (as an invalid number)")
			}
		}
	}
}

type Direction int

const (
	UP Direction = iota
	DOWN
	LEFT
	RIGHT
)

func TestSlideTile(t *testing.T) {
	size := 4
	g := NewGameBoard(size)

	type testcase struct {
		x1, y1       int
		x2, y2       int
		move         bool
		d            Direction
		ansx1, ansy1 int
		ans1         int
		ansx2, ansy2 int
		ans2         int
	}

	clearBoard := func() {
		for i, row := range g.board {
			for j, _ := range row {
				g.board[i][j] = -1
			}
		}
	}

	test := func(tc testcase) bool {
		clearBoard()
		g.board[tc.y1][tc.x1] = 2
		g.board[tc.y2][tc.x2] = 2
		var result bool
		switch tc.d {
		case UP:
			result = g.SlideUp()
		case DOWN:
			result = g.SlideDown()
		case LEFT:
			result = g.SlideLeft()
		case RIGHT:
			result = g.SlideRight()

		}
		if result != tc.move {
			t.Error("Slidable check")
			return false
		}
		if tc.ans1 != -1 && g.board[tc.ansy1][tc.ansx1] != tc.ans1 {
			t.Error("ans1 <->", tc.ans1)
			return false
		}
		if tc.ans2 != -1 && g.board[tc.ansy2][tc.ansx2] != tc.ans2 {
			t.Error("ans2 <->", g.board[tc.ansy2][tc.ansx2], tc.ans2)
			return false
		}
		for i, row := range g.board {
			for j, v := range row {
				if i == tc.ansy1 && j == tc.ansx1 {
					continue
				}
				if i == tc.ansy2 && j == tc.ansx2 {
					continue
				}
				if v != -1 {
					t.Error("Initial value should -1 (as an invalid number) => (x, y) = (", j, ",", i, ",", g.board[i][j], ")")
					return false
				}
			}
		}
		return true
	}

	testcases := []testcase{
		testcase{0, 0, 1, 0, true, LEFT, 0, 0, 4, 0, 0, -1},
		testcase{0, 0, 1, 1, true, LEFT, 0, 0, 2, 0, 1, 2},
		testcase{0, 0, 0, 1, false, LEFT, 0, 0, 2, 0, 1, 2},
		testcase{0, 0, 3, 3, true, LEFT, 0, 0, 2, 0, 3, 2},
		testcase{1, 2, 2, 3, true, LEFT, 0, 2, 2, 0, 3, 2},
		testcase{2, 2, 3, 2, true, LEFT, 0, 2, 4, 0, 0, -1},

		testcase{0, 0, 1, 0, true, RIGHT, 3, 0, 4, 0, 0, -1},
		testcase{0, 0, 1, 1, true, RIGHT, 3, 0, 2, 3, 1, 2},
		testcase{3, 0, 3, 1, false, RIGHT, 3, 0, 2, 3, 1, 2},
		testcase{0, 0, 3, 3, true, RIGHT, 3, 0, 2, 3, 3, 2},
		testcase{1, 2, 2, 3, true, RIGHT, 3, 2, 2, 3, 3, 2},
		testcase{2, 2, 3, 2, true, RIGHT, 3, 2, 4, 0, 0, -1},

		testcase{0, 0, 1, 0, false, UP, 0, 0, 2, 1, 0, 2},
		testcase{0, 0, 1, 1, true, UP, 0, 0, 2, 1, 0, 2},
		testcase{3, 0, 3, 1, true, UP, 3, 0, 4, 3, 1, -1},
		testcase{0, 0, 3, 3, true, UP, 0, 0, 2, 3, 0, 2},
		testcase{1, 2, 2, 3, true, UP, 1, 0, 2, 2, 0, 2},
		testcase{2, 2, 3, 2, true, UP, 2, 0, 2, 3, 0, 2},

		testcase{0, 0, 1, 0, true, DOWN, 0, 3, 2, 1, 3, 2},
		testcase{0, 0, 1, 1, true, DOWN, 0, 3, 2, 1, 3, 2},
		testcase{3, 0, 3, 1, true, DOWN, 3, 3, 4, 3, 1, -1},
		testcase{0, 0, 3, 3, true, DOWN, 0, 3, 2, 3, 3, 2},
		testcase{1, 2, 2, 3, true, DOWN, 1, 3, 2, 2, 3, 2},
		testcase{2, 3, 3, 3, false, DOWN, 2, 3, 2, 3, 3, 2},
	}

	for _, v := range testcases {
		if test(v) {
			fmt.Printf("%v is passed\n", v)
		} else {
			fmt.Printf("%v is failed\n", v)
		}
	}
}

func TestSlideTileRight(t *testing.T) {
	size := 4
	g := NewGameBoard(size)
	clearBoard := func() {
		for i, row := range g.board {
			for j, _ := range row {
				g.board[i][j] = -1
			}
		}
	}
	clearBoard()
	start := [4][4]int{{4, 2, 2, -1},
		{2, 2, 8, -1},
		{2, 2, -1, 2},
		{4, -1, 2, -1},
	}
	ans := [4][4]int{{-1, -1, 4, 4},
		{-1, -1, 4, 8},
		{-1, -1, 2, 4},
		{-1, -1, 4, 2},
	}
	for i, row := range start {
		for j, v := range row {
			g.board[i][j] = v
		}
	}
	if !g.SlideRight() {
		t.Error("Slidable Check")
	}
	for i, row := range ans {
		for j, v := range row {
			if g.board[i][j] != v {
				t.Error("merge miss")
				fmt.Printf("(%d, %d) has %d, but expected %d\n", i, j, g.board[i][j], v)
			}
		}
	}
}

func TestLoseCheck(t *testing.T) {
	size := 4
	g := NewGameBoard(size)
	clearBoard := func() {
		for i, row := range g.board {
			for j, _ := range row {
				g.board[i][j] = -1
			}
		}
	}
	clearBoard()
	g.board = [][]int{{4, 2, 4, 2},
		{8, 16, 128, 8},
		{512, 128, 1024, 32},
		{4, 8, 256, 128}}
	if g.CheckGameEnd() {
		t.Error("End check")
	}
	if g.pstate != LOSE {
		t.Error("Result check")
	}
}

func TestContinueCheck(t *testing.T) {
	size := 4
	g := NewGameBoard(size)
	clearBoard := func() {
		for i, row := range g.board {
			for j, _ := range row {
				g.board[i][j] = -1
			}
		}
	}
	clearBoard()
	g.board = [][]int{{4, 2, 4, 2},
		{4, 16, 128, 8},
		{512, 128, 1024, 32},
		{4, 8, 256, 128}}
	if !g.CheckGameEnd() {
		t.Error("End check")
	}
	if g.pstate == LOSE {
		t.Error("Result check")
	}
}

func TestPtileSlide(t *testing.T) {
	size := 4
	g := NewGameBoard(size)
	clearBoard := func() {
		for i, row := range g.board {
			for j, _ := range row {
				g.board[i][j] = -1
			}
		}
	}
	clearBoard()
	g.board = [][]int{{-1, 4, 2, 32},
		{-1, 2, 0, 0},
		{-1, -1, 4, 32},
		{-1, -1, -1, 4}}
	mid := [4][4]int{{-1, 4, 2, 32},
		{-1, -1, 2, 1},
		{-1, -1, 4, 32},
		{-1, -1, -1, 4}}
	ans := [4][4]int{{-1, 4, 4, 64},
		{-1, -1, 4, -1},
		{-1, -1, 8, 64},
		{-1, -1, -1, 4}}
	if !g.CheckGameEnd() {
		t.Error("End check")
	}
	if g.pstate == LOSE {
		t.Error("Result check")
	}
	if !g.SlideRight() {
		t.Error("Slide Check")
	}
	for i, row := range mid {
		for j, v := range row {
			if g.board[i][j] != v {
				t.Error("merge miss at mid")
				fmt.Printf("(%d, %d) has %d, but expected %d\n", i, j, g.board[i][j], v)
			}
		}
	}
	g.ClearPtile()
	for i, row := range ans {
		for j, v := range row {
			if g.board[i][j] != v {
				t.Error("clear miss")
				fmt.Printf("(%d, %d) has %d, but expected %d\n", i, j, g.board[i][j], v)
			}
		}
	}
	if !g.CheckGameEnd() {
		t.Error("End check after merging tiles")
	}
}
