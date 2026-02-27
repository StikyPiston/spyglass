package lens

type Entry struct {
	ID          string
	Title       string
	Icon        string
	Description string
}

type Action struct {
	Name string
	Run  func(Entry) error
}

type Lens interface {
	Name() string
	Search(query string) ([]Entry, error)
	Enter(entry Entry) error
	ContextActions(entry Entry) []Action
}
