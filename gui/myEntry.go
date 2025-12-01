package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type SubmitEntry struct {
	widget.Entry
	OnSubmit func()

	shiftDown bool
}

func NewSubmitEntry() *SubmitEntry {
	entry := &SubmitEntry{}
	entry.ExtendBaseWidget(entry)
	entry.MultiLine = true
	entry.Wrapping = fyne.TextWrapWord
	return entry
}

func (e *SubmitEntry) KeyDown(key *fyne.KeyEvent) {
	if fyne.KeyModifier(key.Physical.ScanCode) == 50 {
		e.shiftDown = true
	}
}

func (e *SubmitEntry) KeyUp(key *fyne.KeyEvent) {
	if fyne.KeyModifier(key.Physical.ScanCode) == 50 {
		e.shiftDown = false
	}
}

func (e *SubmitEntry) TypedKey(key *fyne.KeyEvent) {
	if key.Name == fyne.KeyReturn && !e.shiftDown {
		if e.OnSubmit != nil {
			e.OnSubmit()
		}
		return
	}

	e.Entry.TypedKey(key)
}
