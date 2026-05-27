package syntax

import (
	"testing"
)

func makeBuf(s string) *Buf {
	runes := [][]rune{}
	line := []rune{}
	for _, ch := range s {
		if ch == '\n' {
			runes = append(runes, line)
			line = []rune{}
		} else {
			line = append(line, ch)
		}
	}
	runes = append(runes, line)
	return NewBuf("test", runes)
}

func checkTokens(t *testing.T, got []ComputedToken, want []ComputedToken) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("token count: got %d, want %d\n  got:  %v\n  want: %v", len(got), len(want), got, want)
		return
	}
	for i, g := range got {
		w := want[i]
		if g.Offset != w.Offset || g.Length != w.Length || g.Role != w.Role {
			t.Errorf("token[%d]: got {%d, %d, %d}, want {%d, %d, %d}",
				i, g.Offset, g.Length, g.Role, w.Offset, w.Length, w.Role)
		}
	}
}

func TestParseAllEmpty(t *testing.T) {
	buf := makeBuf("")
	p := New(GolangParseFunc())
	p.ParseAll(buf)

	tokens := p.Tokens()
	if len(tokens) != 0 {
		t.Errorf("expected no tokens for empty doc, got %v", tokens)
	}
}

func TestParseAllSingleKeyword(t *testing.T) {
	buf := makeBuf("func")
	p := New(GolangParseFunc())
	p.ParseAll(buf)

	tokens := p.Tokens()
	checkTokens(t, tokens, []ComputedToken{
		{Offset: 0, Length: 4, Role: TokenRoleKeyword},
	})
}

func TestParseAllKeywordAndOperator(t *testing.T) {
	buf := makeBuf("func()")
	p := New(GolangParseFunc())
	p.ParseAll(buf)

	tokens := p.Tokens()
	checkTokens(t, tokens, []ComputedToken{
		{Offset: 0, Length: 4, Role: TokenRoleKeyword},
		{Offset: 4, Length: 1, Role: TokenRoleOperator},
		{Offset: 5, Length: 1, Role: TokenRoleOperator},
	})
}

func TestParseAllLineComment(t *testing.T) {
	buf := makeBuf("// hello\nfunc")
	p := New(GolangParseFunc())
	p.ParseAll(buf)

	tokens := p.Tokens()
	checkTokens(t, tokens, []ComputedToken{
		{Offset: 0, Length: 8, Role: TokenRoleComment},
		{Offset: 9, Length: 4, Role: TokenRoleKeyword},
	})
}

func TestParseAllBlockComment(t *testing.T) {
	buf := makeBuf("/* x */")
	p := New(GolangParseFunc())
	p.ParseAll(buf)

	tokens := p.Tokens()
	checkTokens(t, tokens, []ComputedToken{
		{Offset: 0, Length: 7, Role: TokenRoleComment},
	})
}

func TestParseAllStringLiteral(t *testing.T) {
	buf := makeBuf(`"hello"`)
	p := New(GolangParseFunc())
	p.ParseAll(buf)

	tokens := p.Tokens()
	checkTokens(t, tokens, []ComputedToken{
		{Offset: 0, Length: 7, Role: TokenRoleString},
	})
}

func TestParseAllRawStringLiteral(t *testing.T) {
	buf := makeBuf("`hello`")
	p := New(GolangParseFunc())
	p.ParseAll(buf)

	tokens := p.Tokens()
	checkTokens(t, tokens, []ComputedToken{
		{Offset: 0, Length: 7, Role: TokenRoleString},
	})
}

func TestParseAllNumberLiteral(t *testing.T) {
	buf := makeBuf("42")
	p := New(GolangParseFunc())
	p.ParseAll(buf)

	tokens := p.Tokens()
	checkTokens(t, tokens, []ComputedToken{
		{Offset: 0, Length: 2, Role: TokenRoleNumber},
	})
}

