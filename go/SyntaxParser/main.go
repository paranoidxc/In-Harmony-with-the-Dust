package main

import (
	"fmt"
	"os"
	"strings"
	"sy/syntax"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

type editor struct {
	filename string
	lines    [][]rune
	row      int
	col      int
	scroll   int
	parser   *syntax.P
	tokens   []syntax.ComputedToken
	status   string
}

const sampleSource = `package main

func main() {
	// open a .go file: go run . syntax/parser.go
	message := "hello syntax"
	println(message, 42)
}
`

func main() {
	filename := ""
	lines := splitLines(sampleSource)
	if len(os.Args) > 1 {
		filename = os.Args[1]
		loaded, err := loadLines(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "open %s: %v\n", filename, err)
			os.Exit(1)
		}
		lines = loaded
	}

	e := &editor{
		filename: filename,
		lines:    lines,
		parser:   syntax.ParserForLanguage(syntax.LanguageGo),
	}
	e.parser.ParseAll(e.buf())
	e.tokens = e.parser.Tokens()
	e.status = "Ctrl+S save, Ctrl+Q/Esc quit"

	if err := termbox.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "termbox init: %v\n", err)
		os.Exit(1)
	}
	defer termbox.Close()

	for {
		e.render()
		ev := termbox.PollEvent()
		if ev.Type != termbox.EventKey {
			continue
		}
		if !e.handleKey(ev) {
			break
		}
	}
}

func loadLines(path string) ([][]rune, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return splitLines(string(data)), nil
}

func saveLines(path string, lines [][]rune) error {
	return os.WriteFile(path, []byte(joinLines(lines)), 0644)
}

func splitLines(s string) [][]rune {
	parts := strings.Split(s, "\n")
	lines := make([][]rune, len(parts))
	for i, part := range parts {
		lines[i] = []rune(part)
	}
	return lines
}

func joinLines(lines [][]rune) string {
	var b strings.Builder
	for i, line := range lines {
		if i > 0 {
			b.WriteRune('\n')
		}
		b.WriteString(string(line))
	}
	return b.String()
}

func (e *editor) buf() *syntax.Buf {
	return syntax.NewBuf(e.filename, e.lines)
}

func (e *editor) handleKey(ev termbox.Event) bool {
	switch ev.Key {
	case termbox.KeyCtrlQ, termbox.KeyEsc:
		return false
	case termbox.KeyCtrlS:
		e.save()
	case termbox.KeyArrowUp:
		e.moveUp()
	case termbox.KeyArrowDown:
		e.moveDown()
	case termbox.KeyArrowLeft:
		e.moveLeft()
	case termbox.KeyArrowRight:
		e.moveRight()
	case termbox.KeyEnter:
		e.insertNewline()
	case termbox.KeyBackspace, termbox.KeyBackspace2:
		e.deleteBeforeCursor()
	case termbox.KeyDelete:
		e.deleteAtCursor()
	case termbox.KeySpace:
		e.insertRune(' ')
	case termbox.KeyTab:
		e.insertRune('\t')
	default:
		if ev.Ch != 0 {
			e.insertRune(ev.Ch)
		}
	}
	return true
}

func (e *editor) save() {
	if e.filename == "" {
		e.status = "no file name; run with: go run . path/to/file.go"
		return
	}
	if err := saveLines(e.filename, e.lines); err != nil {
		e.status = "save failed: " + err.Error()
		return
	}
	e.status = "saved " + e.filename
}

func (e *editor) insertRune(r rune) {
	offset := absoluteOffset(e.lines, e.row, e.col)
	line := e.lines[e.row]
	e.lines[e.row] = append(line[:e.col], append([]rune{r}, line[e.col:]...)...)
	e.col++
	e.reparse(syntax.Edit{Offset: offset, NumInserted: 1})
}

func (e *editor) insertNewline() {
	offset := absoluteOffset(e.lines, e.row, e.col)
	line := e.lines[e.row]
	right := append([]rune(nil), line[e.col:]...)
	e.lines[e.row] = append([]rune(nil), line[:e.col]...)
	e.lines = append(e.lines[:e.row+1], append([][]rune{right}, e.lines[e.row+1:]...)...)
	e.row++
	e.col = 0
	e.reparse(syntax.Edit{Offset: offset, NumInserted: 1})
}

func (e *editor) deleteBeforeCursor() {
	if e.row == 0 && e.col == 0 {
		return
	}
	if e.col > 0 {
		offset := absoluteOffset(e.lines, e.row, e.col-1)
		line := e.lines[e.row]
		e.lines[e.row] = append(line[:e.col-1], line[e.col:]...)
		e.col--
		e.reparse(syntax.Edit{Offset: offset, NumDeleted: 1})
		return
	}

	prevLen := len(e.lines[e.row-1])
	offset := absoluteOffset(e.lines, e.row-1, prevLen)
	e.lines[e.row-1] = append(e.lines[e.row-1], e.lines[e.row]...)
	e.lines = append(e.lines[:e.row], e.lines[e.row+1:]...)
	e.row--
	e.col = prevLen
	e.reparse(syntax.Edit{Offset: offset, NumDeleted: 1})
}

