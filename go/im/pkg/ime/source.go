package ime

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Source interface {
	Name() string
	CollectSyllables(*Engine) error
	Load(*Engine) error
}

type SourceConfig struct {
	FCITXPath string
	SogouPath string
	Sources   []Source
}

type SourceSpec struct {
	Kind string
	Path string
}

type SourceFactory func(SourceSpec) (Source, error)

var sourceRegistry = map[string]SourceFactory{}

func RegisterSource(kind string, factory SourceFactory) {
	kind = strings.TrimSpace(kind)
	if kind == "" {
		panic("ime: source kind must not be empty")
	}
	if factory == nil {
		panic("ime: source factory must not be nil")
	}
	if _, exists := sourceRegistry[kind]; exists {
		panic("ime: source already registered: " + kind)
	}
	sourceRegistry[kind] = factory
}

func RegisteredSourceKinds() []string {
	kinds := make([]string, 0, len(sourceRegistry))
	for kind := range sourceRegistry {
		kinds = append(kinds, kind)
	}
	sort.Strings(kinds)
	return kinds
}

func DefaultSourceConfig(fcitxPath, sogouPath string) SourceConfig {
	sources := []Source{
		NewSogouSyllableSource(sogouPath),
		NewFCITXSource(fcitxPath),
		NewSogouSource(sogouPath),
	}
	if hasDefaultIBusDict() {
		sources = append(sources, NewIBusSource(defaultIBusPath()))
	}
	return SourceConfig{
		FCITXPath: fcitxPath,
		SogouPath: sogouPath,
		Sources:   sources,
	}
}

func ResolveDefaultDictPaths() (fcitxPath, sogouPath string) {
	if dictDir := strings.TrimSpace(os.Getenv("IM_DICT_DIR")); dictDir != "" {
		return filepath.Join(dictDir, "fcitx", "vimim.pinyin.txt"),
			filepath.Join(dictDir, "sogou", "vimim.pinyin.txt")
	}
	return filepath.Join("data", "dicts", "fcitx", "vimim.pinyin.txt"),
		filepath.Join("data", "dicts", "sogou", "vimim.pinyin.txt")
}

func DefaultSourceConfigFromEnv() SourceConfig {
	fcitxPath, sogouPath := ResolveDefaultDictPaths()
	return DefaultSourceConfig(fcitxPath, sogouPath)
}

func ParseSourceSpec(spec string) (SourceSpec, error) {
	kind, path, ok := strings.Cut(spec, "=")
	if !ok {
		return SourceSpec{}, fmt.Errorf("invalid source spec %q, expected kind=path", spec)
	}
	kind = strings.TrimSpace(kind)
	path = strings.TrimSpace(path)
	if kind == "" || path == "" {
		return SourceSpec{}, fmt.Errorf("invalid source spec %q, expected non-empty kind and path", spec)
	}
	return SourceSpec{Kind: kind, Path: path}, nil
}

func NewSource(spec SourceSpec) (Source, error) {
	factory, ok := sourceRegistry[spec.Kind]
	if !ok {
		return nil, fmt.Errorf("unsupported source kind %q", spec.Kind)
	}
	return factory(spec)
}

func BuildSourceConfig(specs []string) (SourceConfig, error) {
	sources := make([]Source, 0, len(specs))
	for _, raw := range specs {
		spec, err := ParseSourceSpec(raw)
		if err != nil {
			return SourceConfig{}, err
		}
		source, err := NewSource(spec)
		if err != nil {
			return SourceConfig{}, err
		}
		sources = append(sources, source)
	}
	return SourceConfig{Sources: sources}, nil
}