func TestParseAllMixedCode(t *testing.T) {
	buf := makeBuf("// comment\npackage main\nfunc main() {}\n")
	p := New(GolangParseFunc())
	p.ParseAll(buf)

	tokens := p.Tokens()
	// "// comment" = 10 chars, "\n" = 1 char, "package" starts at offset 11
	// "package" = 7 chars, "main" starts at offset 19
	// the next "\n" is at offset 23, so "func" starts at offset 24
	checkTokens(t, tokens, []ComputedToken{
		{Offset: 0, Length: 10, Role: TokenRoleComment},
		{Offset: 11, Length: 7, Role: TokenRoleKeyword},
		{Offset: 24, Length: 4, Role: TokenRoleKeyword},
		{Offset: 33, Length: 1, Role: TokenRoleOperator}, // (
		{Offset: 34, Length: 1, Role: TokenRoleOperator}, // )
		{Offset: 36, Length: 1, Role: TokenRoleOperator}, // {
		{Offset: 37, Length: 1, Role: TokenRoleOperator}, // }
	})
}

func TestParseAllPredeclaredIdentifiers(t *testing.T) {
	buf := makeBuf("nil true")
	p := New(GolangParseFunc())
	p.ParseAll(buf)

	tokens := p.Tokens()
	checkTokens(t, tokens, []ComputedToken{
		{Offset: 0, Length: 3, Role: TokenRoleKeyword},
		{Offset: 4, Length: 4, Role: TokenRoleKeyword},
	})
}

func TestParseAllMultipleStrings(t *testing.T) {
	buf := makeBuf(`"foo" "bar"`)
	p := New(GolangParseFunc())
	p.ParseAll(buf)

	tokens := p.Tokens()
	checkTokens(t, tokens, []ComputedToken{
		{Offset: 0, Length: 5, Role: TokenRoleString},
		{Offset: 6, Length: 5, Role: TokenRoleString},
	})
}

func TestParseAllOperators(t *testing.T) {
	buf := makeBuf("x := y + 1")
	p := New(GolangParseFunc())
	p.ParseAll(buf)

	tokens := p.Tokens()
	checkTokens(t, tokens, []ComputedToken{
		{Offset: 2, Length: 2, Role: TokenRoleOperator}, // :=
		{Offset: 7, Length: 1, Role: TokenRoleOperator}, // +
		{Offset: 9, Length: 1, Role: TokenRoleNumber},   // 1
	})
}

func TestParseAllComputationMerging(t *testing.T) {
	// Build a large enough document so that computation merging actually happens
	// (small segments get merged when consumedLength < minInitialConsumedLen)
	line := "package main\n"
	var s string
	for i := 0; i < 100; i++ {
		s += line
	}
	buf := makeBuf(s)
	p := New(GolangParseFunc())
	p.ParseAll(buf)

	// Should have far fewer computations than 100 due to merging
	if len(p.computations) >= 100 {
		t.Errorf("expected computations to be merged (< 100), got %d", len(p.computations))
	}

	tokens := p.Tokens()
	// Each "package" keyword should produce a token
	keywordCount := 0
	for _, tok := range tokens {
		if tok.Role == TokenRoleKeyword {
			keywordCount++
		}
	}
	if keywordCount != 100 {
		t.Errorf("expected 100 keyword tokens, got %d", keywordCount)
	}
}

func TestParseAllComputationStateContinuity(t *testing.T) {
	buf := makeBuf("// comment\npackage main\n")
	p := New(GolangParseFunc())
	p.ParseAll(buf)

	// Each computation's startState should equal the previous computation's endState
	for i := 1; i < len(p.computations); i++ {
		prev := p.computations[i-1]
		curr := p.computations[i]
		if prev.endState != curr.startState {
			t.Errorf("state discontinuity at computation[%d]: prev.endState=%v, curr.startState=%v",
				i, prev.endState, curr.startState)
		}
	}
}

func TestParseAllTotalConsumedEqualsDocLength(t *testing.T) {
	s := "func main() { println(\"hi\") }\n"
	buf := makeBuf(s)
	p := New(GolangParseFunc())
	p.ParseAll(buf)

	var totalConsumed uint64
	for _, c := range p.computations {
		totalConsumed += c.consumedLength
	}
	if totalConsumed != uint64(len(s)) {
		t.Errorf("total consumed length %d != doc length %d", totalConsumed, len(s))
	}
}

