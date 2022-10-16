package testscommon

import (
	"github.com/aligang/go-musthave-diploma/internal/gofemart/storage/memory"
	"time"
)

type Expected struct {
	Path        string
	Code        int
	ContentType string
	Payload     string
	DBdump      *memory.Storage
}

type Input struct {
	Account     string
	Method      string
	Path        string
	ContentType string
	Payload     string
}

func GenTimeStamps() []time.Time {
	var timeStamps []time.Time
	for i := 0; i < 10; i++ {
		timeStamp, _ := time.Parse(time.RFC3339, "2021-09-19T15:59:41+03:00")
		timeStamps = append(timeStamps, timeStamp.Add(time.Duration(i)*time.Second))
	}
	return timeStamps
}
