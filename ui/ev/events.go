package ev

type Event interface{}

type AppLoadEvent struct{}

type WindowSizeEvent struct {
	Width, Height int
}

type AddModalEvent struct {
	Modal interface{}
}

type QREvent struct{}

type ExitEvent struct{}

type FatalEvent struct {
	Err error
}

type NotifyEvent struct {
	Level int
	Text  string
}

var (
	uiEvent = make(chan Event)
)

func RequestEvent(e Event) {
	go func() {
		uiEvent <- e
	}()
}

func Events() <-chan Event {
	return uiEvent
}
