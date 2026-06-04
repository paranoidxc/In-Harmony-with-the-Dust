package ime

import (
	"sort"
	"strings"
	"time"
)

type Candidate struct {
	Word    string
	Key     string
	Source  string
	Order   int
	Pattern string
	Match   string
}

type SearchResult struct {
	Items []Candidate
	Log   MatchLog
}

type MatchLog struct {
	Query                 string
	Normalized            string
	Segmentations         []string
	ExactFullMatches      int
	ExactInitialMatches   int
	DerivedInitialMatches int
	CombinedMatches       int
	TailComposedMatches   int
	TailFallbackMatches   int
	PrefixFullMatches     int
	PrefixInitialMatches  int
	Returned              int
	Elapsed               time.Duration
}

type Engine struct {
	exact               map[string][]Candidate
	initials            map[string][]Candidate
	derivedInitials     map[string][]Candidate
	exactKeys           []string
	initialsKeys        []string
	derivedInitialsKeys []string
	exactKeySet         map[string]struct{}
	initialsKeySet      map[string]struct{}
	derivedKeySet       map[string]struct{}
	syllables           map[string]struct{}
	segmenter           *Segmenter
}

func NewEngine() *Engine {
	return &Engine{
		exact:           make(map[string][]Candidate),
		initials:        make(map[string][]Candidate),
		derivedInitials: make(map[string][]Candidate),
		exactKeySet:     make(map[string]struct{}),
		initialsKeySet:  make(map[string]struct{}),
		derivedKeySet:   make(map[string]struct{}),
		syllables:       make(map[string]struct{}),
	}
}

func LoadEngine(fcitxPath, sogouPath string) (*Engine, error) {
	return LoadEngineWithConfig(DefaultSourceConfig(fcitxPath, sogouPath))
}

func LoadDefaultEngine() (*Engine, error) {
	return LoadEngineWithConfig(DefaultSourceConfigFromEnv())
}

func LoadEngineWithConfig(config SourceConfig) (*Engine, error) {
	engine := NewEngine()
	for _, source := range config.Sources {
		if err := source.CollectSyllables(engine); err != nil {
			return nil, err
		}
	}
	engine.segmenter = NewSegmenter(engine.syllables)
	for _, source := range config.Sources {
		if err := source.Load(engine); err != nil {
			return nil, err
		}
	}
	sort.Strings(engine.exactKeys)
	sort.Strings(engine.initialsKeys)
	sort.Strings(engine.derivedInitialsKeys)
	return engine, nil
}

