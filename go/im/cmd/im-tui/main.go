package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"

	"im/cmd/internal/sourceflags"
	"im/pkg/ime"
)

func main() {
	fcitxPath := flag.String("fcitx", "", "override path to fcitx pinyin dictionary; empty uses ime default")
	sogouPath := flag.String("sogou", "", "override path to sogou pinyin dictionary; empty uses ime default")
	limit := flag.Int("n", 8, "max candidates to render")
	var sources sourceflags.MultiValue
	flag.Var(&sources, "source", "dictionary source spec, repeatable: fcitx=path, sogou-syllables=path, sogou=path")
	flag.Parse()

	start := time.Now()
	var (
		engine *ime.Engine
		err    error
	)
	if len(sources) == 0 && *fcitxPath == "" && *sogouPath == "" {
		engine, err = ime.LoadDefaultEngine()
	} else {
		config, buildErr := sourceflags.BuildConfig(sources, *fcitxPath, *sogouPath)
		if buildErr != nil {
			log.Fatalf("build source config: %v", buildErr)
		}
		engine, err = ime.LoadEngineWithConfig(config)
	}
	if err != nil {
		log.Fatalf("load dictionary: %v", err)
	}
	loadElapsed := time.Since(start).Round(time.Millisecond)

	session := ime.NewSession(engine, *limit)
	if err := termbox.Init(); err != nil {
		log.Fatalf("init termbox: %v", err)
	}
	defer termbox.Close()

	render(session, loadElapsed, "")
	for {
		ev := termbox.PollEvent()
		if ev.Type == termbox.EventError {
			log.Fatalf("termbox event: %v", ev.Err)
		}
		if ev.Type != termbox.EventKey {
			continue
		}

		message := ""
		switch {
		case ev.Key == termbox.KeyCtrlC || ev.Key == termbox.KeyEsc:
			return
		case ev.Key == termbox.KeyBackspace || ev.Key == termbox.KeyBackspace2:
			session.Backspace()
		case ev.Key == termbox.KeyArrowUp:
			session.MoveSelection(-1)
		case ev.Key == termbox.KeyArrowDown:
			session.MoveSelection(1)
		case ev.Key == termbox.KeySpace || ev.Key == termbox.KeyEnter:
			if committed := session.CommitSelection(); committed != "" {
				message = "committed: " + committed
			}
		case ev.Ch >= '1' && ev.Ch <= '9':
			if committed := session.CommitIndex(int(ev.Ch - '1')); committed != "" {
				message = "committed: " + committed
			}
		case ev.Ch != 0:
			session.InputRune(ev.Ch)
		}
		render(session, loadElapsed, message)
	}
}

func render(session *ime.Session, loadElapsed time.Duration, message string) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	width, _ := termbox.Size()

	drawLine(0, 0, termbox.ColorCyan|termbox.AttrBold, termbox.ColorDefault, fmt.Sprintf("IM TUI Demo  load=%s  Esc/Ctrl+C quit", loadElapsed))
	drawLine(0, 2, termbox.ColorGreen|termbox.AttrBold, termbox.ColorDefault, "Committed:")
	drawWrapped(0, 3, width, termbox.ColorWhite, termbox.ColorDefault, session.CommittedText())

	drawLine(0, 6, termbox.ColorYellow|termbox.AttrBold, termbox.ColorDefault, "Input:")
	drawLine(0, 7, termbox.ColorWhite|termbox.AttrBold, termbox.ColorDefault, "> "+session.Buffer())

	drawLine(0, 9, termbox.ColorMagenta|termbox.AttrBold, termbox.ColorDefault, "Candidates:")
	for i, item := range session.Candidates() {
		fg := termbox.ColorWhite
		bg := termbox.ColorDefault
		if i == session.SelectedIndex() {
			fg = termbox.ColorBlack | termbox.AttrBold
			bg = termbox.ColorCyan
		}
		label := fmt.Sprintf("%d. %s [%s]", i+1, item.Word, item.Match)
		drawLine(0, 10+i, fg, bg, label)
	}

	logLine := formatLog(session.Log())
	drawLine(0, 20, termbox.ColorBlue|termbox.AttrBold, termbox.ColorDefault, "Log:")
	drawWrapped(0, 21, width, termbox.ColorWhite, termbox.ColorDefault, logLine)

	if message != "" {
		drawLine(0, 24, termbox.ColorGreen|termbox.AttrBold, termbox.ColorDefault, message)
	}
	termbox.Flush()
}

func formatLog(log ime.MatchLog) string {
	segments := "-"
	if len(log.Segmentations) > 0 {
		segments = strings.Join(log.Segmentations, ",")
	}
	return fmt.Sprintf(
		"query=%q normalized=%q segments=%q exact=%d initials=%d derived=%d combined=%d tail_composed=%d tail=%d prefix=%d/%d elapsed=%s",
		log.Query,
		log.Normalized,
		segments,
		log.ExactFullMatches,
		log.ExactInitialMatches,
		log.DerivedInitialMatches,
		log.CombinedMatches,
		log.TailComposedMatches,
		log.TailFallbackMatches,
		log.PrefixFullMatches,
		log.PrefixInitialMatches,
		log.Elapsed.Round(time.Microsecond),
	)
}

func drawLine(x, y int, fg, bg termbox.Attribute, text string) {
	col := x
	for _, r := range text {
		termbox.SetCell(col, y, r, fg, bg)
		col += runewidth.RuneWidth(r)
	}
}

func drawWrapped(x, y, width int, fg, bg termbox.Attribute, text string) {
	if width <= 0 {
		return
	}
	col := x
	row := y
	for _, r := range text {
		if r == '\n' {
			row++
			col = x
			continue
		}
		rw := runewidth.RuneWidth(r)
		if rw == 0 {
			rw = 1
		}
		if col+rw > width {
			row++
			col = x
		}
		termbox.SetCell(col, row, r, fg, bg)
		col += rw
	}
}
