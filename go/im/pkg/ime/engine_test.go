package ime

import "testing"

type fakeSource struct {
	name           string
	syllables      []string
	candidates     []Candidate
	collectInvoked bool
	loadInvoked    bool
}

func (s *fakeSource) Name() string {
	return s.name
}

func (s *fakeSource) CollectSyllables(engine *Engine) error {
	s.collectInvoked = true
	for _, syllable := range s.syllables {
		engine.syllables[syllable] = struct{}{}
	}
	return nil
}

func (s *fakeSource) Load(engine *Engine) error {
	s.loadInvoked = true
	for _, candidate := range s.candidates {
		engine.addExact(candidate.Key, candidate)
	}
	return nil
}

func TestSegmentPatterns(t *testing.T) {
	segmenter := NewSegmenter(map[string]struct{}{
		"wo":   {},
		"ai":   {},
		"ni":   {},
		"xian": {},
		"xi":   {},
		"an":   {},
	})
	patterns := segmenter.Patterns("woaini", 8)
	if len(patterns) == 0 || patterns[0] != "wo'ai'ni" {
		t.Fatalf("expected wo'ai'ni, got %#v", patterns)
	}
	xian := segmenter.Patterns("xian", 8)
	if len(xian) < 2 {
		t.Fatalf("expected at least 2 segmentations for xian, got %#v", xian)
	}
}

func TestParseSogouKey(t *testing.T) {
	full, initials := parseSogouKey("ni'hao")
	if full != "nihao" || initials != "nh" {
		t.Fatalf("unexpected parse result full=%q initials=%q", full, initials)
	}
}

func TestParseQueryPatternWithRepair(t *testing.T) {
	segmenter := NewSegmenter(map[string]struct{}{
		"zai": {},
		"b":   {},
		"z":   {},
	})
	query := parseQueryPattern("zaib'z", segmenter)
	if got := query.Joined; got != "zaibz" {
		t.Fatalf("expected joined zaibz, got %q", got)
	}
	if len(query.RawParts) != 3 {
		t.Fatalf("expected 3 raw parts, got %#v", query.RawParts)
	}
	if query.RawParts[0] != "zai" || query.RawParts[1] != "b" || query.RawParts[2] != "z" {
		t.Fatalf("expected zai,b,z got %#v", query.RawParts)
	}
}

func TestAutoQueryPatternRepair(t *testing.T) {
	segmenter := NewSegmenter(map[string]struct{}{
		"zai": {},
		"b":   {},
		"z":   {},
	})
	query := autoQueryPattern("zaibz", segmenter)
	if !query.HasBoundary {
		t.Fatalf("expected auto query to have boundary, got %#v", query)
	}
	if len(query.RawParts) != 3 || query.RawParts[0] != "zai" || query.RawParts[1] != "b" || query.RawParts[2] != "z" {
		t.Fatalf("expected zai,b,z got %#v", query.RawParts)
	}
}

func TestSearchOrder(t *testing.T) {
	engine := NewEngine()
	engine.addExact("nihao", Candidate{Word: "你好", Key: "nihao", Source: "fcitx"})
	engine.addInitials("nh", Candidate{Word: "你好", Key: "nihao", Source: "sogou"})
	engine.addExact("nihaoma", Candidate{Word: "你好吗", Key: "nihaoma", Source: "fcitx"})
	engine.addInitials("nhm", Candidate{Word: "你好吗", Key: "nihaoma", Source: "sogou"})
	engine.addExact("ni", Candidate{Word: "你", Key: "ni", Source: "fcitx"})
	engine.exactKeys = []string{"ni", "nihao", "nihaoma"}
	engine.initialsKeys = []string{"nh", "nhm"}

	result := engine.Search("nihao", 10)
	if len(result.Items) < 2 {
		t.Fatalf("expected at least 2 items, got %d", len(result.Items))
	}
	if result.Items[0].Word != "你好" {
		t.Fatalf("expected first item to be 你好, got %q", result.Items[0].Word)
	}
	if result.Items[1].Word != "你好吗" {
		t.Fatalf("expected second item to be 你好吗, got %q", result.Items[1].Word)
	}
}

func TestInitialsSearch(t *testing.T) {
	engine := NewEngine()
	engine.addExact("beijing", Candidate{Word: "北京", Key: "beijing", Source: "fcitx"})
	engine.addInitials("bj", Candidate{Word: "北京", Key: "beijing", Source: "sogou"})
	engine.exactKeys = []string{"beijing"}
	engine.initialsKeys = []string{"bj"}

	result := engine.Search("bj", 10)
	if len(result.Items) != 1 || result.Items[0].Word != "北京" {
		t.Fatalf("expected 北京, got %#v", result.Items)
	}
}