func TestParseAfterEditAppend(t *testing.T) {
	line := "package main\n"
	var s string
	for i := 0; i < 100; i++ {
		s += line
	}
	buf := makeBuf(s)
	p := New(GolangParseFunc())
	p.ParseAll(buf)
	beforeFirst := p.computations[0]

	newBuf := makeBuf(s + "func main() {}")
	p.ParseAfterEdit(newBuf, Edit{Offset: uint64(len(s)), NumInserted: uint64(len("func main() {}"))})

	if len(p.computations) == 0 || p.computations[0].consumedLength != beforeFirst.consumedLength {
		t.Errorf("expected append to reuse first computation")
	}

	tokens := p.Tokens()
	if len(tokens) == 0 {
		t.Fatalf("expected tokens after append")
	}
	last := tokens[len(tokens)-5:]
	checkTokens(t, last, []ComputedToken{
		{Offset: uint64(len(s)), Length: 4, Role: TokenRoleKeyword},
		{Offset: uint64(len(s) + len("func main")), Length: 1, Role: TokenRoleOperator},
		{Offset: uint64(len(s) + len("func main(")), Length: 1, Role: TokenRoleOperator},
		{Offset: uint64(len(s) + len("func main() ")), Length: 1, Role: TokenRoleOperator},
		{Offset: uint64(len(s) + len("func main() {")), Length: 1, Role: TokenRoleOperator},
	})
}

func TestParseAfterEditInsert(t *testing.T) {
	buf := makeBuf("package main\nfunc main() {}")
	p := New(GolangParseFunc())
	p.ParseAll(buf)

	newBuf := makeBuf("package main\nfunc test() {}")
	p.ParseAfterEdit(newBuf, Edit{Offset: uint64(len("package main\nfunc ")), NumInserted: uint64(len("test")), NumDeleted: uint64(len("main"))})

	checkTokens(t, p.Tokens(), []ComputedToken{
		{Offset: 0, Length: 7, Role: TokenRoleKeyword},
		{Offset: 13, Length: 4, Role: TokenRoleKeyword},
		{Offset: 22, Length: 1, Role: TokenRoleOperator},
		{Offset: 23, Length: 1, Role: TokenRoleOperator},
		{Offset: 25, Length: 1, Role: TokenRoleOperator},
		{Offset: 26, Length: 1, Role: TokenRoleOperator},
	})
}

func TestParseAfterEditDelete(t *testing.T) {
	buf := makeBuf("package main\nfunc main() {}")
	p := New(GolangParseFunc())
	p.ParseAll(buf)

	newBuf := makeBuf("package main\n")
	p.ParseAfterEdit(newBuf, Edit{Offset: uint64(len("package main\n")), NumDeleted: uint64(len("func main() {}"))})

	checkTokens(t, p.Tokens(), []ComputedToken{
		{Offset: 0, Length: 7, Role: TokenRoleKeyword},
	})
}

func TestParseAfterEditAddDocument(t *testing.T) {
	buf := makeBuf("")
	p := New(GolangParseFunc())
	p.ParseAll(buf)

	newBuf := makeBuf("func main() {}")
	p.ParseAfterEdit(newBuf, Edit{Offset: 0, NumInserted: uint64(len("func main() {}"))})

	checkTokens(t, p.Tokens(), []ComputedToken{
		{Offset: 0, Length: 4, Role: TokenRoleKeyword},
		{Offset: 9, Length: 1, Role: TokenRoleOperator},
		{Offset: 10, Length: 1, Role: TokenRoleOperator},
		{Offset: 12, Length: 1, Role: TokenRoleOperator},
		{Offset: 13, Length: 1, Role: TokenRoleOperator},
	})
}

func parseTokensFromAll(s string) []ComputedToken {
	p := New(GolangParseFunc())
	p.ParseAll(makeBuf(s))
	return p.Tokens()
}

func makeLargeGoSource(numFuncs int) string {
	var s string
	s += "package main\n\n"
	s += "const banner = `large\nsource`\n\n"
	for i := 0; i < numFuncs; i++ {
		s += "func generated() {\n"
		s += "\t// repeated parser test line\n"
		s += "\tvalue := 42 + 7\n"
		s += "\tprintln(\"value\", value)\n"
		s += "}\n\n"
	}
	return s
}

