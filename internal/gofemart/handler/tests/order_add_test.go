package tests

import (
	"bytes"
	"fmt"
	accural_handler "github.com/aligang/go-musthave-diploma/internal/accural/handler"
	"github.com/aligang/go-musthave-diploma/internal/accural/message"
	accural_storage "github.com/aligang/go-musthave-diploma/internal/accural/storage/memory"
	"github.com/aligang/go-musthave-diploma/internal/config"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/account"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/auth"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/handler"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/order"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/storage/memory"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/storage/memory/orderrecord"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/testscommon"
	"github.com/aligang/go-musthave-diploma/internal/withdrawn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestAddOrder(t *testing.T) {
	tests := []struct {
		name     string
		input    testscommon.Input
		expected testscommon.Expected
	}{
		{
			name: "PROCESSED ORDER REGISTRATION CORRECT",
			input: testscommon.Input{Method: http.MethodPost, Path: "/api/user/orders", ContentType: "application/json",
				Payload: "18", Account: "user1"},
			expected: testscommon.Expected{
				Code:        202,
				ContentType: "text/plain",
				DbDump: memory.Init(
					account.CustomerAccounts{
						"user1": account.CustomerAccount{
							Login:    "user1",
							Password: "pass1",
							AccountBalance: account.AccountBalance{
								Current:  10.5,
								Withdraw: 0,
							},
						},
						"user2": account.CustomerAccount{
							Login:    "user2",
							Password: "pass2",
							AccountBalance: account.AccountBalance{
								Current:  0,
								Withdraw: 0,
							},
						},
					},
					orderrecord.Orders{
						"18": orderrecord.Order{
							Order: &order.Order{
								Number:     "18",
								Accural:    10.5,
								Status:     "PROCESSED",
								UploadedAt: time.Now().Round(time.Second),
							},
							Owner: "user1",
						},
					},
					withdrawn.Withdrawns{},
					map[string]*struct{}{},
					map[string][]string{"user1": {"18"}},
					map[string][]string{},
				),
			},
		},

		{
			name: "PROCESSING ORDER REGISTRATION CORRECT",
			input: testscommon.Input{Method: http.MethodPost, Path: "/api/user/orders", ContentType: "application/json",
				Payload: "26", Account: "user1"},
			expected: testscommon.Expected{
				Code:        202,
				ContentType: "text/plain",
				DbDump: memory.Init(
					account.CustomerAccounts{
						"user1": account.CustomerAccount{
							Login:    "user1",
							Password: "pass1",
							AccountBalance: account.AccountBalance{
								Current:  10.5,
								Withdraw: 0,
							},
						},
						"user2": account.CustomerAccount{
							Login:    "user2",
							Password: "pass2",
							AccountBalance: account.AccountBalance{
								Current:  0,
								Withdraw: 0,
							},
						},
					},
					orderrecord.Orders{
						"18": orderrecord.Order{
							Order: &order.Order{
								Number:     "18",
								Accural:    10.5,
								Status:     "PROCESSED",
								UploadedAt: time.Now().Round(time.Second),
							},
							Owner: "user1",
						},
						"26": orderrecord.Order{
							Order: &order.Order{
								Number:     "26",
								Accural:    10.5,
								Status:     "PROCESSING",
								UploadedAt: time.Now().Round(time.Second),
							},
							Owner: "user1",
						},
					},
					withdrawn.Withdrawns{},
					map[string]*struct{}{"26": nil},
					map[string][]string{"user1": {"18", "26"}},
					map[string][]string{},
				),
			},
		},

		{
			name: "NEW ORDER REGISTRATION CORRECT",
			input: testscommon.Input{Method: http.MethodPost, Path: "/api/user/orders", ContentType: "application/json",
				Payload: "34", Account: "user1"},
			expected: testscommon.Expected{
				Code:        202,
				ContentType: "text/plain",
				DbDump: memory.Init(
					account.CustomerAccounts{
						"user1": account.CustomerAccount{
							Login:    "user1",
							Password: "pass1",
							AccountBalance: account.AccountBalance{
								Current:  10.5,
								Withdraw: 0,
							},
						},
						"user2": account.CustomerAccount{
							Login:    "user2",
							Password: "pass2",
							AccountBalance: account.AccountBalance{
								Current:  0,
								Withdraw: 0,
							},
						},
					},
					orderrecord.Orders{
						"18": orderrecord.Order{
							Order: &order.Order{
								Number:     "18",
								Accural:    10.5,
								Status:     "PROCESSED",
								UploadedAt: time.Now().Round(time.Second),
							},
							Owner: "user1",
						},
						"26": orderrecord.Order{
							Order: &order.Order{
								Number:     "26",
								Accural:    10.5,
								Status:     "PROCESSING",
								UploadedAt: time.Now().Round(time.Second),
							},
							Owner: "user1",
						},
						"34": orderrecord.Order{
							Order: &order.Order{
								Number:     "34",
								Accural:    10.5,
								Status:     "NEW",
								UploadedAt: time.Now().Round(time.Second),
							},
							Owner: "user1",
						},
					},
					withdrawn.Withdrawns{},
					map[string]*struct{}{"26": nil, "34": nil},
					map[string][]string{"user1": {"18", "26", "34"}},
					map[string][]string{},
				),
			},
		},

		{
			name: "INVALID ORDER REGISTRATION CORRECT",
			input: testscommon.Input{Method: http.MethodPost, Path: "/api/user/orders", ContentType: "application/json",
				Payload: "42", Account: "user1"},
			expected: testscommon.Expected{
				Code:        202,
				ContentType: "text/plain",
				DbDump: memory.Init(
					account.CustomerAccounts{
						"user1": account.CustomerAccount{
							Login:    "user1",
							Password: "pass1",
							AccountBalance: account.AccountBalance{
								Current:  10.5,
								Withdraw: 0,
							},
						},
						"user2": account.CustomerAccount{
							Login:    "user2",
							Password: "pass2",
							AccountBalance: account.AccountBalance{
								Current:  0,
								Withdraw: 0,
							},
						},
					},
					orderrecord.Orders{
						"18": orderrecord.Order{
							Order: &order.Order{
								Number:     "18",
								Accural:    10.5,
								Status:     "PROCESSED",
								UploadedAt: time.Now().Round(time.Second),
							},
							Owner: "user1",
						},
						"26": orderrecord.Order{
							Order: &order.Order{
								Number:     "26",
								Accural:    10.5,
								Status:     "PROCESSING",
								UploadedAt: time.Now().Round(time.Second),
							},
							Owner: "user1",
						},
						"34": orderrecord.Order{
							Order: &order.Order{
								Number:     "34",
								Accural:    10.5,
								Status:     "NEW",
								UploadedAt: time.Now().Round(time.Second),
							},
							Owner: "user1",
						},
						"42": orderrecord.Order{
							Order: &order.Order{
								Number:     "42",
								Accural:    10.5,
								Status:     "INVALID",
								UploadedAt: time.Now().Round(time.Second),
							},
							Owner: "user1",
						},
					},
					withdrawn.Withdrawns{},
					map[string]*struct{}{"26": nil, "34": nil},
					map[string][]string{"user1": {"18", "26", "34", "42"}},
					map[string][]string{},
				),
			},
		},

		{
			name: "ORDER REAPPLY",
			input: testscommon.Input{Method: http.MethodPost, Path: "/api/user/orders", ContentType: "application/json",
				Payload: "18", Account: "user1"},
			expected: testscommon.Expected{
				Code:        200,
				ContentType: "text/plain",
				DbDump: memory.Init(
					account.CustomerAccounts{
						"user1": account.CustomerAccount{
							Login:    "user1",
							Password: "pass1",
							AccountBalance: account.AccountBalance{
								Current:  10.5,
								Withdraw: 0,
							},
						},
						"user2": account.CustomerAccount{
							Login:    "user2",
							Password: "pass2",
							AccountBalance: account.AccountBalance{
								Current:  0,
								Withdraw: 0,
							},
						},
					},
					orderrecord.Orders{
						"18": orderrecord.Order{
							Order: &order.Order{
								Number:     "18",
								Accural:    10.5,
								Status:     "PROCESSED",
								UploadedAt: time.Now().Round(time.Second),
							},
							Owner: "user1",
						},
					},
					withdrawn.Withdrawns{},
					map[string]*struct{}{},
					map[string][]string{"user1": {"18"}},
					map[string][]string{},
				),
			},
		},

		{
			name: "ORDER REGISTRATION INCORRECT REQUEST FORMAT",
			input: testscommon.Input{Method: http.MethodPost, Path: "/api/user/orders", ContentType: "application/json",
				Payload: "{\"aa\":\"bb\"}", Account: "user1"},
			expected: testscommon.Expected{
				Code:        400,
				ContentType: "text/plain",
				DbDump: memory.Init(
					account.CustomerAccounts{
						"user1": account.CustomerAccount{
							Login:    "user1",
							Password: "pass1",
							AccountBalance: account.AccountBalance{
								Current:  10.5,
								Withdraw: 0,
							},
						},
						"user2": account.CustomerAccount{
							Login:    "user2",
							Password: "pass2",
							AccountBalance: account.AccountBalance{
								Current:  0,
								Withdraw: 0,
							},
						},
					},
					orderrecord.Orders{
						"18": orderrecord.Order{
							Order: &order.Order{
								Number:     "18",
								Accural:    10.5,
								Status:     "PROCESSED",
								UploadedAt: time.Now().Round(time.Second),
							},
							Owner: "user1",
						},
					},
					withdrawn.Withdrawns{},
					map[string]*struct{}{},
					map[string][]string{"user1": {"18"}},
					map[string][]string{},
				),
			},
		},

		{
			name: "ORDER REGISTRATION INCORRECT ORDER ID FORMAT",
			input: testscommon.Input{Method: http.MethodPost, Path: "/api/user/orders", ContentType: "application/json",
				Payload: "20", Account: "user1"},
			expected: testscommon.Expected{
				Code:        422,
				ContentType: "text/plain",
				DbDump: memory.Init(
					account.CustomerAccounts{
						"user1": account.CustomerAccount{
							Login:    "user1",
							Password: "pass1",
							AccountBalance: account.AccountBalance{
								Current:  10.5,
								Withdraw: 0,
							},
						},
						"user2": account.CustomerAccount{
							Login:    "user2",
							Password: "pass2",
							AccountBalance: account.AccountBalance{
								Current:  0,
								Withdraw: 0,
							},
						},
					},
					orderrecord.Orders{
						"1": orderrecord.Order{
							Order: &order.Order{
								Number:     "1",
								Accural:    10.5,
								Status:     "PROCESSED",
								UploadedAt: time.Now().Round(time.Second),
							},
							Owner: "user1",
						},
					},
					withdrawn.Withdrawns{},
					map[string]*struct{}{},
					map[string][]string{"user1": {"1"}},
					map[string][]string{},
				),
			},
		},

		{
			name: "DUBLICATE ORDER",
			input: testscommon.Input{Method: http.MethodPost, Path: "/api/user/orders", ContentType: "application/json",
				Payload: "18", Account: "user2"},
			expected: testscommon.Expected{
				Code:        409,
				ContentType: "text/plain",
				DbDump: memory.Init(
					account.CustomerAccounts{
						"user1": account.CustomerAccount{
							Login:    "user1",
							Password: "pass1",
							AccountBalance: account.AccountBalance{
								Current:  10.5,
								Withdraw: 0,
							},
						},
						"user2": account.CustomerAccount{
							Login:    "user2",
							Password: "pass2",
							AccountBalance: account.AccountBalance{
								Current:  0,
								Withdraw: 0,
							},
						},
					},
					orderrecord.Orders{
						"18": orderrecord.Order{
							Order: &order.Order{
								Number:     "18",
								Accural:    10.5,
								Status:     "PROCESSED",
								UploadedAt: time.Now().Round(time.Second),
							},
							Owner: "user1",
						},
					},
					withdrawn.Withdrawns{},
					map[string]*struct{}{},
					map[string][]string{"user1": {"18"}},
					map[string][]string{},
				),
			},
		},
	}

	//Starting Test AccuralServer
	accuralStorage := accural_storage.Init(message.AccuralMessageMap{
		"18": message.AccuralMessage{
			Order:   "18",
			Status:  "PROCESSED",
			Accural: 10.5,
		},
		"26": message.AccuralMessage{
			Order:   "26",
			Status:  "PROCESSING",
			Accural: 10.5,
		},
		"34": message.AccuralMessage{
			Order:   "34",
			Status:  "NEW",
			Accural: 10.5,
		},
		"42": message.AccuralMessage{
			Order:   "42",
			Status:  "INVALID",
			Accural: 10.5,
		},
	})
	accuralMux := accural_handler.New(accuralStorage)
	accuralMux.Get("/api/orders/{order}", accuralMux.Fetch)
	accuralServer := httptest.NewServer(accuralMux)

	//Starting Test ApplicationServer
	cfg := &config.Config{
		AccuralSystemAddress: accuralServer.URL,
		DatabaseURI:          "",
		RunAddress:           "",
	}
	Auth := auth.New()

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
			"user2": account.CustomerAccount{
				Login:    "user2",
				Password: "pass2",
				AccountBalance: account.AccountBalance{
					Current:  0,
					Withdraw: 0,
				},
			},
		},
		orderrecord.Orders{},
		withdrawn.Withdrawns{},
		map[string]*struct{}{},
		map[string][]string{},
		map[string][]string{},
	)
	mux := handler.New(storage, Auth, cfg)
	mux.Post("/api/user/orders", mux.AddOrder)
	ts := httptest.NewServer(mux)
	defer ts.Close()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			request, err := http.NewRequest(test.input.Method, ts.URL+test.input.Path,
				bytes.NewBuffer([]byte(test.input.Payload)))
			if err != nil {
				fmt.Println(err.Error())
			}
			require.NoError(t, err)
			request.AddCookie(&http.Cookie{Name: "username", Value: test.input.Account, MaxAge: 1000})
			request.Header.Add("Content-Type", test.input.ContentType)
			res, err := http.DefaultClient.Do(request)
			if err != nil {
				fmt.Println(err)
			}
			defer res.Body.Close()
			require.NoError(t, err)
			assert.Equal(t, test.expected.Code, res.StatusCode)
			if res.StatusCode == http.StatusAccepted {
				assert.Equal(t, test.expected.ContentType, res.Header.Get("Content-Type"))
				//fmt.Printf("%+v\n", test.expected.DbDump)
				//fmt.Printf("%+v\n", storage)
				assert.Equal(t, reflect.DeepEqual(*test.expected.DbDump, *storage), true)
			}
		})
	}
}
