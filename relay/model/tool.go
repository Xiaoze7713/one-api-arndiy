package model

type Tool struct {
	Id       string   `json:"id,omitempty"`
	Type     string   `json:"type,omitempty"` // when splicing claude tools stream messages, it is empty
	Function Function `json:"function"`
}

type Function struct {
	Description string     `json:"description,omitempty"`
	Name        string     `json:"name,omitempty"`       // when splicing claude tools stream messages, it is empty
	Parameters  *Parameter `json:"parameters,omitempty"` // request
	Arguments   any        `json:"arguments,omitempty"`  // response
}

type Property struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Enum        []string `json:"enum"`
}

type Parameter struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
	Required   []string            `json:"required"`
}