func checkSourceHasAtLeastLines(t *testing.T, s string, want int) {
	t.Helper()
	lines := 0
	for _, r := range s {
		if r == '\n' {
			lines++
		}
	}
	if lines < want {
		t.Fatalf("source has %d lines, want at least %d", lines, want)
	}
}

func checkParseAfterEditMatchesParseAll(t *testing.T, oldText string, newText string, edit Edit) {
	t.Helper()
	p := New(GolangParseFunc())
	p.ParseAll(makeBuf(oldText))
	p.ParseAfterEdit(makeBuf(newText), edit)
	checkTokens(t, p.Tokens(), parseTokensFromAll(newText))
}

func TestParseAfterEditMatchesParseAllForKeywordInsertion(t *testing.T) {
	oldText := "package main\n\nfunc main() {\n\tprintln(\"hi\")\n}\n"
	newText := "package main\n\nfunc main() {\n\treturn\n\tprintln(\"hi\")\n}\n"
	inserted := "return\n\t"
	checkParseAfterEditMatchesParseAll(t, oldText, newText, Edit{
		Offset:      uint64(len("package main\n\nfunc main() {\n\t")),
		NumInserted: uint64(len(inserted)),
	})
}

func TestParseAfterEditMatchesParseAllForLineDeletion(t *testing.T) {
	oldText := "package main\n\nfunc main() {\n\tvar x = 42\n\tprintln(x)\n}\n"
	deleted := "\tvar x = 42\n"
	newText := "package main\n\nfunc main() {\n\tprintln(x)\n}\n"
	checkParseAfterEditMatchesParseAll(t, oldText, newText, Edit{
		Offset:     uint64(len("package main\n\nfunc main() {\n")),
		NumDeleted: uint64(len(deleted)),
	})
}

func TestParseAfterEditMatchesParseAllForStringReplacement(t *testing.T) {
	oldText := "package main\n\nconst message = \"hello\"\n"
	newText := "package main\n\nconst message = `hello\nworld`\n"
	checkParseAfterEditMatchesParseAll(t, oldText, newText, Edit{
		Offset:      uint64(len("package main\n\nconst message = ")),
		NumInserted: uint64(len("`hello\nworld`")),
		NumDeleted:  uint64(len("\"hello\"")),
	})
}

func TestParseAfterEditMatchesParseAllForCommentInsertion(t *testing.T) {
	oldText := "package main\n\nfunc main() {}\n"
	inserted := "// generated file\n"
	newText := inserted + oldText
	checkParseAfterEditMatchesParseAll(t, oldText, newText, Edit{
		Offset:      0,
		NumInserted: uint64(len(inserted)),
	})
}

func TestParseAfterEditMatchesParseAllForBlockCommentReplacement(t *testing.T) {
	oldText := "package main\n\nvar x = 1\n"
	newText := "package main\n\n/* var x = 1 */\n"
	checkParseAfterEditMatchesParseAll(t, oldText, newText, Edit{
		Offset:      uint64(len("package main\n\n")),
		NumInserted: uint64(len("/* var x = 1 */")),
		NumDeleted:  uint64(len("var x = 1")),
	})
}

func TestParseAfterEditMatchesParseAllForOperatorReplacement(t *testing.T) {
	oldText := "package main\n\nfunc main() {\n\tx := 1 + 2\n}\n"
	newText := "package main\n\nfunc main() {\n\tx := 1 << 2\n}\n"
	checkParseAfterEditMatchesParseAll(t, oldText, newText, Edit{
		Offset:      uint64(len("package main\n\nfunc main() {\n\tx := 1 ")),
		NumInserted: uint64(len("<<")),
		NumDeleted:  uint64(len("+")),
	})
}

func TestParseAfterEditMatchesParseAllForMultilineDeletion(t *testing.T) {
	oldText := "package main\n\nfunc main() {\n\tprintln(\"a\")\n\tprintln(\"b\")\n\tprintln(\"c\")\n}\n"
	deleted := "\tprintln(\"b\")\n\tprintln(\"c\")\n"
	newText := "package main\n\nfunc main() {\n\tprintln(\"a\")\n}\n"
	checkParseAfterEditMatchesParseAll(t, oldText, newText, Edit{
		Offset:     uint64(len("package main\n\nfunc main() {\n\tprintln(\"a\")\n")),
		NumDeleted: uint64(len(deleted)),
	})
}

