package ui

type Signal uint16

const (
	SignalShowedMenu Signal = iota
	SignalHidMenu
)

var SignalChannel chan Signal

func setupSignalChannel() {
	SignalChannel = make(chan Signal, 10)
}
