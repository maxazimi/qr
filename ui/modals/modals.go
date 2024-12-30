package modals

import "gioui.org/layout"

type (
	C = layout.Context
	D = layout.Dimensions
)

type Modal interface {
	IsVisible() bool
	Layout(gtx C) D
}

var (
	modals []Modal
)

func Items() []Modal {
	return modals
}

func Add(m Modal) {
	if m != nil {
		modals = appendIfMissing(modals, m)
	}
}

func appendIfMissing(slice []Modal, m Modal) []Modal {
	for _, elem := range slice {
		if elem == m {
			return slice
		}
	}
	return append(slice, m)
}
