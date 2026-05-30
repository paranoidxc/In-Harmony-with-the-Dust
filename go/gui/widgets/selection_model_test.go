package widgets

import "testing"

func TestSelectionModelToggleRetainsFirstVisibleLead(t *testing.T) {
	model := newSelectionModel[int]()
	order := selectionOrder[int]{
		Len:    func() int { return 4 },
		ItemAt: func(index int) int { return index },
		IndexOf: func(index int) int {
			if index < 0 || index >= 4 {
				return -1
			}
			return index
		},
	}

	model.SelectOnly(1)
	model.Toggle(3, order)
	if !model.Toggle(3, order) {
		t.Fatal("toggling selected lead off should change selection")
	}

	lead, ok := model.Lead()
	if !ok || lead != 1 {
		t.Fatalf("lead = (%d, %v), want (1, true)", lead, ok)
	}
	if !model.Contains(1) || model.Contains(3) {
		t.Fatal("selection contents should fall back to item 1 only")
	}
}

func TestSelectionModelApplyMarqueeCanUnionBaseSelection(t *testing.T) {
	model := newSelectionModel[int]()
	order := selectionOrder[int]{
		Len:    func() int { return 5 },
		ItemAt: func(index int) int { return index },
		IndexOf: func(index int) int {
			if index < 0 || index >= 5 {
				return -1
			}
			return index
		},
	}

	model.SelectOnly(0)
	model.Toggle(4, order)
	model.CaptureDragBase()

	changed := model.ApplyMarquee(order, func(index int) bool {
		return index == 1 || index == 2
	}, true)
	if !changed {
		t.Fatal("marquee union should change selection")
	}
	for _, index := range []int{0, 1, 2, 4} {
		if !model.Contains(index) {
			t.Fatalf("selection should contain %d", index)
		}
	}
	lead, ok := model.Lead()
	if !ok || lead != 2 {
		t.Fatalf("lead = (%d, %v), want (2, true)", lead, ok)
	}
}
