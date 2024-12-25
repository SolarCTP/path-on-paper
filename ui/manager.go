package ui

import (
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
)

type uiManager struct {
	EUI   *ebitenui.UI
	Stack *UIStack
}

func setupUiManager() {
	Manager = uiManager{
		EUI: &ebitenui.UI{
			Container: widget.NewContainer(),
		},
		Stack: NewUIStack(),
	}
}

// He talks to ebiten ui to provide the user interface
var Manager uiManager
