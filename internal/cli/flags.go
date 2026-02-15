package cli

type Action string

const (
	ActionNone     Action = ""
	ActionAdd      Action = "add"
	ActionShow     Action = "show"
	ActionComplete Action = "complete"
	ActionOpen     Action = "open"
	ActionSync     Action = "sync"
	ActionEdit     Action = "edit"
	ActionDelete   Action = "delete"
)

type Flags struct {
	AddText    string
	Show       bool
	CompleteID int
	OpenID     int
	Sync       bool
	EditID     int
	EditText   string
	DeleteID   int
	Status     string // fica string aqui (parse depois)
}
