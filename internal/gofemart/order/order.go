package order

import (
	"github.com/aligang/go-musthave-diploma/internal/accural/message"
	"time"
)

type Order struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accural    float64   `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}

//func Sort(orders []Order) {
//	sort.Slice(orders, func(i, j int) bool {
//		return orders[i].UploadedAt.After(orders[j].UploadedAt)
//	})
//}

type OrderSlice []Order

func (s OrderSlice) Len() int {
	return len(s)
}

func (s OrderSlice) Less(i, j int) bool {
	return s[i].UploadedAt.After(s[j].UploadedAt)
}

func (s OrderSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func FromAccural(record *message.AccuralMessage) *Order {
	return &Order{
		Number:     record.Order,
		Status:     record.Status,
		Accural:    record.Accural,
		UploadedAt: time.Now().Round(time.Second),
	}
}
