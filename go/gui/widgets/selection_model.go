package widgets

type selectionOrder[T comparable] struct {
	Len     func() int
	ItemAt  func(int) T
	IndexOf func(T) int
}

type selectionModel[T comparable] struct {
	lead        T
	leadOK      bool
	recent      T
	recentOK    bool
	anchor      T
	anchorOK    bool
	selectedSet map[T]struct{}
	dragBaseSet map[T]struct{}
}

type selectionSnapshot[T comparable] struct {
	lead        T
	leadOK      bool
	anchor      T
	anchorOK    bool
	selectedSet map[T]struct{}
}

func newSelectionModel[T comparable]() selectionModel[T] {
	return selectionModel[T]{
		selectedSet: make(map[T]struct{}),
		dragBaseSet: make(map[T]struct{}),
	}
}

func (s *selectionModel[T]) Lead() (T, bool) {
	return s.lead, s.leadOK
}

func (s *selectionModel[T]) Anchor() (T, bool) {
	return s.anchor, s.anchorOK
}

func (s *selectionModel[T]) Recent() (T, bool) {
	return s.recent, s.recentOK
}

func (s *selectionModel[T]) SetLead(item T) {
	s.lead = item
	s.leadOK = true
	s.recent = item
	s.recentOK = true
}

func (s *selectionModel[T]) ClearLead() {
	var zero T
	s.lead = zero
	s.leadOK = false
}

func (s *selectionModel[T]) SetAnchor(item T) {
	s.anchor = item
	s.anchorOK = true
}

func (s *selectionModel[T]) ClearAnchor() {
	var zero T
	s.anchor = zero
	s.anchorOK = false
}

func (s *selectionModel[T]) Count() int {
	return len(s.selectedSet)
}

func (s *selectionModel[T]) Contains(item T) bool {
	_, ok := s.selectedSet[item]
	return ok
}

func (s *selectionModel[T]) Clear() bool {
	snap := s.snapshot()
	clear(s.selectedSet)
	s.ClearLead()
	s.ClearAnchor()
	return snap.changed(s)
}

func (s *selectionModel[T]) SelectOnly(item T) bool {
	snap := s.snapshot()
	clear(s.selectedSet)
	s.selectedSet[item] = struct{}{}
	s.SetLead(item)
	s.SetAnchor(item)
	return snap.changed(s)
}

func (s *selectionModel[T]) EnsureLeadSelected() bool {
	if !s.leadOK {
		return false
	}
	if s.Contains(s.lead) && s.Count() == 1 {
		return false
	}
	snap := s.snapshot()
	clear(s.selectedSet)
	s.selectedSet[s.lead] = struct{}{}
	return snap.changed(s)
}

func (s *selectionModel[T]) Toggle(item T, order selectionOrder[T]) bool {
	if s.Contains(item) {
		if s.Count() == 1 && s.leadOK && s.lead == item {
			return false
		}
		snap := s.snapshot()
		delete(s.selectedSet, item)
		if s.leadOK && s.lead == item {
			if replacement, ok := s.firstSelectedInOrder(order); ok {
				s.SetLead(replacement)
			} else {
				s.selectedSet[item] = struct{}{}
				s.SetLead(item)
				return false
			}
		}
		s.SetAnchor(item)
		return snap.changed(s)
	}
	snap := s.snapshot()
	s.selectedSet[item] = struct{}{}
	s.SetLead(item)
	s.SetAnchor(item)
	return snap.changed(s)
}

func (s *selectionModel[T]) SelectRange(order selectionOrder[T], target T) bool {
	targetIndex := order.IndexOf(target)
	if targetIndex < 0 {
		return false
	}
	anchor := target
	if s.anchorOK && order.IndexOf(s.anchor) >= 0 {
		anchor = s.anchor
	} else if s.leadOK && order.IndexOf(s.lead) >= 0 {
		anchor = s.lead
	}
	start := order.IndexOf(anchor)
	if start < 0 {
		start = targetIndex
	}
	end := targetIndex
	if start > end {
		start, end = end, start
	}
	snap := s.snapshot()
	clear(s.selectedSet)
	for i := start; i <= end; i++ {
		s.selectedSet[order.ItemAt(i)] = struct{}{}
	}
	s.SetLead(target)
	return snap.changed(s)
}

func (s *selectionModel[T]) SelectAll(order selectionOrder[T]) bool {
	if order.Len() == 0 {
		return false
	}
	lead := order.ItemAt(0)
	if s.leadOK && order.IndexOf(s.lead) >= 0 {
		lead = s.lead
	}
	snap := s.snapshot()
	clear(s.selectedSet)
	for i := 0; i < order.Len(); i++ {
		s.selectedSet[order.ItemAt(i)] = struct{}{}
	}
	s.SetLead(lead)
	if !s.anchorOK || order.IndexOf(s.anchor) < 0 {
		s.SetAnchor(lead)
	}
	return snap.changed(s)
}

