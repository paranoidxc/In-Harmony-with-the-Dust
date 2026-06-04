package ime

import "strings"

type QueryPattern struct {
	RawParts    []string
	Joined      string
	Initials    string
	HasBoundary bool
}

func parseQueryPattern(input string, segmenter *Segmenter) QueryPattern {
	trimmed := strings.ToLower(strings.TrimSpace(input))
	if trimmed == "" {
		return QueryPattern{}
	}

	rawParts := strings.Split(trimmed, "'")
	parts := make([]string, 0, len(rawParts))
	var joined strings.Builder
	var initials strings.Builder
	hasBoundary := len(rawParts) > 1
	for _, raw := range rawParts {
		part := normalizeAlphaOnly(raw)
		if part == "" {
			continue
		}
		subParts := []string{part}
		if segmenter != nil {
			subParts = segmenter.BestEffortParts(part)
		}
		for _, subPart := range subParts {
			if subPart == "" {
				continue
			}
			parts = append(parts, subPart)
			joined.WriteString(subPart)
			initials.WriteByte(subPart[0])
		}
	}

	return QueryPattern{
		RawParts:    parts,
		Joined:      joined.String(),
		Initials:    initials.String(),
		HasBoundary: hasBoundary && len(parts) > 0,
	}
}

func autoQueryPattern(input string, segmenter *Segmenter) QueryPattern {
	normalized := normalizeKey(input)
	if normalized == "" {
		return QueryPattern{}
	}

	parts := []string{normalized}
	if segmenter != nil {
		parts = segmenter.BestEffortParts(normalized)
	}

	var joined strings.Builder
	var initials strings.Builder
	for _, part := range parts {
		if part == "" {
			continue
		}
		joined.WriteString(part)
		initials.WriteByte(part[0])
	}

	return QueryPattern{
		RawParts:    parts,
		Joined:      joined.String(),
		Initials:    initials.String(),
		HasBoundary: len(parts) > 1,
	}
}

func patternMatchesQuery(candidatePattern string, query QueryPattern) bool {
	if !query.HasBoundary || len(query.RawParts) == 0 || candidatePattern == "" {
		return false
	}

	candidateParts := strings.Split(candidatePattern, "'")
	filtered := make([]string, 0, len(candidateParts))
	for _, part := range candidateParts {
		part = normalizeAlphaOnly(part)
		if part == "" {
			continue
		}
		filtered = append(filtered, part)
	}
	if len(filtered) < len(query.RawParts) {
		return false
	}

	pos := 0
	last := len(query.RawParts) - 1
	for i, queryPart := range query.RawParts {
		matched := false
		for pos < len(filtered) {
			candidatePart := filtered[pos]
			pos++
			if len(queryPart) == 1 && candidatePart[0] == queryPart[0] {
				matched = true
				break
			}
			if i == last {
				if strings.HasPrefix(candidatePart, queryPart) {
					matched = true
					break
				}
				continue
			}
			if candidatePart == queryPart {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	return true
}