func (e *Engine) Search(query string, limit int) SearchResult {
	start := time.Now()
	queryPattern := parseQueryPattern(query, e.segmenter)
	normalized := queryPattern.Joined
	if normalized == "" {
		normalized = normalizeKey(query)
	}
	autoPattern := QueryPattern{}
	hasExplicitBoundary := strings.Contains(query, "'")
	if !hasExplicitBoundary {
		autoPattern = autoQueryPattern(query, e.segmenter)
	}
	result := SearchResult{
		Log: MatchLog{
			Query:      query,
			Normalized: normalized,
		},
	}
	if normalized == "" {
		result.Log.Elapsed = time.Since(start)
		return result
	}
	if queryPattern.HasBoundary {
		result.Log.Segmentations = []string{strings.Join(queryPattern.RawParts, "'")}
	} else {
		result.Log.Segmentations = e.segmentPatterns(normalized)
	}

	seen := make(map[string]struct{}, limit)
	appendCandidates := func(items []Candidate, counter *int, match string) {
		for _, item := range items {
			*counter = *counter + 1
			if _, ok := seen[item.Word]; ok {
				continue
			}
			seen[item.Word] = struct{}{}
			copyItem := item
			copyItem.Match = match
			result.Items = append(result.Items, copyItem)
			if len(result.Items) >= limit {
				return
			}
		}
	}

	preferInitials := e.shouldPreferInitials(normalized)
	if preferInitials {
		appendCandidates(e.derivedInitials[normalized], &result.Log.DerivedInitialMatches, "derived_initials")
		if len(result.Items) < limit {
			appendCandidates(e.initials[normalized], &result.Log.ExactInitialMatches, "exact_initials")
		}
		if len(result.Items) < limit {
			appendCandidates(e.exact[normalized], &result.Log.ExactFullMatches, "exact_full")
		}
	} else {
		appendCandidates(e.exact[normalized], &result.Log.ExactFullMatches, "exact_full")
		for _, pattern := range result.Log.Segmentations {
			initials := initialsFromPattern(pattern)
			if initials == "" || initials == normalized || len(result.Items) >= limit {
				continue
			}
			appendCandidates(e.initials[initials], &result.Log.ExactInitialMatches, "segmented_initials")
		}
		if len(result.Items) < limit {
			appendCandidates(e.initials[normalized], &result.Log.ExactInitialMatches, "exact_initials")
		}
		if len(result.Items) < limit {
			appendCandidates(e.derivedInitials[normalized], &result.Log.DerivedInitialMatches, "derived_initials")
		}
	}
	if len(result.Items) < limit {
		e.appendCombinedMatches(&result, &seen, limit)
	}
	if len(result.Items) < limit && queryPattern.HasBoundary {
		e.appendPatternMatches(&result, &seen, queryPattern, limit)
	}
	if len(result.Items) == 0 && autoPattern.HasBoundary {
		result.Log.Segmentations = []string{strings.Join(autoPattern.RawParts, "'")}
		e.appendPatternMatches(&result, &seen, autoPattern, limit)
	}
	if len(result.Items) == 0 {
		e.appendTailComposedMatches(&result, &seen, limit)
	}
	if len(result.Items) == 0 && autoPattern.HasBoundary {
		e.appendTailFallbackMatches(&result, &seen, autoPattern, limit)
	}
	if len(result.Items) < limit {
		for _, key := range prefixKeys(e.exactKeys, normalized) {
			appendCandidates(e.exact[key], &result.Log.PrefixFullMatches, "prefix_full")
			if len(result.Items) >= limit {
				break
			}
		}
	}
	if len(result.Items) < limit {
		for _, key := range prefixKeys(e.initialsKeys, normalized) {
			appendCandidates(e.initials[key], &result.Log.PrefixInitialMatches, "prefix_initials")
			if len(result.Items) >= limit {
				break
			}
		}
	}
	if len(result.Items) < limit {
		for _, key := range prefixKeys(e.derivedInitialsKeys, normalized) {
			appendCandidates(e.derivedInitials[key], &result.Log.PrefixInitialMatches, "prefix_derived_initials")
			if len(result.Items) >= limit {
				break
			}
		}
	}

	result.Log.Returned = len(result.Items)
	sortCandidates(result.Items)
	result.Log.Elapsed = time.Since(start)
	return result
}

func (e *Engine) appendTailComposedMatches(result *SearchResult, seen *map[string]struct{}, limit int) {
	for _, pattern := range result.Log.Segmentations {
		parts := strings.Split(pattern, "'")
		if len(parts) < 2 {
			continue
		}
		for split := len(parts) - 1; split >= 1; split-- {
			headKey := strings.Join(parts[:split], "")
			tailKey := strings.Join(parts[split:], "")
			if headKey == "" || tailKey == "" {
				continue
			}
			headWords := e.exact[headKey]
			if len(headWords) == 0 {
				continue
			}
			tailWords := e.lookupTailCandidates(tailKey, 5)
			if len(tailWords) == 0 {
				continue
			}
			for _, head := range headWords {
				for _, tail := range tailWords {
					result.Log.TailComposedMatches++
					word := head.Word + tail.Word
					if _, ok := (*seen)[word]; ok {
						continue
					}
					(*seen)[word] = struct{}{}
					result.Items = append(result.Items, Candidate{
						Word:    word,
						Key:     head.Key + "+" + tail.Key,
						Source:  head.Source,
						Order:   head.Order,
						Pattern: head.Pattern + "+" + tail.Pattern,
						Match:   "tail_composed",
					})
					if len(result.Items) >= limit {
						return
					}
				}
			}
			if len(result.Items) > 0 {
				return
			}
		}
	}
}

func (e *Engine) appendTailFallbackMatches(result *SearchResult, seen *map[string]struct{}, query QueryPattern, limit int) {
	if len(query.RawParts) < 2 {
		return
	}
	for end := len(query.RawParts) - 1; end >= 1; end-- {
		key := strings.Join(query.RawParts[:end], "")
		if key == "" {
			continue
		}
		for _, item := range e.exact[key] {
			result.Log.TailFallbackMatches++
			if _, ok := (*seen)[item.Word]; ok {
				continue
			}
			(*seen)[item.Word] = struct{}{}
			copyItem := item
			copyItem.Match = "tail_fallback"
			result.Items = append(result.Items, copyItem)
			if len(result.Items) >= limit {
				return
			}
		}
		if len(result.Items) > 0 {
			return
		}
	}
}