func TestParseAfterEditMatchesParseAllForNumberReplacement(t *testing.T) {
	oldText := "package main\n\nconst n = 42\nconst f = 1.25\n"
	newText := "package main\n\nconst n = 0x2a\nconst f = 1.25\n"
	checkParseAfterEditMatchesParseAll(t, oldText, newText, Edit{
		Offset:      uint64(len("package main\n\nconst n = ")),
		NumInserted: uint64(len("0x2a")),
		NumDeleted:  uint64(len("42")),
	})
}

func TestParseAfterEditMatchesParseAllForManySequentialEdits(t *testing.T) {
	p := New(GolangParseFunc())
	text := "package main\n\nfunc main() {}\n"
	p.ParseAll(makeBuf(text))

	edits := []struct {
		newText string
		edit    Edit
	}{
		{
			newText: "package main\n\nfunc main() {\n}\n",
			edit: Edit{
				Offset:      uint64(len("package main\n\nfunc main() {")),
				NumInserted: uint64(len("\n")),
			},
		},
		{
			newText: "package main\n\nfunc main() {\n\tprintln(\"ok\")\n}\n",
			edit: Edit{
				Offset:      uint64(len("package main\n\nfunc main() {\n")),
				NumInserted: uint64(len("\tprintln(\"ok\")\n")),
			},
		},
		{
			newText: "package main\n\nfunc run() {\n\tprintln(\"ok\")\n}\n",
			edit: Edit{
				Offset:      uint64(len("package main\n\nfunc ")),
				NumInserted: uint64(len("run")),
				NumDeleted:  uint64(len("main")),
			},
		},
	}

	for _, tt := range edits {
		p.ParseAfterEdit(makeBuf(tt.newText), tt.edit)
		checkTokens(t, p.Tokens(), parseTokensFromAll(tt.newText))
		text = tt.newText
	}

	if text != edits[len(edits)-1].newText {
		t.Fatalf("unexpected final text")
	}
}

func TestParseAllLargeSource(t *testing.T) {
	source := makeLargeGoSource(70)
	checkSourceHasAtLeastLines(t, source, 300)

	p := New(GolangParseFunc())
	p.ParseAll(makeBuf(source))

	checkTokens(t, p.Tokens(), parseTokensFromAll(source))
	if len(p.computations) < 2 {
		t.Fatalf("expected large source to produce multiple computations, got %d", len(p.computations))
	}
}

func TestParseAfterEditLargeSourceAppendFunction(t *testing.T) {
	oldText := makeLargeGoSource(70)
	inserted := "func appended() {\n\tprintln(\"appended\")\n}\n"
	newText := oldText + inserted
	checkSourceHasAtLeastLines(t, newText, 300)

	checkParseAfterEditMatchesParseAll(t, oldText, newText, Edit{
		Offset:      uint64(len(oldText)),
		NumInserted: uint64(len(inserted)),
	})
}

func TestParseAfterEditLargeSourceInsertNearStart(t *testing.T) {
	oldText := makeLargeGoSource(70)
	inserted := "// file header inserted during edit\n"
	newText := inserted + oldText
	checkSourceHasAtLeastLines(t, newText, 300)

	checkParseAfterEditMatchesParseAll(t, oldText, newText, Edit{
		Offset:      0,
		NumInserted: uint64(len(inserted)),
	})
}

func TestParseAfterEditLargeSourceReplaceMiddleBlock(t *testing.T) {
	oldText := makeLargeGoSource(70)
	oldBlock := "\tvalue := 42 + 7\n\tprintln(\"value\", value)\n"
	newBlock := "\tvalue := 0x2a << 1\n\tprintln(`value`, value)\n"
	offset := len("package main\n\nconst banner = `large\nsource`\n\n") + len("func generated() {\n\t// repeated parser test line\n")
	newText := oldText[:offset] + newBlock + oldText[offset+len(oldBlock):]
	checkSourceHasAtLeastLines(t, newText, 300)

	checkParseAfterEditMatchesParseAll(t, oldText, newText, Edit{
		Offset:      uint64(offset),
		NumInserted: uint64(len(newBlock)),
		NumDeleted:  uint64(len(oldBlock)),
	})
}

