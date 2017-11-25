package main

import (
    "fmt"
    "math/rand"
    "time"
    "strings"
    "os"
    "github.com/nsf/termbox-go"
)


type BoardState uint
const (
    BOARD_INIT_STATE BoardState = iota
    UNDO_OK
    REDO_OK
)

type PlayerState uint
const (
    PLAYER_DEFAULT_STATE PlayerState = iota
    WIN
    LOSE
)

type GameBoard struct {
    board [][]int
    prevBoard [][]int
    size int
    score int
    prevscore int
    highscore int
    bstate BoardState
    pstate PlayerState
}

func genSquareSlice(size, inival int) [][]int {
    s := make([][]int, size)
    for i := 0; i < size; i++ {
        s[i] = make([]int, size)
        for j := 0; j < size; j++ {
            s[i][j] = inival
        }
    }
    return s
}

func GetTileColor(v int) [2]termbox.Attribute {
    var color [2]termbox.Attribute
    switch v {
    case 2:
        color[0] = termbox.Attribute(232)
        color[1] = termbox.Attribute(23)
    case 4:
        color[0] = termbox.Attribute(231)
        color[1] = termbox.Attribute(24)
    case 8:
        color[0] = termbox.Attribute(160)
        color[1] = termbox.Attribute(95)
    case 16:
        color[0] = termbox.Attribute(17)
        color[1] = termbox.Attribute(133)
    case 32:
        color[0] = termbox.Attribute(1)
        color[1] = termbox.Attribute(76)
    case 64:
        color[0] = termbox.Attribute(1)
        color[1] = termbox.Attribute(78)
    case 128:
        color[0] = termbox.Attribute(1)
        color[1] = termbox.Attribute(167)
    case 256:
        color[0] = termbox.Attribute(1)
        color[1] = termbox.Attribute(168)
    case 512:
        color[0] = termbox.Attribute(192)
        color[1] = termbox.Attribute(54)
    case 1024:
        color[0] = termbox.Attribute(192)
        color[1] = termbox.Attribute(53)
    case 2048:
        color[0] = termbox.Attribute(192)
        color[1] = termbox.Attribute(89)
    default:
        color[0] = termbox.ColorDefault
        color[1] = termbox.ColorDefault
    }
    return color
}

func NewGameBoard(size int) *GameBoard {
    g := new(GameBoard)
    g.size = 4
    g.board = genSquareSlice(size, -1)
    g.prevBoard = genSquareSlice(size, -1)
    rand.Seed(time.Now().UTC().UnixNano())
    x1 := rand.Intn(g.size)
    y1 := rand.Intn(g.size)
    x2 := rand.Intn(g.size)
    y2 := rand.Intn(g.size)
    if x1 != x2 || y1 != y2 {
        g.board[x1][y1] = 2
        g.board[x2][y2] = 2
    }
    g.bstate = BOARD_INIT_STATE
    g.pstate = PLAYER_DEFAULT_STATE
    g.score = 0
    return g
}

func (g *GameBoard) PrintBoard(offset int) int {
    drawline := func(y int) {
        str := []rune("+" + strings.Repeat("----+", g.size))
        for i, v := range str {
            termbox.SetCell(i, y, v, termbox.ColorDefault, termbox.ColorDefault)
        }
    }
    lineCount := 0
    drawline(offset)
    lineCount++ 
    for i, row := range g.board {
        y := i * 2 + offset
        for j, v := range row {
            termbox.SetCell(j*5, y+1, '|', termbox.ColorDefault, termbox.ColorDefault)
            if v != -1 {
                color := GetTileColor(v)
                for k, r := range []rune(fmt.Sprintf("%4d", v)) {
                    termbox.SetCell(j*5+1+k, y+1, r, color[0], color[1])
                }
            } else {
                for k, r := range []rune(fmt.Sprintf("    ")) {
                    termbox.SetCell(j*5+1+k, y+1, r, termbox.ColorDefault, termbox.ColorDefault)
                }
            }
        }
        termbox.SetCell((g.size+1)*4, y+1, '|', termbox.ColorDefault, termbox.ColorDefault)
        drawline(y+2)
        lineCount += 2
    }
    if g.pstate == WIN {
        g.Win()
    } else if g.pstate == LOSE {
        g.Lose()
    }
    lineCount++
    return lineCount
}

func drawLine(str string, yoff int) {
    for i, r := range []rune(str) {
        termbox.SetCell(i, yoff, r, termbox.ColorDefault, termbox.ColorDefault)
    }
}

func drawLineColor(str string, yoff int, fg, bg termbox.Attribute) {
    for i, r := range []rune(str) {
        termbox.SetCell(i, yoff, r, fg, bg)
    }
}

