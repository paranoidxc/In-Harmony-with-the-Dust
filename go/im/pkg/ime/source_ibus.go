package ime

import (
	"fmt"
	"os"
)

type IBusSource struct {
	Path string
}

func init() {
	RegisterSource("ibus", func(spec SourceSpec) (Source, error) {
		return NewIBusSource(spec.Path), nil
	})
}

func NewIBusSource(path string) *IBusSource {
	return &IBusSource{Path: path}
}

func (s *IBusSource) Name() string {
	return "ibus"
}

func (s *IBusSource) CollectSyllables(engine *Engine) error {
	if err := collectQuotedTextSyllables(engine, s.Path, s.Name()); err != nil {
		return fmt.Errorf("collect ibus syllables: %w", err)
	}
	return nil
}

func (s *IBusSource) Load(engine *Engine) error {
	if err := loadQuotedTextSource(engine, s.Path, s.Name(), false); err != nil {
		return fmt.Errorf("load ibus source: %w", err)
	}
	return nil
}

func defaultIBusPath() string {
	return "data/dicts/ibus/vimim.pinyin.txt"
}

func hasDefaultIBusDict() bool {
	_, err := os.Stat(defaultIBusPath())
	return err == nil
}