func TestParseAfterEditLargeSourceDeleteManyFunctions(t *testing.T) {
	oldText := makeLargeGoSource(90)
	prefix := makeLargeGoSource(20)
	deletedStart := len(prefix)
	deleted := oldText[deletedStart:len(makeLargeGoSource(50))]
	newText := oldText[:deletedStart] + oldText[deletedStart+len(deleted):]
	checkSourceHasAtLeastLines(t, oldText, 300)
	checkSourceHasAtLeastLines(t, newText, 150)

	checkParseAfterEditMatchesParseAll(t, oldText, newText, Edit{
		Offset:     uint64(deletedStart),
		NumDeleted: uint64(len(deleted)),
	})
}

func TestParseAfterEditLargeSourceSequentialMixedEdits(t *testing.T) {
	p := New(GolangParseFunc())
	text := makeLargeGoSource(70)
	checkSourceHasAtLeastLines(t, text, 300)
	p.ParseAll(makeBuf(text))

	insertedHeader := "// sequential edit header\n"
	nextText := insertedHeader + text
	p.ParseAfterEdit(makeBuf(nextText), Edit{Offset: 0, NumInserted: uint64(len(insertedHeader))})
	checkTokens(t, p.Tokens(), parseTokensFromAll(nextText))
	text = nextText

	oldSnippet := "const banner = `large\nsource`"
	newSnippet := "const banner = \"large source\""
	offset := len(insertedHeader) + len("package main\n\n")
	nextText = text[:offset] + newSnippet + text[offset+len(oldSnippet):]
	p.ParseAfterEdit(makeBuf(nextText), Edit{
		Offset:      uint64(offset),
		NumInserted: uint64(len(newSnippet)),
		NumDeleted:  uint64(len(oldSnippet)),
	})
	checkTokens(t, p.Tokens(), parseTokensFromAll(nextText))
	text = nextText

	deleteStart := len(insertedHeader) + len("package main\n\n") + len(newSnippet) + len("\n\n")
	deleteEnd := deleteStart + len("func generated() {\n\t// repeated parser test line\n\tvalue := 42 + 7\n\tprintln(\"value\", value)\n}\n\n")
	nextText = text[:deleteStart] + text[deleteEnd:]
	p.ParseAfterEdit(makeBuf(nextText), Edit{
		Offset:     uint64(deleteStart),
		NumDeleted: uint64(deleteEnd - deleteStart),
	})
	checkTokens(t, p.Tokens(), parseTokensFromAll(nextText))
	checkSourceHasAtLeastLines(t, nextText, 300)
}

func TestParseAllMultilineBlockComment(t *testing.T) {
	source := "package main\n\n/*\nline one\nline two\n*/\nfunc main() {}\n"
	p := New(GolangParseFunc())
	p.ParseAll(makeBuf(source))

	checkTokens(t, p.Tokens(), []ComputedToken{
		{Offset: 0, Length: 7, Role: TokenRoleKeyword},
		{Offset: uint64(len("package main\n\n")), Length: uint64(len("/*\nline one\nline two\n*/")), Role: TokenRoleComment},
		{Offset: uint64(len("package main\n\n/*\nline one\nline two\n*/\n")), Length: 4, Role: TokenRoleKeyword},
		{Offset: uint64(len("package main\n\n/*\nline one\nline two\n*/\nfunc main")), Length: 1, Role: TokenRoleOperator},
		{Offset: uint64(len("package main\n\n/*\nline one\nline two\n*/\nfunc main(")), Length: 1, Role: TokenRoleOperator},
		{Offset: uint64(len("package main\n\n/*\nline one\nline two\n*/\nfunc main() ")), Length: 1, Role: TokenRoleOperator},
		{Offset: uint64(len("package main\n\n/*\nline one\nline two\n*/\nfunc main() {")), Length: 1, Role: TokenRoleOperator},
	})
}