func TestSegmentedInitialsSearch(t *testing.T) {
	engine := NewEngine()
	engine.syllables = map[string]struct{}{
		"wo": {},
		"ai": {},
		"ni": {},
	}
	engine.segmenter = NewSegmenter(engine.syllables)
	engine.addExact("woaini", Candidate{Word: "我爱你", Key: "woaini", Source: "fcitx"})
	engine.addDerivedInitials("wan", Candidate{Word: "我爱你", Key: "woaini", Source: "fcitx"})
	engine.exactKeys = []string{"woaini"}
	engine.derivedInitialsKeys = []string{"wan"}

	result := engine.Search("woaini", 10)
	if len(result.Log.Segmentations) == 0 || result.Log.Segmentations[0] != "wo'ai'ni" {
		t.Fatalf("expected segmentation wo'ai'ni, got %#v", result.Log.Segmentations)
	}

	initialsResult := engine.Search("wan", 10)
	if len(initialsResult.Items) != 1 || initialsResult.Items[0].Word != "我爱你" {
		t.Fatalf("expected 我爱你 for wan, got %#v", initialsResult.Items)
	}
}

func TestPatternQueryMatch(t *testing.T) {
	query := parseQueryPattern("zhen's'd", nil)
	if !query.HasBoundary || query.Joined != "zhensd" {
		t.Fatalf("unexpected query pattern %#v", query)
	}
	if !patternMatchesQuery("zhen'shi'de", query) {
		t.Fatal("expected pattern to match zhen'shi'de")
	}
	if !patternMatchesQuery("zhen'shi'di", query) {
		t.Fatal("expected pattern to match zhen'shi'di")
	}
	if patternMatchesQuery("zhen'de", query) {
		t.Fatal("did not expect pattern to match zhen'de")
	}
}

func TestStructuredBoundarySearch(t *testing.T) {
	engine := NewEngine()
	engine.syllables = map[string]struct{}{
		"zhen": {},
		"shi":  {},
		"de":   {},
		"di":   {},
	}
	engine.segmenter = NewSegmenter(engine.syllables)
	engine.addExact("zhenshide", Candidate{Word: "真实的", Key: "zhenshide", Source: "fcitx", Order: 0, Pattern: "zhen'shi'de"})
	engine.addExact("zhenshidi", Candidate{Word: "真实地", Key: "zhenshidi", Source: "fcitx", Order: 1, Pattern: "zhen'shi'di"})
	engine.addExact("zhende", Candidate{Word: "真的", Key: "zhende", Source: "fcitx", Order: 2, Pattern: "zhen'de"})
	engine.exactKeys = []string{"zhende", "zhenshide", "zhenshidi"}

	result := engine.Search("zhen's'd", 10)
	if len(result.Items) != 2 {
		t.Fatalf("expected 2 items, got %#v", result.Items)
	}
	if result.Items[0].Word != "真实的" || result.Items[1].Word != "真实地" {
		t.Fatalf("unexpected results %#v", result.Items)
	}
}

func TestStructuredBoundaryRepairSearch(t *testing.T) {
	engine := NewEngine()
	engine.syllables = map[string]struct{}{
		"zai":  {},
		"bian": {},
		"zhe":  {},
		"b":    {},
		"z":    {},
	}
	engine.segmenter = NewSegmenter(engine.syllables)
	engine.addExact("zaibianzhe", Candidate{Word: "在编者", Key: "zaibianzhe", Source: "fcitx", Order: 0, Pattern: "zai'bian'zhe"})
	engine.exactKeys = []string{"zaibianzhe"}

	result := engine.Search("zaib'z", 10)
	if len(result.Items) != 1 || result.Items[0].Word != "在编者" {
		t.Fatalf("expected 在编者, got %#v", result.Items)
	}
	if len(result.Log.Segmentations) == 0 || result.Log.Segmentations[0] != "zai'b'z" {
		t.Fatalf("expected repaired segmentation zai'b'z, got %#v", result.Log.Segmentations)
	}
}

