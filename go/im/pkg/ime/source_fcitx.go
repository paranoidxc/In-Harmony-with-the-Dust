package ime

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type FCITXSource struct {
	Path string
}

func init() {
	RegisterSource("fcitx", func(spec SourceSpec) (Source, error) {
		return NewFCITXSource(spec.Path), nil
	})
}

func NewFCITXSource(path string) *FCITXSource {
	return &FCITXSource{Path: path}
}

func (s *FCITXSource) Name() string {
	return "fcitx"
}

func (s *FCITXSource) CollectSyllables(*Engine) error {
	return nil
}

func (s *FCITXSource) Load(engine *Engine) error {
	file, err := os.Open(s.Path)
	if err != nil {
		return fmt.Errorf("open fcitx dict: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	order := 0
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 2 {
			continue
		}
		key := normalizeAlphaOnly(fields[0])
		if key == "" {
			continue
		}
		bestPattern := engine.bestSegmentPattern(key)
		initials := initialsFromPattern(bestPattern)
		for _, word := range fields[1:] {
			candidate := Candidate{
				Word:    word,
				Key:     key,
				Source:  s.Name(),
				Order:   order,
				Pattern: bestPattern,
			}
			engine.addExact(key, candidate)
			if initials != "" {
				engine.addDerivedInitials(initials, candidate)
			}
			order++
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan fcitx dict: %w", err)
	}
	return nil
}
