package ime

import "fmt"

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
	if err := loadQuotedTextSource(engine, s.Path, s.Name(), true); err != nil {
		return fmt.Errorf("load sogou source: %w", err)
	}
	return nil
}
