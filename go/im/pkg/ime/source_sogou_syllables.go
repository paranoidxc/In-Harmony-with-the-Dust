package ime

import "fmt"

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
	if err := collectQuotedTextSyllables(engine, s.Path, s.Name()); err != nil {
		return fmt.Errorf("collect sogou syllables: %w", err)
	}
	return nil
}

func (s *SogouSyllableSource) Load(*Engine) error {
	return nil
}
