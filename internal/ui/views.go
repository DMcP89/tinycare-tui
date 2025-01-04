package ui

import "github.com/rivo/tview"

type TabTextView struct {
	*tview.TextView
	next *TabTextView
}

func (view *TabTextView) SetNext(next *TabTextView) *TabTextView {
	view.next = next
	return view
}

func (view *TabTextView) GetNext() *TabTextView {
	return view.next
}

func NewTabTextView(next *TabTextView) *TabTextView {
	return &TabTextView{
		TextView: tview.NewTextView(),
		next:     next,
	}
}
