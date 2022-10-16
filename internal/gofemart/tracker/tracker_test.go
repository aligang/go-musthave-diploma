package tracker

import (
	"context"
	accural_handler "github.com/aligang/go-musthave-diploma/internal/accural/handler"
	"github.com/aligang/go-musthave-diploma/internal/accural/message"
	accural_storage "github.com/aligang/go-musthave-diploma/internal/accural/storage/memory"
	"github.com/aligang/go-musthave-diploma/internal/config"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/account"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/order"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/storage/memory"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/storage/memory/internal_order"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/tests_common"
	"github.com/aligang/go-musthave-diploma/internal/withdrawn"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestTracking(t *testing.T) {
	tests := []struct {
		name     string
		input    tests_common.Input
		expected tests_common.Expected
	}{
		{
			name: "TEST TRACKER",
			expected: tests_common.Expected{
				Code:        202,
				ContentType: "text/plain",
				DbDump: memory.Init(
					account.CustomerAccounts{
						"user1": account.CustomerAccount{
							Login:    "user1",
							Password: "pass1",
							AccountBalance: account.AccountBalance{
								Current:  21,
								Withdraw: 0,
							},
						},
					},
					internal_order.Orders{
						"19": internal_order.Order{
							Order: &order.Order{
								Number:     "19",
								Accural:    10.5,
								Status:     "PROCESSED",
								UploadedAt: time.Now().Round(time.Second),
							},
							Owner: "user1",
						},
						"20": internal_order.Order{
							Order: &order.Order{
								Number:     "20",
								Accural:    10.5,
								Status:     "PROCESSED",
								UploadedAt: time.Now().Round(time.Second),
							},
							Owner: "user1",
						},
						"21": internal_order.Order{
							Order: &order.Order{
								Number:     "21",
								Accural:    10.5,
								Status:     "PROCESSED",
								UploadedAt: time.Now().Round(time.Second),
							},
							Owner: "user1",
						},
						"22": internal_order.Order{
							Order: &order.Order{
								Number:     "22",
								Accural:    10.5,
								Status:     "INVALID",
								UploadedAt: time.Now().Round(time.Second),
							},
							Owner: "user1",
						},
					},
					withdrawn.Withdrawns{},
					map[string]*struct{}{},
					map[string][]string{"user1": {"19", "20", "21", "22"}},
					map[string][]string{},
				),
			},
		},
	}

	//Starting Test AccuralServer
	accuralStorage := accural_storage.Init(message.AccuralMessageMap{
		"19": message.AccuralMessage{
			Order:   "19",
			Status:  "PROCESSED",
			Accural: 10.5,
		},
		"20": message.AccuralMessage{
			Order:   "20",
			Status:  "PROCESSED",
			Accural: 10.5,
		},
		"21": message.AccuralMessage{
			Order:   "21",
			Status:  "PROCESSED",
			Accural: 10.5,
		},
		"22": message.AccuralMessage{
			Order:   "22",
			Status:  "INVALID",
			Accural: 10.5,
		},
	})
	accuralMux := accural_handler.New(accuralStorage)
	accuralMux.Get("/api/orders/{order}", accuralMux.Fetch)
	accuralServer := httptest.NewServer(accuralMux)

	//Starting Test ApplicationServer
	cfg := &config.Config{
		AccuralSystemAddress: strings.Split(accuralServer.URL, "/")[2],
		DatabaseURI:          "",
		RunAddress:           "",
	}

	storage := memory.Init(
		account.CustomerAccounts{
			"user1": account.CustomerAccount{
				Login:    "user1",
				Password: "pass1",
				AccountBalance: account.AccountBalance{
					Current:  0,
					Withdraw: 0,
				},
			},
		},
		internal_order.Orders{
			"19": internal_order.Order{
				Order: &order.Order{
					Number:     "19",
					Accural:    10.5,
					Status:     "PROCESSED",
					UploadedAt: time.Now().Round(time.Second),
				},
				Owner: "user1",
			},
			"20": internal_order.Order{
				Order: &order.Order{
					Number:     "20",
					Accural:    10.5,
					Status:     "PROCESSING",
					UploadedAt: time.Now().Round(time.Second),
				},
				Owner: "user1",
			},
			"21": internal_order.Order{
				Order: &order.Order{
					Number:     "21",
					Accural:    10.5,
					Status:     "NEW",
					UploadedAt: time.Now().Round(time.Second),
				},
				Owner: "user1",
			},
			"22": internal_order.Order{
				Order: &order.Order{
					Number:     "22",
					Accural:    10.5,
					Status:     "INVALID",
					UploadedAt: time.Now().Round(time.Second),
				},
				Owner: "user1",
			},
		},
		withdrawn.Withdrawns{},
		map[string]*struct{}{"20": nil, "21": nil},
		map[string][]string{"user1": {"19", "20", "21", "22"}},
		map[string][]string{},
	)
	Tracker := New(storage, cfg)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			Tracker.Sync(context.Background())
			assert.Equal(t, reflect.DeepEqual(*test.expected.DbDump, *storage), true)
		},
		)
	}
}
