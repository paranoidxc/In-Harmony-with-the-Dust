package syntax

import (
	"unicode"
)

func GolangParseFunc() Func {
	return golangIdentifierOrKeywordParseFunc()
}

func golangIdentifierOrKeywordParseFunc() Func {
	isLetter := func(r rune) bool { return unicode.IsLetter(r) || r == '_' }
	// isLetterOrDigit := func(r rune) bool { return isLetter(r) || unicode.IsDigit(r) }
	// keywords := []string{
	// 	"break", "default", "func", "interface", "select", "case",
	// 	"defer", "go", "map", "struct", "chan", "else", "goto", "package",
	// 	"switch", "const", "fallthrough", "if", "range", "type", "continue",
	// 	"for", "import", "return", "var",
	// }
	// predeclaredIdentifiers := []string{
	// 	"bool", "byte", "complex64", "complex128", "error", "float32",
	// 	"float64", "int", "int8", "int16", "int32", "int64", "rune", "string",
	// 	"uint", "uint8", "uint16", "uint32", "uint64", "uintptr", "true",
	// 	"false", "iota", "nil", "append", "cap", "close", "complex", "copy",
	// 	"delete", "imag", "len", "make", "new", "panic", "print", "println",
	// 	"real", "recover", "any", "comparable",
	// }
	return consumeSingleRuneLike(isLetter)
	//.
	//ThenMaybe(consumeRunesLike(isLetterOrDigit)).
	//MapWithInput(recognizeKeywordOrConsume(append(keywords, predeclaredIdentifiers...)))
}