func (s *selectionModel[T]) SelectDragRange(order selectionOrder[T], anchor, target T) bool {
	start := order.IndexOf(anchor)
	end := order.IndexOf(target)
	if start < 0 || end < 0 {
		return false
	}
	if start > end {
		start, end = end, start
	}
	snap := s.snapshot()
	clear(s.selectedSet)
	for i := start; i <= end; i++ {
		s.selectedSet[order.ItemAt(i)] = struct{}{}
	}
	s.SetLead(target)
	s.SetAnchor(anchor)
	return snap.changed(s)
}

func (s *selectionModel[T]) SelectDragUnion(order selectionOrder[T], anchor, target T) bool {
	start := order.IndexOf(anchor)
	end := order.IndexOf(target)
	if start < 0 || end < 0 {
		return false
	}
	if start > end {
		start, end = end, start
	}
	snap := s.snapshot()
	clear(s.selectedSet)
	for item := range s.dragBaseSet {
		if order.IndexOf(item) >= 0 {
			s.selectedSet[item] = struct{}{}
		}
	}
	for i := start; i <= end; i++ {
		s.selectedSet[order.ItemAt(i)] = struct{}{}
	}
	s.SetLead(target)
	s.SetAnchor(anchor)
	return snap.changed(s)
}

func (s *selectionModel[T]) ApplyMarquee(order selectionOrder[T], intersects func(T) bool, unionBase bool) bool {
	snap := s.snapshot()
	lead, leadOK := s.lead, s.leadOK
	anchor, anchorOK := s.anchor, s.anchorOK
	clear(s.selectedSet)
	if unionBase {
		for item := range s.dragBaseSet {
			if order.IndexOf(item) >= 0 {
				s.selectedSet[item] = struct{}{}
			}
		}
	} else {
		leadOK = false
		anchorOK = false
	}
	for i := 0; i < order.Len(); i++ {
		item := order.ItemAt(i)
		if !intersects(item) {
			continue
		}
		s.selectedSet[item] = struct{}{}
		if !anchorOK {
			anchor = item
			anchorOK = true
		}
		lead = item
		leadOK = true
	}
	if len(s.selectedSet) == 0 {
		s.ClearLead()
		s.ClearAnchor()
	} else {
		s.lead = lead
		s.leadOK = leadOK
		s.anchor = anchor
		s.anchorOK = anchorOK
	}
	return snap.changed(s)
}

func (s *selectionModel[T]) CaptureDragBase() {
	clear(s.dragBaseSet)
	for item := range s.selectedSet {
		s.dragBaseSet[item] = struct{}{}
	}
}

func (s *selectionModel[T]) DropInvalid(valid func(T) bool) {
	for item := range s.selectedSet {
		if !valid(item) {
			delete(s.selectedSet, item)
		}
	}
	for item := range s.dragBaseSet {
		if !valid(item) {
			delete(s.dragBaseSet, item)
		}
	}
	if s.leadOK && !valid(s.lead) {
		s.ClearLead()
	}
	if s.anchorOK && !valid(s.anchor) {
		s.ClearAnchor()
	}
	if s.recentOK && !valid(s.recent) {
		var zero T
		s.recent = zero
		s.recentOK = false
	}
}

func (s *selectionModel[T]) snapshot() selectionSnapshot[T] {
	set := make(map[T]struct{}, len(s.selectedSet))
	for item := range s.selectedSet {
		set[item] = struct{}{}
	}
	return selectionSnapshot[T]{
		lead:        s.lead,
		leadOK:      s.leadOK,
		anchor:      s.anchor,
		anchorOK:    s.anchorOK,
		selectedSet: set,
	}
}

func (snap selectionSnapshot[T]) changed(s *selectionModel[T]) bool {
	if snap.leadOK != s.leadOK || snap.anchorOK != s.anchorOK || len(snap.selectedSet) != len(s.selectedSet) {
		return true
	}
	if snap.leadOK && snap.lead != s.lead {
		return true
	}
	if snap.anchorOK && snap.anchor != s.anchor {
		return true
	}
	for item := range s.selectedSet {
		if _, ok := snap.selectedSet[item]; !ok {
			return true
		}
	}
	return false
}

func (s *selectionModel[T]) firstSelectedInOrder(order selectionOrder[T]) (T, bool) {
	for i := 0; i < order.Len(); i++ {
		item := order.ItemAt(i)
		if s.Contains(item) {
			return item, true
		}
	}
	var zero T
	return zero, false
}
