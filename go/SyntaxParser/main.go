package main

import (
	"bufio"
	"io"
	"log"
	"log/slog"
	"os"
	"path"
	"sy/syntax"
)

var buf *syntax.Buf

func main() {
	f, err := os.OpenFile("./log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	defer f.Close()
	InitLog(f)

	file := "xxx.php"
	runes, err := readFile(file)
	if err != nil {
	}

	buf = syntax.NewBuf(file, runes)
	p := syntax.ParserForLanguage(syntax.LanguageGo)
	p.ParseAll(buf)

	slog.Info("buf", slog.Any("buf", buf))
}

func readFile(filename string) ([][]rune, error) {
	runes := [][]rune{}
	file, err := os.Open(filename)
	if err != nil {
		return runes, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	line_number := 0
	for scanner.Scan() {
		line := scanner.Text()
		runes = append(runes, []rune{})
		for _, ch := range line {
			_rune := rune(ch)
			//slog.Println(" rune:", rn, string(rn))
			if string(_rune) == "\t" {
				//slog.Info(" tab:")
				for i := 0; i < 4; i++ {
					runes[line_number] = append(runes[line_number], rune(' '))
				}
			} else {
				runes[line_number] = append(runes[line_number], rune(ch))
			}
		}
		line_number++
	}
	if line_number == 0 {
		runes = append(runes, []rune{})
	}

	return runes, nil
}

func InitLog(f io.Writer) {
	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
		//Level: slog.LevelError,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				s := a.Value.Any().(*slog.Source)
				s.File = path.Base(s.File)
			} else if a.Key == slog.TimeKey {
				t := a.Value.Time()
				a.Value = slog.StringValue(t.Format("15:04:05"))
			}
			return a
		},
	}
	// NewTextHandler
	logger := slog.New(slog.NewJSONHandler(f, opts))
	slog.SetDefault(logger)
}