func TestParseAfterEditMultilineBlockCommentInsertion(t *testing.T) {
	oldText := "package main\n\nfunc main() {}\n"
	inserted := "/*\nline one\nline two\n*/\n"
	newText := "package main\n\n" + inserted + "func main() {}\n"
	checkParseAfterEditMatchesParseAll(t, oldText, newText, Edit{
		Offset:      uint64(len("package main\n\n")),
		NumInserted: uint64(len(inserted)),
	})
}

func TestParseAfterEditMultilineBlockCommentContentReplacement(t *testing.T) {
	oldComment := "/*\nline one\nline two\n*/"
	newComment := "/*\nalpha\nbeta\ngamma\n*/"
	oldText := "package main\n\n" + oldComment + "\nfunc main() {}\n"
	newText := "package main\n\n" + newComment + "\nfunc main() {}\n"
	checkParseAfterEditMatchesParseAll(t, oldText, newText, Edit{
		Offset:      uint64(len("package main\n\n")),
		NumInserted: uint64(len(newComment)),
		NumDeleted:  uint64(len(oldComment)),
	})
}

func TestParseAfterEditMultilineBlockCommentDeletion(t *testing.T) {
	comment := "/*\nremove me\nacross lines\n*/\n"
	oldText := "package main\n\n" + comment + "func main() {}\n"
	newText := "package main\n\nfunc main() {}\n"
	checkParseAfterEditMatchesParseAll(t, oldText, newText, Edit{
		Offset:     uint64(len("package main\n\n")),
		NumDeleted: uint64(len(comment)),
	})
}

func TestParseAfterEditLargeSourceMultilineBlockCommentInsertion(t *testing.T) {
	oldText := makeLargeGoSource(70)
	inserted := "/*\nlarge source inserted comment\nwith several lines\nand operators like := + <<\n*/\n"
	offset := len("package main\n\nconst banner = `large\nsource`\n\n")
	newText := oldText[:offset] + inserted + oldText[offset:]
	checkSourceHasAtLeastLines(t, newText, 300)

	checkParseAfterEditMatchesParseAll(t, oldText, newText, Edit{
		Offset:      uint64(offset),
		NumInserted: uint64(len(inserted)),
	})
}

func TestParseAfterEditLargeSourceMultilineBlockCommentDeletion(t *testing.T) {
	comment := "/*\nlarge comment to delete\nline 2\nline 3\n*/\n"
	oldText := makeLargeGoSource(70)
	offset := len("package main\n\nconst banner = `large\nsource`\n\n")
	oldText = oldText[:offset] + comment + oldText[offset:]
	newText := oldText[:offset] + oldText[offset+len(comment):]
	checkSourceHasAtLeastLines(t, oldText, 300)
	checkSourceHasAtLeastLines(t, newText, 300)

	checkParseAfterEditMatchesParseAll(t, oldText, newText, Edit{
		Offset:     uint64(offset),
		NumDeleted: uint64(len(comment)),
	})
}

func TestAdvancePos(t *testing.T) {
	buf := makeBuf("hello\nworld\n!")

	tests := []struct {
		start Pos
		n     uint64
		want  Pos
	}{
		{Pos{0, 0}, 0, Pos{0, 0}},
		{Pos{0, 0}, 3, Pos{0, 3}},
		{Pos{0, 0}, 5, Pos{0, 5}},
		{Pos{0, 2}, 3, Pos{0, 5}},
		{Pos{0, 0}, 6, Pos{1, 0}},
		{Pos{0, 0}, 11, Pos{1, 5}},
		{Pos{0, 0}, 12, Pos{2, 0}},
		{Pos{1, 0}, 5, Pos{1, 5}},
	}

	for _, tt := range tests {
		got := advancePos(buf, tt.start, tt.n)
		if got != tt.want {
			t.Errorf("advancePos(%v, %d) = %v, want %v", tt.start, tt.n, got, tt.want)
		}
	}
}

func TestTotalChars(t *testing.T) {
	buf := makeBuf("hello\nworld")
	got := totalChars(buf)
	want := uint64(11) // 5 + "\n" + 5
	if got != want {
		t.Errorf("totalChars = %d, want %d", got, want)
	}

	emptyBuf := makeBuf("")
	if totalChars(emptyBuf) != 0 {
		t.Errorf("totalChars of empty buf should be 0")
	}
}
