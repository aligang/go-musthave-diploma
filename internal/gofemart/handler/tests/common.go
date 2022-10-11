package tests

import (
	"github.com/aligang/go-musthave-diploma/internal/gofemart/storage/memory"
	"time"
)

type expected struct {
	path        string
	code        int
	contentType string
	payload     string
	dbDump      *memory.Storage
}

type input struct {
	account     string
	method      string
	path        string
	contentType string
	payload     string
}

func genTimeStamps() []time.Time {
	var timeStamps []time.Time
	for i := 0; i < 10; i++ {
		timeStamp, _ := time.Parse(time.RFC3339, "2021-09-19T15:59:41+03:00")
		timeStamps = append(timeStamps, timeStamp.Add(time.Duration(i)*time.Second))
	}
	return timeStamps
}
