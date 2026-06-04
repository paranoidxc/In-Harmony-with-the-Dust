package ime

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type SogouSource struct {
	Path string
}

func init() {
	RegisterSource("sogou", func(spec SourceSpec) (Source, error) {
		return NewSogouSource(spec.Path), nil
	})
}

func NewSogouSource(path string) *SogouSource {
	return &SogouSource{Path: path}
}

func (s *SogouSource) Name() string {
	return "sogou"
}

func (s *SogouSource) CollectSyllables(*Engine) error {
	return nil
}

func (s *SogouSource) Load(engine *Engine) error {
	file, err := os.Open(s.Path)
	if err != nil {
		return fmt.Errorf("open sogou dict: %w", err)
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
		full, initials := parseSogouKey(fields[0])
		if full == "" {
			continue
		}
		for _, word := range fields[1:] {
			candidate := Candidate{
				Word:    word,
				Key:     full,
				Source:  s.Name(),
				Order:   order,
				Pattern: fields[0],
			}
			engine.addExact(full, candidate)
			if initials != "" {
				engine.addInitials(initials, candidate)
			}
			order++
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan sogou dict: %w", err)
	}
	return nil
}
