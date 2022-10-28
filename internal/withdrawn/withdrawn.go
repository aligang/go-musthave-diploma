package withdrawn

import (
	"time"
)

type Withdrawn struct {
	OrderID string  `json:"order"`
	Sum     float64 `json:"sum"`
}

func New() *Withdrawn {
	return &Withdrawn{}
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

//func Sort(withdraws []WithdrawnRecord) {
//	sort.Slice(withdraws, func(i, j int) bool {
//		return withdraws[i].ProcessedAt.After(withdraws[j].ProcessedAt)
//	})
//}

type WithdrawnSlice []WithdrawnRecord

func (s WithdrawnSlice) Len() int {
	return len(s)
}

func (s WithdrawnSlice) Less(i, j int) bool {
	return s[i].ProcessedAt.After(s[j].ProcessedAt)
}

func (s WithdrawnSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