func TestAutoStructuredRepairSearch(t *testing.T) {
	engine := NewEngine()
	engine.syllables = map[string]struct{}{
		"zai":  {},
		"bian": {},
		"zhe":  {},
		"b":    {},
		"z":    {},
	}
	engine.segmenter = NewSegmenter(engine.syllables)
	engine.addExact("zaibianzhe", Candidate{Word: "在编者", Key: "zaibianzhe", Source: "fcitx", Order: 0, Pattern: "zai'bian'zhe"})
	engine.exactKeys = []string{"zaibianzhe"}

	result := engine.Search("zaibz", 10)
	if len(result.Items) != 1 || result.Items[0].Word != "在编者" {
		t.Fatalf("expected 在编者, got %#v", result.Items)
	}
	if len(result.Log.Segmentations) == 0 || result.Log.Segmentations[0] != "zai'b'z" {
		t.Fatalf("expected repaired segmentation zai'b'z, got %#v", result.Log.Segmentations)
	}
}

func TestCombinedMatches(t *testing.T) {
	engine := NewEngine()
	engine.syllables = map[string]struct{}{
		"wo":   {},
		"ai":   {},
		"bei":  {},
		"jing": {},
	}
	engine.segmenter = NewSegmenter(engine.syllables)
	engine.addExact("woai", Candidate{Word: "我爱", Key: "woai", Source: "fcitx", Order: 0, Pattern: "wo'ai"})
	engine.addExact("beijing", Candidate{Word: "北京", Key: "beijing", Source: "fcitx", Order: 1, Pattern: "bei'jing"})
	engine.exactKeys = []string{"beijing", "woai"}

	result := engine.Search("woaibeijing", 10)
	found := false
	for _, item := range result.Items {
		if item.Word == "我爱北京" && item.Match == "combined" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected combined 我爱北京, got %#v", result.Items)
	}
}

func TestTailFallbackMatches(t *testing.T) {
	engine := NewEngine()
	engine.syllables = map[string]struct{}{
		"hen": {},
		"bu":  {},
		"cuo": {},
		"d":   {},
	}
	engine.segmenter = NewSegmenter(engine.syllables)
	engine.addExact("henbucuo", Candidate{
		Word:    "很不错",
		Key:     "henbucuo",
		Source:  "fcitx",
		Order:   0,
		Pattern: "hen'bu'cuo",
	})
	engine.exactKeys = []string{"henbucuo"}

	result := engine.Search("henbucuod", 5)
	if len(result.Items) == 0 || result.Items[0].Word != "很不错" {
		t.Fatalf("expected 很不错 fallback, got %#v", result.Items)
	}
	if result.Items[0].Match != "tail_fallback" {
		t.Fatalf("expected tail_fallback match, got %#v", result.Items[0])
	}
}

func TestTailComposedMatches(t *testing.T) {
	engine := NewEngine()
	engine.syllables = map[string]struct{}{
		"hen": {},
		"bu":  {},
		"cuo": {},
		"d":   {},
		"de":  {},
		"di":  {},
	}
	engine.segmenter = NewSegmenter(engine.syllables)
	engine.addExact("henbucuo", Candidate{
		Word:    "很不错",
		Key:     "henbucuo",
		Source:  "fcitx",
		Order:   0,
		Pattern: "hen'bu'cuo",
	})
	engine.addExact("de", Candidate{
		Word:    "的",
		Key:     "de",
		Source:  "fcitx",
		Order:   0,
		Pattern: "de",
	})
	engine.addExact("di", Candidate{
		Word:    "地",
		Key:     "di",
		Source:  "fcitx",
		Order:   1,
		Pattern: "di",
	})
	engine.exactKeys = []string{"de", "di", "henbucuo"}

	result := engine.Search("henbucuod", 5)
	if len(result.Items) == 0 || result.Items[0].Word != "很不错的" {
		t.Fatalf("expected 很不错的, got %#v", result.Items)
	}
	if result.Items[0].Match != "tail_composed" {
		t.Fatalf("expected tail_composed match, got %#v", result.Items[0])
	}
}

func TestSessionCommit(t *testing.T) {
	engine := NewEngine()
	engine.addExact("nihao", Candidate{Word: "你好", Key: "nihao", Source: "fcitx", Order: 0, Pattern: "ni'hao"})
	engine.exactKeys = []string{"nihao"}

	session := NewSession(engine, 5)
	for _, r := range "nihao" {
		session.InputRune(r)
	}
	if len(session.Candidates()) != 1 {
		t.Fatalf("expected 1 candidate, got %#v", session.Candidates())
	}
	if committed := session.CommitSelection(); committed != "你好" {
		t.Fatalf("expected 你好, got %q", committed)
	}
	if session.Buffer() != "" {
		t.Fatalf("expected empty buffer, got %q", session.Buffer())
	}
	if session.CommittedText() != "你好" {
		t.Fatalf("expected committed text 你好, got %q", session.CommittedText())
	}
}

