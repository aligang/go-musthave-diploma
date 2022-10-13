package tests

import (
	"bytes"
	"fmt"
	accural_handler "github.com/aligang/go-musthave-diploma/internal/accural/handler"
	"github.com/aligang/go-musthave-diploma/internal/accural/message"
	accural_storage "github.com/aligang/go-musthave-diploma/internal/accural/storage/memory"
	"github.com/aligang/go-musthave-diploma/internal/config"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/auth"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/customer_account"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/handler"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/order"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/storage/memory"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/storage/memory/internal_order"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/tests_common"
	"github.com/aligang/go-musthave-diploma/internal/withdrawn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestAddOrder(t *testing.T) {
	tests := []struct {
		name     string
		input    tests_common.Input
		expected tests_common.Expected
	}{
		{
			name: "PROCESSED ORDER REGISTRATION CORRECT",
			input: tests_common.Input{Method: http.MethodPost, Path: "/api/user/orders", ContentType: "application/json",
				Payload: "19", Account: "user1"},
			expected: tests_common.Expected{
				Code:        202,
				ContentType: "text/plain",
				DbDump: memory.Init(
					customer_account.CustomerAccounts{
						"user1": customer_account.CustomerAccount{
							Login:    "user1",
							Password: "pass1",
							AccountBalance: customer_account.AccountBalance{
								Balance:  10.5,
								Withdraw: 0,
							},
						},
						"user2": customer_account.CustomerAccount{
							Login:    "user2",
							Password: "pass2",
							AccountBalance: customer_account.AccountBalance{
								Balance:  0,
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
					},
					withdrawn.Withdrawns{},
					map[string]*struct{}{},
					map[string][]string{"user1": {"19"}},
					map[string][]string{},
				),
			},
		},

		{
			name: "PROCESSING ORDER REGISTRATION CORRECT",
			input: tests_common.Input{Method: http.MethodPost, Path: "/api/user/orders", ContentType: "application/json",
				Payload: "24", Account: "user1"},
			expected: tests_common.Expected{
				Code:        202,
				ContentType: "text/plain",
				DbDump: memory.Init(
					customer_account.CustomerAccounts{
						"user1": customer_account.CustomerAccount{
							Login:    "user1",
							Password: "pass1",
							AccountBalance: customer_account.AccountBalance{
								Balance:  10.5,
								Withdraw: 0,
							},
						},
						"user2": customer_account.CustomerAccount{
							Login:    "user2",
							Password: "pass2",
							AccountBalance: customer_account.AccountBalance{
								Balance:  0,
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
						"24": internal_order.Order{
							Order: &order.Order{
								Number:     "24",
								Accural:    10.5,
								Status:     "PROCESSING",
								UploadedAt: time.Now().Round(time.Second),
							},
							Owner: "user1",
						},
					},
					withdrawn.Withdrawns{},
					map[string]*struct{}{"24": nil},
					map[string][]string{"user1": {"19", "24"}},
					map[string][]string{},
				),
			},
		},

		{
			name: "NEW ORDER REGISTRATION CORRECT",
			input: tests_common.Input{Method: http.MethodPost, Path: "/api/user/orders", ContentType: "application/json",
				Payload: "38", Account: "user1"},
			expected: tests_common.Expected{
				Code:        202,
				ContentType: "text/plain",
				DbDump: memory.Init(
					customer_account.CustomerAccounts{
						"user1": customer_account.CustomerAccount{
							Login:    "user1",
							Password: "pass1",
							AccountBalance: customer_account.AccountBalance{
								Balance:  10.5,
								Withdraw: 0,
							},
						},
						"user2": customer_account.CustomerAccount{
							Login:    "user2",
							Password: "pass2",
							AccountBalance: customer_account.AccountBalance{
								Balance:  0,
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
						"24": internal_order.Order{
							Order: &order.Order{
								Number:     "24",
								Accural:    10.5,
								Status:     "PROCESSING",
								UploadedAt: time.Now().Round(time.Second),
							},
							Owner: "user1",
						},
						"38": internal_order.Order{
							Order: &order.Order{
								Number:     "38",
								Accural:    10.5,
								Status:     "NEW",
								UploadedAt: time.Now().Round(time.Second),
							},
							Owner: "user1",
						},
					},
					withdrawn.Withdrawns{},
					map[string]*struct{}{"24": nil, "38": nil},
					map[string][]string{"user1": {"19", "24", "38"}},
					map[string][]string{},
				),
			},
		},

		{
			name: "INVALID ORDER REGISTRATION CORRECT",
			input: tests_common.Input{Method: http.MethodPost, Path: "/api/user/orders", ContentType: "application/json",
				Payload: "43", Account: "user1"},
			expected: tests_common.Expected{
				Code:        202,
				ContentType: "text/plain",
				DbDump: memory.Init(
					customer_account.CustomerAccounts{
						"user1": customer_account.CustomerAccount{
							Login:    "user1",
							Password: "pass1",
							AccountBalance: customer_account.AccountBalance{
								Balance:  10.5,
								Withdraw: 0,
							},
						},
						"user2": customer_account.CustomerAccount{
							Login:    "user2",
							Password: "pass2",
							AccountBalance: customer_account.AccountBalance{
								Balance:  0,
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
						"24": internal_order.Order{
							Order: &order.Order{
								Number:     "24",
								Accural:    10.5,
								Status:     "PROCESSING",
								UploadedAt: time.Now().Round(time.Second),
							},
							Owner: "user1",
						},
						"38": internal_order.Order{
							Order: &order.Order{
								Number:     "38",
								Accural:    10.5,
								Status:     "NEW",
								UploadedAt: time.Now().Round(time.Second),
							},
							Owner: "user1",
						},
						"43": internal_order.Order{
							Order: &order.Order{
								Number:     "43",
								Accural:    10.5,
								Status:     "INVALID",
								UploadedAt: time.Now().Round(time.Second),
							},
							Owner: "user1",
						},
					},
					withdrawn.Withdrawns{},
					map[string]*struct{}{"24": nil, "38": nil},
					map[string][]string{"user1": {"19", "24", "38", "43"}},
					map[string][]string{},
				),
			},
		},

		{
			name: "ORDER REAPPLY",
			input: tests_common.Input{Method: http.MethodPost, Path: "/api/user/orders", ContentType: "application/json",
				Payload: "19", Account: "user1"},
			expected: tests_common.Expected{
				Code:        200,
				ContentType: "text/plain",
				DbDump: memory.Init(
					customer_account.CustomerAccounts{
						"user1": customer_account.CustomerAccount{
							Login:    "user1",
							Password: "pass1",
							AccountBalance: customer_account.AccountBalance{
								Balance:  10.5,
								Withdraw: 0,
							},
						},
						"user2": customer_account.CustomerAccount{
							Login:    "user2",
							Password: "pass2",
							AccountBalance: customer_account.AccountBalance{
								Balance:  0,
								Withdraw: 0,
							},
						},
					},
					internal_order.Orders{
						"1": internal_order.Order{
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
			name: "ORDER REGISTRATION INCORRECT REQUEST FORMAT",
			input: tests_common.Input{Method: http.MethodPost, Path: "/api/user/orders", ContentType: "application/json",
				Payload: "{\"aa\":\"bb\"}", Account: "user1"},
			expected: tests_common.Expected{
				Code:        400,
				ContentType: "text/plain",
				DbDump: memory.Init(
					customer_account.CustomerAccounts{
						"user1": customer_account.CustomerAccount{
							Login:    "user1",
							Password: "pass1",
							AccountBalance: customer_account.AccountBalance{
								Balance:  10.5,
								Withdraw: 0,
							},
						},
						"user2": customer_account.CustomerAccount{
							Login:    "user2",
							Password: "pass2",
							AccountBalance: customer_account.AccountBalance{
								Balance:  0,
								Withdraw: 0,
							},
						},
					},
					internal_order.Orders{
						"1": internal_order.Order{
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
			name: "ORDER REGISTRATION INCORRECT ORDER ID FORMAT",
			input: tests_common.Input{Method: http.MethodPost, Path: "/api/user/orders", ContentType: "application/json",
				Payload: "20", Account: "user1"},
			expected: tests_common.Expected{
				Code:        422,
				ContentType: "text/plain",
				DbDump: memory.Init(
					customer_account.CustomerAccounts{
						"user1": customer_account.CustomerAccount{
							Login:    "user1",
							Password: "pass1",
							AccountBalance: customer_account.AccountBalance{
								Balance:  10.5,
								Withdraw: 0,
							},
						},
						"user2": customer_account.CustomerAccount{
							Login:    "user2",
							Password: "pass2",
							AccountBalance: customer_account.AccountBalance{
								Balance:  0,
								Withdraw: 0,
							},
						},
					},
					internal_order.Orders{
						"1": internal_order.Order{
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
			input: tests_common.Input{Method: http.MethodPost, Path: "/api/user/orders", ContentType: "application/json",
				Payload: "19", Account: "user2"},
			expected: tests_common.Expected{
				Code:        409,
				ContentType: "text/plain",
				DbDump: memory.Init(
					customer_account.CustomerAccounts{
						"user1": customer_account.CustomerAccount{
							Login:    "user1",
							Password: "pass1",
							AccountBalance: customer_account.AccountBalance{
								Balance:  10.5,
								Withdraw: 0,
							},
						},
						"user2": customer_account.CustomerAccount{
							Login:    "user2",
							Password: "pass2",
							AccountBalance: customer_account.AccountBalance{
								Balance:  0,
								Withdraw: 0,
							},
						},
					},
					internal_order.Orders{
						"1": internal_order.Order{
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
	}

	//Starting Test AccuralServer
	accuralStorage := accural_storage.Init(message.AccuralMessageMap{
		"19": message.AccuralMessage{
			Order:   "19",
			Status:  "PROCESSED",
			Accural: 10.5,
		},
		"24": message.AccuralMessage{
			Order:   "24",
			Status:  "PROCESSING",
			Accural: 10.5,
		},
		"38": message.AccuralMessage{
			Order:   "38",
			Status:  "NEW",
			Accural: 10.5,
		},
		"43": message.AccuralMessage{
			Order:   "43",
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
	Auth := auth.New()
	storage := memory.Init(
		customer_account.CustomerAccounts{
			"user1": customer_account.CustomerAccount{
				Login:    "user1",
				Password: "pass1",
				AccountBalance: customer_account.AccountBalance{
					Balance:  0,
					Withdraw: 0,
				},
			},
			"user2": customer_account.CustomerAccount{
				Login:    "user2",
				Password: "pass2",
				AccountBalance: customer_account.AccountBalance{
					Balance:  0,
					Withdraw: 0,
				},
			},
		},
		internal_order.Orders{},
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
