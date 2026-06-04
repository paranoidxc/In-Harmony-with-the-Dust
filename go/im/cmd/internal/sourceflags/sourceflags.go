package sourceflags

import (
	"fmt"
	"strings"

	"im/pkg/ime"
)

type MultiValue []string

func (m *MultiValue) String() string {
	if m == nil {
		return ""
	}
	return fmt.Sprint([]string(*m))
}

func (m *MultiValue) Set(value string) error {
	*m = append(*m, value)
	return nil
}

func BuildConfig(specs []string, fcitxPath, sogouPath string) (ime.SourceConfig, error) {
	if len(specs) == 0 {
		defaultFCITX, defaultSogou := ime.ResolveDefaultDictPaths()
		if strings.TrimSpace(fcitxPath) == "" {
			fcitxPath = defaultFCITX
		}
		if strings.TrimSpace(sogouPath) == "" {
			sogouPath = defaultSogou
		}
		return ime.DefaultSourceConfig(fcitxPath, sogouPath), nil
	}
	return ime.BuildSourceConfig(specs)
}