func (g *GameBoard) PrintScore(offset int) int {
    str := fmt.Sprintf("Score: %d / High Score: %d", g.score, g.highscore)
    drawLine(str, offset)
    return 1
}

func PrintUsage(offset int) int {
    usage := []string{
                "Slide tiles with arrow keys. If the two tiles are the same, ",
                "they merge into one and its number becomes double.",
                "GOAL: generate a '2048' tile.",
                "ESC: Exit the game / SPACE: Reset the game / PgDn: Undo / PgUp: Redo",
    }
    for i, v := range usage {
        drawLine(v, offset+i)
    }
    return len(usage)
}


func (g *GameBoard) Draw() {
    termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
    offset := 0
    offset += g.PrintScore(offset)
    offset += g.PrintBoard(offset)
    offset += PrintUsage(offset)

    termbox.Flush()
}

func (g *GameBoard) KeepPrevBoard() {
    g.bstate = REDO_OK
    for i, row := range g.board {
        for j, v := range row {
            g.prevBoard[i][j] = v
        }
    }
    g.prevscore = g.score
}

func (g *GameBoard) swapBoard() {
    for i, row := range g.board {
        for j, _ := range row {
            g.prevBoard[i][j], g.board[i][j] = g.board[i][j], g.prevBoard[i][j]
        }
    }
}

func (g *GameBoard) mirrorLR() {
    for i := 0; i < g.size; i++ {
        for j := 0; j < g.size/2; j++ {
            g.board[i][j], g.board[i][g.size-1-j] = g.board[i][g.size-1-j], g.board[i][j]
        }
    }
}

func (g *GameBoard) rrot90() {
    board := genSquareSlice(g.size, -1)
    for i := 0; i < g.size; i++ {
        for j := 0; j < g.size; j++ {
            board[i][j] = g.board[g.size - 1 - j][i]
        }
    }
    g.board = board
}

func (g *GameBoard) lrot90() {
    board := genSquareSlice(g.size, -1)
    for i := 0; i < g.size; i++ {
        for j := 0; j < g.size; j++ {
            board[i][j] = g.board[j][g.size - 1 - i]
        }
    }
    g.board = board
}

func (g *GameBoard) Undo() {
    if g.bstate == UNDO_OK {
        g.swapBoard()
        g.bstate = REDO_OK
    }
}

func (g *GameBoard) Redo() {
    if g.bstate == REDO_OK {
        g.swapBoard()
        g.bstate = UNDO_OK
    }
}

func (g *GameBoard) mergeTile(dsty, dstx, srcy, srcx int) {
    g.board[srcy][srcx] = -1
    g.board[dsty][dstx] *= 2
    if g.board[dsty][dstx] == 2048 {
        g.pstate = WIN
    }
    g.score += g.board[dsty][dstx]
    if g.score > g.highscore {
        g.highscore = g.score
    }
}

func (g *GameBoard) canSlideRight(y, x, limit int) bool {
    if x == g.size - 1 {
        return false
    }
    if g.board[y][x] == g.board[y][x+1] {
        return true
    }
    for i := x + 1; i < limit; i++ {
        if g.board[y][i] == -1 {
            return true
        }
    }
    return false
}

func (g *GameBoard) _slideRight(y, x, limit int) int {
    for i := x + 1; i < limit; i++ {
        if g.board[y][i] == -1 {
            continue
        } else if g.board[y][i] == g.board[y][x] {
            g.mergeTile(y, i, y, x)
            return i
        } else {
            g.board[y][i-1] = g.board[y][x]
            if !(i == 1 && x == 0) {
                g.board[y][x] = -1
            }
            return i
        }
    }
    g.board[y][limit - 1] = g.board[y][x]
    g.board[y][x] = -1
    return limit
}

func (g *GameBoard) SlideRight() bool {
    moved := false
    limit := g.size
    for i, row := range g.board {
        for j := g.size - 1; j >= 0; j-- {
            if row[j] == -1 {
                continue
            }
            if g.canSlideRight(i, j, limit) {
                limit = g._slideRight(i, j, limit)
                moved = true
            }
        }
        limit = g.size
    }
    return moved
}

func (g *GameBoard) PopNewTile(){
    for {
        x := rand.Intn(g.size)
        y := rand.Intn(g.size)
        if g.board[x][y] == -1 {
            if rand.Intn(100) < 90 {
                g.board[x][y] = 2
            } else {
                g.board[x][y] = 4
            }
            break
        }
    }
}

func (g *GameBoard) canSlideRightAll() bool {
    for i := 0; i < g.size; i++ {
        for j := 0; j < g.size; j++ {
            if g.canSlideRight(i, j, g.size) {
                return true
            }
        }
    }
    return false
}

