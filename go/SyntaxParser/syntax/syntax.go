package syntax

type Buf struct {
	runes [][]rune
}

func NewBuf(sourceFile string, runes [][]rune) *Buf {
	buf := &Buf{
		runes: runes,
	}

	return buf
}

type Language string

var AllLanguages []Language

const (
	LanguagePlaintext = Language("plaintext")
	LanguageGo        = Language("go")
	LanguagePHP       = Language("php")
)

var languageToParseFunc map[Language]Func

func init() {
	languageToParseFunc = map[Language]Func{
		LanguagePlaintext: nil,
		LanguageGo:        GolangParseFunc(),
		//LanguagePHP:       languages.PHPParseFunc(),
	}

	for language := range languageToParseFunc {
		AllLanguages = append(AllLanguages, language)
	}
}

func ParserForLanguage(language Language) *P {
	//parseFunc := languageToParseFunc[language]
	parseFunc := languageToParseFunc[LanguageGo]
	if parseFunc == nil {
		return nil
	}
	return New(parseFunc)
}
