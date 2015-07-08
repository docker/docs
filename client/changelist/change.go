package changelist

// tufChange represents a change to a TUF repo
type tufChange struct {
	// Abbreviated because Go doesn't permit a field and method of the same name
	Actn       int    `json:"action"`
	Role       string `json:"role"`
	ChangeType string `json:"type"`
	ChangePath string `json:"path"`
	Data       []byte `json:"data"`
}

// NewTufChange initializes a tufChange object
func NewTufChange(action int, role, changeType, changePath string, content []byte) *tufChange {
	return &tufChange{
		Actn:       action,
		Role:       role,
		ChangeType: changeType,
		ChangePath: changePath,
		Data:       content,
	}
}

// Action return c.Actn
func (c tufChange) Action() int {
	return c.Actn
}

// Scope returns c.Role
func (c tufChange) Scope() string {
	return c.Role
}

// Type returns c.ChangeType
func (c tufChange) Type() string {
	return c.ChangeType
}

// Path return c.ChangePath
func (c tufChange) Path() string {
	return c.ChangePath
}

func (c tufChange) Content() []byte {
	return c.Data
}
