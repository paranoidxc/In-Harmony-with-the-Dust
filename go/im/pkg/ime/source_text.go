package ime

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func loadQuotedTextSource(engine *Engine, path string, sourceName string, withInitials bool) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open %s dict: %w", sourceName, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	order := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		word := fields[0]
		full, initials := parseSogouKey(fields[1])
		if full == "" {
			continue
		}
		candidate := Candidate{
			Word:    word,
			Key:     full,
			Source:  sourceName,
			Order:   order,
			Pattern: fields[1],
		}
		engine.addExact(full, candidate)
		if withInitials && initials != "" {
			engine.addInitials(initials, candidate)
		}
		order++
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan %s dict: %w", sourceName, err)
	}
	return nil
}

func collectQuotedTextSyllables(engine *Engine, path string, sourceName string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open %s dict for syllables: %w", sourceName, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		parts := strings.Split(fields[1], "'")
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
		return fmt.Errorf("scan %s dict for syllables: %w", sourceName, err)
	}
	return nil
}
