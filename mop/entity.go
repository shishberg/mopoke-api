package mop

type Entity struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Rel         []Rel  `json:"rel,omitempty"`
}

type Rel struct {
	From Entity `json:"from,omitempty"`
	To   Entity `json:"to,omitempty"`
	Type string `json:"type,omitempty"`
}
