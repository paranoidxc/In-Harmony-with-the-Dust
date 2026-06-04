package ime

import "strings"

type Segmenter struct {
	syllables map[string]struct{}
	maxLen    int
	cache     map[string][]string
	bestCache map[string]string
}

func NewSegmenter(syllables map[string]struct{}) *Segmenter {
	maxLen := 0
	for syllable := range syllables {
		if len(syllable) > maxLen {
			maxLen = len(syllable)
		}
	}
	return &Segmenter{
		syllables: syllables,
		maxLen:    maxLen,
		cache:     make(map[string][]string),
		bestCache: make(map[string]string),
	}
}

func (s *Segmenter) Patterns(input string, limit int) []string {
	if input == "" || limit <= 0 || len(s.syllables) == 0 {
		return nil
	}
	if cached, ok := s.cache[input]; ok {
		if len(cached) > limit {
			return cached[:limit]
		}
		return cached
	}

	var results []string
	var path []string
	var dfs func(pos int)
	dfs = func(pos int) {
		if len(results) >= limit {
			return
		}
		if pos == len(input) {
			results = append(results, strings.Join(path, "'"))
			return
		}
		end := pos + s.maxLen
		if end > len(input) {
			end = len(input)
		}
		for i := end; i > pos; i-- {
			part := input[pos:i]
			if _, ok := s.syllables[part]; !ok {
				continue
			}
			path = append(path, part)
			dfs(i)
			path = path[:len(path)-1]
			if len(results) >= limit {
				return
			}
		}
	}
	dfs(0)
	s.cache[input] = append([]string(nil), results...)
	return results
}

func (s *Segmenter) BestPattern(input string) string {
	if cached, ok := s.bestCache[input]; ok {
		return cached
	}
	var path []string
	var best string
	var dfs func(pos int) bool
	dfs = func(pos int) bool {
		if pos == len(input) {
			best = strings.Join(path, "'")
			return true
		}
		end := pos + s.maxLen
		if end > len(input) {
			end = len(input)
		}
		for i := end; i > pos; i-- {
			part := input[pos:i]
			if _, ok := s.syllables[part]; !ok {
				continue
			}
			path = append(path, part)
			if dfs(i) {
				return true
			}
			path = path[:len(path)-1]
		}
		return false
	}
	dfs(0)
	s.bestCache[input] = best
	return best
}

func (s *Segmenter) BestEffortParts(input string) []string {
	if input == "" {
		return nil
	}
	if best := s.BestPattern(input); best != "" {
		return strings.Split(best, "'")
	}

	var path []string
	if s.bestEffortDFS(input, 0, &path) {
		return append([]string(nil), path...)
	}
	return splitLetters(input)
}

func (s *Segmenter) bestEffortDFS(input string, pos int, path *[]string) bool {
	if pos == len(input) {
		return true
	}
	end := pos + s.maxLen
	if end > len(input) {
		end = len(input)
	}
	for i := end; i > pos; i-- {
		part := input[pos:i]
		if _, ok := s.syllables[part]; !ok {
			continue
		}
		baseLen := len(*path)
		*path = append(*path, part)
		if s.bestEffortDFS(input, i, path) {
			return true
		}
		if i < len(input) {
			*path = append(*path, splitLetters(input[i:])...)
			return true
		}
		*path = (*path)[:baseLen]
	}
	return false
}

func initialsFromPattern(pattern string) string {
	if pattern == "" {
		return ""
	}
	parts := strings.Split(pattern, "'")
	var builder strings.Builder
	for _, part := range parts {
		if part == "" {
			continue
		}
		builder.WriteByte(part[0])
	}
	return builder.String()
}

func splitLetters(input string) []string {
	parts := make([]string, 0, len(input))
	for i := 0; i < len(input); i++ {
		parts = append(parts, input[i:i+1])
	}
	return parts
}
