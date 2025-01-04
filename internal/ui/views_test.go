package ui

import (
	"testing"
)

func TestTabTextView_SetNext(t *testing.T) {
	view1 := NewTabTextView(nil)
	view2 := NewTabTextView(nil)
	view1.SetNext(view2)
	if view1.GetNext() != view2 {
		t.Errorf("Expected view2 to be the next view of view1")
	}
}

func TestTabTextView_GetNext(t *testing.T) {
	view1 := NewTabTextView(nil)
	view2 := NewTabTextView(nil)
	view1.SetNext(view2)
	if view1.GetNext() != view2 {
		t.Errorf("Expected view2 to be the next view of view1")
	}
}

func TestNewTabTextView(t *testing.T) {
	view := NewTabTextView(nil)
	if view == nil {
		t.Errorf("Expected new TabTextView to be not nil")
	}
	if view.GetNext() != nil {
		t.Errorf("Expected next view to be nil")
	}
}