func (e *editor) deleteAtCursor() {
	line := e.lines[e.row]
	if e.col < len(line) {
		offset := absoluteOffset(e.lines, e.row, e.col)
		e.lines[e.row] = append(line[:e.col], line[e.col+1:]...)
		e.reparse(syntax.Edit{Offset: offset, NumDeleted: 1})
		return
	}
	if e.row >= len(e.lines)-1 {
		return
	}
	offset := absoluteOffset(e.lines, e.row, e.col)
	e.lines[e.row] = append(e.lines[e.row], e.lines[e.row+1]...)
	e.lines = append(e.lines[:e.row+1], e.lines[e.row+2:]...)
	e.reparse(syntax.Edit{Offset: offset, NumDeleted: 1})
}

func (e *editor) reparse(edit syntax.Edit) {
	e.parser.ParseAfterEdit(e.buf(), edit)
	e.tokens = e.parser.Tokens()
	e.status = fmt.Sprintf("incremental parse: offset=%d inserted=%d deleted=%d", edit.Offset, edit.NumInserted, edit.NumDeleted)
}

func (e *editor) moveUp() {
	if e.row > 0 {
		e.row--
		e.clampCol()
	}
}

func (e *editor) moveDown() {
	if e.row < len(e.lines)-1 {
		e.row++
		e.clampCol()
	}
}

func (e *editor) moveLeft() {
	if e.col > 0 {
		e.col--
		return
	}
	if e.row > 0 {
		e.row--
		e.col = len(e.lines[e.row])
	}
}

func (e *editor) moveRight() {
	if e.col < len(e.lines[e.row]) {
		e.col++
		return
	}
	if e.row < len(e.lines)-1 {
		e.row++
		e.col = 0
	}
}

func (e *editor) clampCol() {
	if e.col > len(e.lines[e.row]) {
		e.col = len(e.lines[e.row])
	}
}

func (e *editor) render() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	width, height := termbox.Size()
	if height <= 1 || width <= 0 {
		termbox.Flush()
		return
	}

	e.ensureCursorVisible(height)
	visibleRows := height - 1
	tokenIndex := 0
	lineStartOffset := uint64(0)
	for row := 0; row < len(e.lines); row++ {
		if row >= e.scroll && row < e.scroll+visibleRows {
			e.renderLine(row, row-e.scroll, width, lineStartOffset, &tokenIndex)
		}
		lineStartOffset += uint64(len(e.lines[row]))
		if row < len(e.lines)-1 {
			lineStartOffset++
		}
	}

	e.renderStatus(width, height-1)
	cursorX := screenX(e.lines[e.row], e.col)
	termbox.SetCursor(cursorX, e.row-e.scroll)
	termbox.Flush()
}

func (e *editor) renderLine(row int, y int, width int, lineStartOffset uint64, tokenIndex *int) {
	x := 0
	for col, r := range e.lines[row] {
		if x >= width {
			break
		}
		abs := lineStartOffset + uint64(col)
		fg := colorForRole(roleAt(e.tokens, tokenIndex, abs))
		if r == '\t' {
			for i := 0; i < 4 && x < width; i++ {
				termbox.SetCell(x, y, ' ', fg, termbox.ColorDefault)
				x++
			}
			continue
		}
		termbox.SetCell(x, y, r, fg, termbox.ColorDefault)
		w := runewidth.RuneWidth(r)
		if w < 1 {
			w = 1
		}
		x += w
	}
}

func (e *editor) renderStatus(width int, y int) {
	name := e.filename
	if name == "" {
		name = "<sample>"
	}
	status := fmt.Sprintf(" %s  Ln %d, Col %d  %s ", name, e.row+1, e.col+1, e.status)
	for x := 0; x < width; x++ {
		ch := ' '
		if x < len([]rune(status)) {
			ch = []rune(status)[x]
		}
		termbox.SetCell(x, y, ch, termbox.ColorBlack|termbox.AttrBold, termbox.ColorWhite)
	}
}

func (e *editor) ensureCursorVisible(height int) {
	visibleRows := height - 1
	if e.row < e.scroll {
		e.scroll = e.row
	}
	if e.row >= e.scroll+visibleRows {
		e.scroll = e.row - visibleRows + 1
	}
	if e.scroll < 0 {
		e.scroll = 0
	}
}

func roleAt(tokens []syntax.ComputedToken, tokenIndex *int, offset uint64) syntax.TokenRole {
	for *tokenIndex < len(tokens) && tokens[*tokenIndex].Offset+tokens[*tokenIndex].Length <= offset {
		*tokenIndex++
	}
	if *tokenIndex < len(tokens) {
		tok := tokens[*tokenIndex]
		if tok.Offset <= offset && offset < tok.Offset+tok.Length {
			return tok.Role
		}
	}
	return syntax.TokenRoleNone
}

func colorForRole(role syntax.TokenRole) termbox.Attribute {
	switch role {
	case syntax.TokenRoleKeyword:
		return termbox.ColorCyan | termbox.AttrBold
	case syntax.TokenRoleString:
		return termbox.ColorYellow
	case syntax.TokenRoleComment:
		return termbox.ColorGreen
	case syntax.TokenRoleNumber:
		return termbox.ColorMagenta
	case syntax.TokenRoleOperator:
		return termbox.ColorWhite | termbox.AttrBold
	default:
		return termbox.ColorDefault
	}
}

func absoluteOffset(lines [][]rune, row int, col int) uint64 {
	var offset uint64
	for i := 0; i < row; i++ {
		offset += uint64(len(lines[i]) + 1)
	}
	return offset + uint64(col)
}

func screenX(line []rune, col int) int {
	x := 0
	for _, r := range line[:col] {
		if r == '\t' {
			x += 4
			continue
		}
		w := runewidth.RuneWidth(r)
		if w < 1 {
			w = 1
		}
		x += w
	}
	return x
}
