package calc

type Calculator struct {
	vars map[string]int64
}

type Instruction struct {
	Type  string      `json:"type"`
	Op    string      `json:"op,omitempty"`
	Var   string      `json:"var,omitempty"`
	Left  interface{} `json:"left,omitempty"`
	Right interface{} `json:"right,omitempty"`
}

type Result struct {
	Var   string `json:"var"`
	Value int64  `json:"value"`
}