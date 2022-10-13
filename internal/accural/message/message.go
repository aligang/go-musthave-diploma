package message

type AccuralMessage struct {
	Order   string  `json:"order"`
	Status  string  `json:"status""`
	Accural float64 `json:"accrual,omitempty"`
}

type AccuralMessageMap map[string]AccuralMessage
