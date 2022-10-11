package withdrawn

import (
	"sort"
	"time"
)

type Withdrawn struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

type WithdrawnRecord struct {
	*Withdrawn
	ProcessedAt time.Time `json:"processed_at"`
}

type Withdrawns map[string]WithdrawnRecord

func NewRecord(withdrawn *Withdrawn) *WithdrawnRecord {
	return &WithdrawnRecord{
		Withdrawn:   withdrawn,
		ProcessedAt: time.Now().Round(time.Second),
	}
}

func Sort(withdraws []WithdrawnRecord) {
	sort.Slice(withdraws, func(i, j int) bool {
		return withdraws[i].ProcessedAt.After(withdraws[j].ProcessedAt)
	})
}
