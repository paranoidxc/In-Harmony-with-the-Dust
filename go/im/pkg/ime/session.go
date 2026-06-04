package ime

import "strings"

type Session struct {
	engine    *Engine
	limit     int
	buffer    []rune
	committed []string
	result    SearchResult
	selected  int
}

func NewSession(engine *Engine, limit int) *Session {
	if limit <= 0 {
		limit = 10
	}
	return &Session{
		engine: engine,
		limit:  limit,
	}
}

func (s *Session) InputRune(r rune) {
	if r < 32 || r > 126 {
		return
	}
	s.buffer = append(s.buffer, r)
	s.refresh()
}

func (s *Session) Backspace() {
	if len(s.buffer) == 0 {
		return
	}
	s.buffer = s.buffer[:len(s.buffer)-1]
	s.refresh()
}

func (s *Session) ClearInput() {
	s.buffer = s.buffer[:0]
	s.result = SearchResult{}
	s.selected = 0
}

func (s *Session) MoveSelection(delta int) {
	if len(s.result.Items) == 0 {
		s.selected = 0
		return
	}
	s.selected += delta
	if s.selected < 0 {
		s.selected = len(s.result.Items) - 1
	}
	if s.selected >= len(s.result.Items) {
		s.selected = 0
	}
}

func (s *Session) CommitSelection() string {
	if len(s.result.Items) == 0 {
		return ""
	}
	word := s.result.Items[s.selected].Word
	s.committed = append(s.committed, word)
	s.ClearInput()
	return word
}

func (s *Session) CommitIndex(index int) string {
	if index < 0 || index >= len(s.result.Items) {
		return ""
	}
	s.selected = index
	return s.CommitSelection()
}

func (s *Session) Buffer() string {
	return string(s.buffer)
}

func (s *Session) Candidates() []Candidate {
	return s.result.Items
}

func (s *Session) Log() MatchLog {
	return s.result.Log
}

func (s *Session) SelectedIndex() int {
	return s.selected
}

func (s *Session) CommittedText() string {
	return strings.Join(s.committed, "")
}

func (s *Session) ResetCommitted() {
	s.committed = nil
}

func (s *Session) refresh() {
	s.selected = 0
	if len(s.buffer) == 0 {
		s.result = SearchResult{}
		return
	}
	s.result = s.engine.Search(string(s.buffer), s.limit)
}