func (e *Engine) lookupTailCandidates(query string, limit int) []Candidate {
	if query == "" || limit <= 0 {
		return nil
	}
	var results []Candidate
	seen := map[string]struct{}{}
	appendItems := func(items []Candidate, match string) bool {
		for _, item := range items {
			if _, ok := seen[item.Word]; ok {
				continue
			}
			seen[item.Word] = struct{}{}
			copyItem := item
			copyItem.Match = match
			results = append(results, copyItem)
			if len(results) >= limit {
				return true
			}
		}
		return false
	}

	if appendItems(e.exact[query], "tail_exact") {
		return results
	}
	if appendItems(e.initials[query], "tail_initials") {
		return results
	}
	if appendItems(e.derivedInitials[query], "tail_derived_initials") {
		return results
	}
	for _, key := range prefixKeys(e.exactKeys, query) {
		if appendItems(e.exact[key], "tail_prefix_full") {
			return results
		}
	}
	for _, key := range prefixKeys(e.initialsKeys, query) {
		if appendItems(e.initials[key], "tail_prefix_initials") {
			return results
		}
	}
	for _, key := range prefixKeys(e.derivedInitialsKeys, query) {
		if appendItems(e.derivedInitials[key], "tail_prefix_derived_initials") {
			return results
		}
	}
	return results
}

func (e *Engine) appendCombinedMatches(result *SearchResult, seen *map[string]struct{}, limit int) {
	for _, pattern := range result.Log.Segmentations {
		parts := strings.Split(pattern, "'")
		if len(parts) < 2 {
			continue
		}
		combined := e.combineByParts(parts, 0, 2)
		for _, item := range combined {
			result.Log.CombinedMatches++
			if _, ok := (*seen)[item.Word]; ok {
				continue
			}
			(*seen)[item.Word] = struct{}{}
			copyItem := item
			copyItem.Match = "combined"
			result.Items = append(result.Items, copyItem)
			if len(result.Items) >= limit {
				return
			}
		}
	}
}

func (e *Engine) combineByParts(parts []string, start int, depth int) []Candidate {
	if start >= len(parts) || depth <= 0 {
		return nil
	}

	var results []Candidate
	for end := len(parts); end > start; end-- {
		key := strings.Join(parts[start:end], "")
		words := e.exact[key]
		if len(words) == 0 {
			continue
		}
		if end == len(parts) {
			results = append(results, words...)
			if len(results) > 0 {
				return results
			}
			continue
		}
		if depth == 1 {
			continue
		}
		tails := e.combineByParts(parts, end, depth-1)
		if len(tails) == 0 {
			continue
		}
		for _, head := range words {
			for _, tail := range tails {
				results = append(results, Candidate{
					Word:    head.Word + tail.Word,
					Key:     head.Key + "+" + tail.Key,
					Source:  head.Source,
					Order:   head.Order,
					Pattern: head.Pattern + "+" + tail.Pattern,
				})
			}
		}
		if len(results) > 0 {
			return results
		}
	}
	return results
}

func (e *Engine) appendPatternMatches(result *SearchResult, seen *map[string]struct{}, query QueryPattern, limit int) {
	prefix := query.Joined
	if len(query.RawParts) > 0 {
		prefix = query.RawParts[0]
	}
	for _, key := range prefixKeys(e.exactKeys, prefix) {
		for _, item := range e.exact[key] {
			if !patternMatchesQuery(item.Pattern, query) {
				continue
			}
			result.Log.PrefixFullMatches++
			if _, ok := (*seen)[item.Word]; ok {
				continue
			}
			(*seen)[item.Word] = struct{}{}
			copyItem := item
			copyItem.Match = "pattern_prefix"
			result.Items = append(result.Items, copyItem)
			if len(result.Items) >= limit {
				return
			}
		}
	}
}

func (e *Engine) addExact(key string, candidate Candidate) {
	e.exact[key] = append(e.exact[key], candidate)
	if _, ok := e.exactKeySet[key]; ok {
		return
	}
	e.exactKeySet[key] = struct{}{}
	e.exactKeys = append(e.exactKeys, key)
}