func (g *GameBoard) canSlideUpAll() bool {
    g.rrot90()
    for i := 0; i < g.size; i++ {
        for j := 0; j < g.size; j++ {
            if g.canSlideRight(i, j, g.size) {
                g.lrot90()
                return true
            }
        }
    }
    g.lrot90()
    return false
}

func (g *GameBoard) canSlideDownAll() bool {
    g.lrot90()
    for i := 0; i < g.size; i++ {
        for j := 0; j < g.size; j++ {
            if g.canSlideRight(i, j, g.size) {
                g.rrot90()
                return true
            }
        }
    }
    g.rrot90()
    return false
}

func (g *GameBoard) canSlideLeftAll() bool {
    g.mirrorLR()
    for i := 0; i < g.size; i++ {
        for j := 0; j < g.size; j++ {
            if g.canSlideRight(i, j, g.size) {
                g.mirrorLR()
                return true
            }
        }
    }
    g.mirrorLR()
    return false
}

func (g *GameBoard) IsSlidable() bool {
    ret := g.canSlideRightAll()
    ret = ret || g.canSlideLeftAll()
    ret = ret || g.canSlideUpAll()
    ret = ret || g.canSlideDownAll()
    return ret
}

func (g *GameBoard) SlideUp() bool {
    g.rrot90()
    ret := g.SlideRight()
    g.lrot90()
    return ret
}

func (g *GameBoard) SlideDown() bool {
    g.lrot90()
    ret := g.SlideRight()
    g.rrot90()
    return ret
}

func (g *GameBoard) SlideLeft() bool {
    g.mirrorLR()
    ret := g.SlideRight()
    g.mirrorLR()
    return ret
}

func (g *GameBoard) CheckGameEnd() bool {
     if g.pstate == WIN {
        return true
     }
     if g.IsSlidable() {
        return true
     }
     g.pstate = LOSE
     return false
}

func (g *GameBoard) Win() {
    if g.highscore < g.score {
        g.highscore = g.score
    }
    msg := "****** YOU WIN ******"
    drawLineColor(msg, 5, 4, 240)
}

func (g *GameBoard) Lose() {
    if g.highscore < g.score {
        g.highscore = g.score
    }
    msg := "***** YOU LOSE  *****"
    drawLineColor(msg, 5, 4, 240)
}

func handleKeyEvent(ev termbox.Event) bool {
    switch ev.Type {
    case termbox.EventKey:
        switch ev.Key {
        case termbox.KeyEsc:
            return false
        case termbox.KeyArrowLeft:
            if g.pstate != PLAYER_DEFAULT_STATE {
                break
            }
            g.KeepPrevBoard()
            if g.SlideLeft() {
                g.PopNewTile()
            }
            g.Draw()
        case termbox.KeyArrowRight:
            if g.pstate != PLAYER_DEFAULT_STATE {
                break
            }
            g.KeepPrevBoard()
            if g.SlideRight() {
                g.PopNewTile()
            }
            g.Draw()
        case termbox.KeyArrowUp:
            if g.pstate != PLAYER_DEFAULT_STATE {
                break
            }
            g.KeepPrevBoard()
            if g.SlideUp() {
                g.PopNewTile()
            }
            g.Draw()
        case termbox.KeyArrowDown:
            if g.pstate != PLAYER_DEFAULT_STATE {
                break
            }
            g.KeepPrevBoard()
            if g.SlideDown() {
                g.PopNewTile()
            }
            g.Draw()
        case termbox.KeyPgup:
            if g.pstate != PLAYER_DEFAULT_STATE {
                break
            }
            g.Undo()
            g.Draw()
        case termbox.KeyPgdn:
            if g.pstate != PLAYER_DEFAULT_STATE {
                break
            }
            g.Redo()
            g.Draw()
        case termbox.KeySpace:
            g = NewGameBoard(g.size)
            g.Draw()
        default:
            g.Draw()
        }
    default:
        g.Draw()
    }
    if g.CheckGameEnd() {
        g.Draw()
    }
    return true
}

var g *GameBoard

func main() {
    err := termbox.Init()
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
    defer termbox.Close()
    termbox.SetOutputMode(termbox.Output256)

    size := 4
    g = NewGameBoard(size)
    evc := make(chan termbox.Event)
    go func() {
        for {
            evc <-termbox.PollEvent()
        }
    }()
    g.Draw()
    for {
        select {
        case ev := <-evc:
            if !handleKeyEvent(ev) {
                return
            }
            g.Draw()
        case <-time.After(1*time.Second):
            g.Draw()
        }
    }
}
