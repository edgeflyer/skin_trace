package analysis

// Result 表示一张照片的分析结果。
type Result struct {
	Summary      string        `json:"summary"`
	Observations []Observation `json:"observations"`
	Limitations  []string      `json:"limitations"`
}

// Observation 表示照片中某个部位的可见现象。
type Observation struct {
	Region      string `json:"region"`
	Description string `json:"description"`
}