func (e *Engine) addInitials(key string, candidate Candidate) {
	e.initials[key] = append(e.initials[key], candidate)
	if _, ok := e.initialsKeySet[key]; ok {
		return
	}
	e.initialsKeySet[key] = struct{}{}
	e.initialsKeys = append(e.initialsKeys, key)
}

func (e *Engine) addDerivedInitials(key string, candidate Candidate) {
	e.derivedInitials[key] = append(e.derivedInitials[key], candidate)
	if _, ok := e.derivedKeySet[key]; ok {
		return
	}
	e.derivedKeySet[key] = struct{}{}
	e.derivedInitialsKeys = append(e.derivedInitialsKeys, key)
}

func (e *Engine) segmentPatterns(input string) []string {
	if e.segmenter == nil {
		return nil
	}
	return e.segmenter.Patterns(input, 8)
}

func (e *Engine) bestSegmentPattern(input string) string {
	if e.segmenter == nil {
		return ""
	}
	return e.segmenter.BestPattern(input)
}

func (e *Engine) shouldPreferInitials(query string) bool {
	if len(query) < 2 || len(query) > 3 {
		return false
	}
	if len(e.initials[query]) == 0 && len(e.derivedInitials[query]) == 0 {
		return false
	}
	for _, item := range e.derivedInitials[query] {
		if len(item.Key) > len(query) {
			return true
		}
	}
	for _, item := range e.initials[query] {
		if len(item.Key) > len(query) {
			return true
		}
	}
	return false
}

func parseSogouKey(raw string) (string, string) {
	parts := strings.Split(strings.ToLower(strings.TrimSpace(raw)), "'")
	filtered := parts[:0]
	for _, part := range parts {
		part = normalizeAlphaOnly(part)
		if part == "" {
			continue
		}
		filtered = append(filtered, part)
	}
	if len(filtered) == 0 {
		return "", ""
	}
	var builder strings.Builder
	var initials strings.Builder
	for _, part := range filtered {
		builder.WriteString(part)
		initials.WriteByte(part[0])
	}
	return builder.String(), initials.String()
}

func normalizeKey(input string) string {
	return normalizeAlphaOnly(strings.ToLower(strings.TrimSpace(input)))
}

func normalizeAlphaOnly(input string) string {
	var builder strings.Builder
	builder.Grow(len(input))
	for i := 0; i < len(input); i++ {
		ch := input[i]
		if ch >= 'A' && ch <= 'Z' {
			ch = ch - 'A' + 'a'
		}
		if ch >= 'a' && ch <= 'z' {
			builder.WriteByte(ch)
		}
	}
	return builder.String()
}

func prefixKeys(keys []string, prefix string) []string {
	if prefix == "" {
		return nil
	}
	start := sort.Search(len(keys), func(i int) bool {
		return keys[i] >= prefix
	})
	if start == len(keys) {
		return nil
	}
	end := start
	for end < len(keys) && strings.HasPrefix(keys[end], prefix) {
		end++
	}
	result := make([]string, 0, end-start)
	for _, key := range keys[start:end] {
		if key == prefix {
			continue
		}
		result = append(result, key)
	}
	return result
}

func sortCandidates(items []Candidate) {
	sort.SliceStable(items, func(i, j int) bool {
		left := candidateScore(items[i])
		right := candidateScore(items[j])
		if left != right {
			return left > right
		}
		if len(items[i].Key) != len(items[j].Key) {
			return len(items[i].Key) < len(items[j].Key)
		}
		return items[i].Order < items[j].Order
	})
}

func candidateScore(item Candidate) int {
	score := 0
	switch item.Match {
	case "exact_full":
		score += 500
	case "exact_initials":
		score += 420
	case "derived_initials":
		score += 360
	case "segmented_initials":
		score += 340
	case "combined":
		score += 320
	case "tail_composed":
		score += 310
	case "tail_fallback":
		score += 300
	case "pattern_prefix":
		score += 280
	case "prefix_full":
		score += 260
	case "prefix_initials", "prefix_derived_initials":
		score += 180
	}
	if item.Pattern != "" {
		score += 20 - strings.Count(item.Pattern, "'")
	}
	if len(item.Word) >= 2 {
		score += 40
	}
	if item.Order < 256 {
		score += 32
	} else if item.Order < 4096 {
		score += 16
	}
	return score
}
