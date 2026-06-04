package ime

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type SogouSyllableSource struct {
	Path string
}

func init() {
	RegisterSource("sogou-syllables", func(spec SourceSpec) (Source, error) {
		return NewSogouSyllableSource(spec.Path), nil
	})
}

func NewSogouSyllableSource(path string) *SogouSyllableSource {
	return &SogouSyllableSource{Path: path}
}

func (s *SogouSyllableSource) Name() string {
	return "sogou-syllables"
}

func (s *SogouSyllableSource) CollectSyllables(engine *Engine) error {
	file, err := os.Open(s.Path)
	if err != nil {
		return fmt.Errorf("open sogou dict for syllables: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) == 0 {
			continue
		}
		parts := strings.Split(fields[0], "'")
		if len(parts) <= 1 {
			continue
		}
		for _, part := range parts {
			part = normalizeAlphaOnly(part)
			if part == "" || len(part) == 1 {
				continue
			}
			engine.syllables[part] = struct{}{}
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan sogou dict for syllables: %w", err)
	}
	return nil
}

func (s *SogouSyllableSource) Load(*Engine) error {
	return nil
}
