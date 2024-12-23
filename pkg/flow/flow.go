package flow

type Flow struct {
	Name  string `json:"name"`
	Steps []Step `json:"steps"`
}

type Step struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Action any    `json:"action"` // rename, filter, etc.
}

type Rename struct {
	Name  string `json:"name"`
	Kind  string `json:"kind"`
	Mode  string `json:"mode"`
	Expr  string `json:"expr"`
	Index int    `json:"index"`
	Value string `json:"value"`
}

func (r Rename) Add(str string) string {
	return str
}

func (r Rename) Delete(str string) string {
	return str
}

func (r Rename) Replace(str string) string {
	return str
}

func (r Rename) Shift(str string) string {
	return str
}
