package syntax

import (
	"unicode"
)

func GolangParseFunc() Func {
	return golangLineCommentParseFunc().
		Or(golangGeneralCommentParseFunc()).
		Or(golangIdentifierOrKeywordParseFunc()).
		Or(golangOperatorParseFunc()).
		Or(golangRuneLiteralParseFunc()).
		Or(golangRawStringLiteralParseFunc()).
		Or(golangInterpretedStringLiteralParseFunc()).
		Or(golangFloatLiteralParseFunc()).
		Or(golangIntegerLiteralParseFunc())
}

func golangLineCommentParseFunc() Func {
	return consumeString("//").
		ThenMaybe(consumeToNextLineFeed).
		Map(recognizeToken(TokenRoleComment))
}

func golangGeneralCommentParseFunc() Func {
	return consumeString("/*").
		Then(consumeToString("*/")).
		Map(recognizeToken(TokenRoleComment))
}

func golangIdentifierOrKeywordParseFunc() Func {
	isLetter := func(r rune) bool { return unicode.IsLetter(r) || r == '_' }
	isLetterOrDigit := func(r rune) bool { return isLetter(r) || unicode.IsDigit(r) }
	keywords := []string{
		"break", "default", "func", "interface", "select", "case",
		"defer", "go", "map", "struct", "chan", "else", "goto", "package",
		"switch", "const", "fallthrough", "if", "range", "type", "continue",
		"for", "import", "return", "var",
	}
	predeclaredIdentifiers := []string{
		"bool", "byte", "complex64", "complex128", "error", "float32",
		"float64", "int", "int8", "int16", "int32", "int64", "rune", "string",
		"uint", "uint8", "uint16", "uint32", "uint64", "uintptr", "true",
		"false", "iota", "nil", "append", "cap", "close", "complex", "copy",
		"delete", "imag", "len", "make", "new", "panic", "print", "println",
		"real", "recover", "any", "comparable",
	}
	return consumeSingleRuneLike(isLetter).
		ThenMaybe(consumeRunesLike(isLetterOrDigit)).
		MapWithInput(recognizeKeywordOrConsume(append(keywords, predeclaredIdentifiers...)))
}

func golangOperatorParseFunc() Func {
	return consumeLongestMatchingOption([]string{
		"+", "&", "+=", "&=", "&&", "==", "!=",
		"-", "|", "-=", "|=", "||", "<", "<=",
		"*", "^", "*=", "^=", "<-", ">", ">=",
		"/", "<<", "/=", "<<=", "++", "=", ":=",
		"%", ">>", "%=", ">>=", "--", "!",
		"&^", "&^=", "~",
	}).Map(recognizeToken(TokenRoleOperator))
}

func golangRuneLiteralParseFunc() Func {
	return parseCStyleString('\'', false)
}

func golangRawStringLiteralParseFunc() Func {
	return consumeString("`").
		Then(consumeToString("`")).
		Map(recognizeToken(TokenRoleString))
}

func golangInterpretedStringLiteralParseFunc() Func {
	return parseCStyleString('"', false)
}

func golangFloatLiteralParseFunc() Func {
	consumeDecimalDigits := consumeDigitsAndSeparators(false, func(r rune) bool {
		return r >= '0' && r <= '9'
	})
	consumeDecimalExponent := consumeSingleRuneLike(func(r rune) bool {
		return r == 'e' || r == 'E'
	}).ThenMaybe(consumeSingleRuneLike(func(r rune) bool {
		return r == '+' || r == '-'
	})).Then(consumeDecimalDigits)

	consumeDecimalFloatLiteralFormA := consumeDecimalDigits.
		Then(consumeString(".")).
		ThenMaybe(consumeDecimalDigits).
		ThenMaybe(consumeDecimalExponent)

	consumeDecimalFloatLiteralFormB := consumeDecimalDigits.
		Then(consumeDecimalExponent)

	consumeDecimalFloatLiteralFormC := consumeString(".").
		Then(consumeDecimalDigits).
		ThenMaybe(consumeDecimalExponent)

	consumeDecimalFloatLiteral := consumeDecimalFloatLiteralFormA.
		Or(consumeDecimalFloatLiteralFormB).
		Or(consumeDecimalFloatLiteralFormC)

	consumeHexDigitsAllowLeadingUnderscore := consumeDigitsAndSeparators(true, func(r rune) bool {
		return (r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')
	})
	consumeHexDigits := consumeDigitsAndSeparators(false, func(r rune) bool {
		return (r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')
	})
	consumeHexExponent := consumeSingleRuneLike(func(r rune) bool {
		return r == 'p' || r == 'P'
	}).ThenMaybe(consumeSingleRuneLike(func(r rune) bool {
		return r == '+' || r == '-'
	})).Then(consumeDecimalDigits)

	consumeHexMantissaFormA := consumeHexDigitsAllowLeadingUnderscore.
		Then(consumeString(".")).
		ThenMaybe(consumeHexDigits)

	consumeHexMantissaFormB := consumeHexDigitsAllowLeadingUnderscore

	consumeHexMantissaFormC := consumeString(".").Then(consumeHexDigits)

	consumeHexMantissa := consumeHexMantissaFormA.
		Or(consumeHexMantissaFormB).
		Or(consumeHexMantissaFormC)

	consumeHexFloatLiteral := consumeString("0").
		Then(consumeSingleRuneLike(func(r rune) bool { return r == 'x' || r == 'X' })).
		Then(consumeHexMantissa).
		Then(consumeHexExponent)

	return consumeHexFloatLiteral.
		Or(consumeDecimalFloatLiteral).
		ThenMaybe(consumeString("i")).
		Map(recognizeToken(TokenRoleNumber))
}

func golangIntegerLiteralParseFunc() Func {
	consumeDecimalLiteral := consumeString("0").
		Or(consumeSingleRuneLike(func(r rune) bool { return r >= '1' && r <= '9' }).
			ThenMaybe(consumeDigitsAndSeparators(true, func(r rune) bool { return r >= '0' && r <= '9' })))

	consumeBinaryLiteral := consumeString("0").
		Then(consumeSingleRuneLike(func(r rune) bool { return r == 'b' || r == 'B' })).
		Then(consumeDigitsAndSeparators(true, func(r rune) bool { return r == '0' || r == '1' }))

	consumeOctalLiteral := consumeString("0").
		ThenMaybe(consumeSingleRuneLike(func(r rune) bool { return r == 'o' || r == 'O' })).
		Then(consumeDigitsAndSeparators(true, func(r rune) bool { return r >= '0' && r <= '7' }))

	consumeHexLiteral := consumeString("0").
		Then(consumeSingleRuneLike(func(r rune) bool { return r == 'x' || r == 'X' })).
		Then(consumeDigitsAndSeparators(true, func(r rune) bool {
			return (r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')
		}))

	return consumeBinaryLiteral.
		Or(consumeOctalLiteral).
		Or(consumeHexLiteral).
		Or(consumeDecimalLiteral).
		ThenMaybe(consumeString("i")).
		Map(recognizeToken(TokenRoleNumber))
}
