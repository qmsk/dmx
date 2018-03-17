package api

type GroupConfig struct {
	Heads []HeadID
	Name  string
}

type GroupID string

type Groups map[GroupID]Group

type Group struct {
	GroupConfig // TODO: separate config field
	ID          GroupID
	Heads       []HeadID
	Colors      Colors

	GroupParams
}

type GroupParams struct {
	Intensity *Intensity `json:",omitempty"`
	Color     *Color     `json:",omitempty"`
}