func TestIBusQuotedFormatSearch(t *testing.T) {
	engine := NewEngine()
	engine.syllables = map[string]struct{}{
		"hen": {},
		"bu":  {},
		"cuo": {},
	}
	engine.segmenter = NewSegmenter(engine.syllables)
	engine.addExact("henbucuo", Candidate{
		Word:    "很不错",
		Key:     "henbucuo",
		Source:  "ibus",
		Order:   0,
		Pattern: "hen'bu'cuo",
	})
	engine.exactKeys = []string{"henbucuo"}

	result := engine.Search("henbucuo", 5)
	if len(result.Items) != 1 || result.Items[0].Word != "很不错" {
		t.Fatalf("expected 很不错, got %#v", result.Items)
	}
}

func TestLoadEngineWithConfig(t *testing.T) {
	source := &fakeSource{
		name:      "fake",
		syllables: []string{"ni", "hao"},
		candidates: []Candidate{
			{Word: "你好", Key: "nihao", Source: "fake", Pattern: "ni'hao"},
		},
	}
	engine, err := LoadEngineWithConfig(SourceConfig{
		Sources: []Source{source},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !source.collectInvoked || !source.loadInvoked {
		t.Fatalf("expected source hooks to be invoked, got collect=%v load=%v", source.collectInvoked, source.loadInvoked)
	}
	result := engine.Search("nihao", 5)
	if len(result.Items) != 1 || result.Items[0].Word != "你好" {
		t.Fatalf("expected fake source result, got %#v", result.Items)
	}
}

func TestParseSourceSpec(t *testing.T) {
	spec, err := ParseSourceSpec("fcitx=data/dicts/fcitx/vimim.pinyin.txt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if spec.Kind != "fcitx" || spec.Path != "data/dicts/fcitx/vimim.pinyin.txt" {
		t.Fatalf("unexpected spec %#v", spec)
	}
}

func TestBuildSourceConfig(t *testing.T) {
	config, err := BuildSourceConfig([]string{
		"sogou-syllables=data/dicts/sogou/vimim.pinyin.txt",
		"fcitx=data/dicts/fcitx/vimim.pinyin.txt",
		"ibus=data/dicts/ibus/vimim.pinyin.txt",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(config.Sources) != 3 {
		t.Fatalf("expected 3 sources, got %d", len(config.Sources))
	}
	if config.Sources[0].Name() != "sogou-syllables" || config.Sources[1].Name() != "fcitx" || config.Sources[2].Name() != "ibus" {
		t.Fatalf("unexpected source order: %s, %s, %s", config.Sources[0].Name(), config.Sources[1].Name(), config.Sources[2].Name())
	}
}

func TestResolveDefaultDictPaths(t *testing.T) {
	t.Setenv("IM_DICT_DIR", "")
	fcitxPath, sogouPath := ResolveDefaultDictPaths()
	if fcitxPath != "data/dicts/fcitx/vimim.pinyin.txt" {
		t.Fatalf("unexpected default fcitx path %q", fcitxPath)
	}
	if sogouPath != "data/dicts/sogou/vimim.pinyin.txt" {
		t.Fatalf("unexpected default sogou path %q", sogouPath)
	}
}

func TestResolveDefaultDictPathsFromEnv(t *testing.T) {
	t.Setenv("IM_DICT_DIR", "/tmp/im-dicts")
	fcitxPath, sogouPath := ResolveDefaultDictPaths()
	if fcitxPath != "/tmp/im-dicts/fcitx/vimim.pinyin.txt" {
		t.Fatalf("unexpected env fcitx path %q", fcitxPath)
	}
	if sogouPath != "/tmp/im-dicts/sogou/vimim.pinyin.txt" {
		t.Fatalf("unexpected env sogou path %q", sogouPath)
	}
}

func TestRegisteredSourceKinds(t *testing.T) {
	kinds := RegisteredSourceKinds()
	if len(kinds) < 3 {
		t.Fatalf("expected builtin source kinds, got %#v", kinds)
	}
	foundIBus := false
	for _, kind := range kinds {
		if kind == "ibus" {
			foundIBus = true
			break
		}
	}
	if !foundIBus {
		t.Fatalf("expected builtin source kinds, got %#v", kinds)
	}
}

func TestNewSourceFromRegistry(t *testing.T) {
	const kind = "test-factory"
	RegisterSource(kind, func(spec SourceSpec) (Source, error) {
		return &fakeSource{name: spec.Kind}, nil
	})

	source, err := NewSource(SourceSpec{Kind: kind, Path: "ignored"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if source.Name() != kind {
		t.Fatalf("expected source name %q, got %q", kind, source.Name())
	}
}
