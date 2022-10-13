package tests

import (
	"bytes"
	"fmt"
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
	"testing"
	"time"
)

func TestAddWithdraw(t *testing.T) {
	tests := []struct {
		name     string
		input    tests_common.Input
		expected tests_common.Expected
	}{
		{
			name: "WITHDRAW REGISTRATION",
			input: tests_common.Input{Method: http.MethodPost, Path: "/api/user/balance/withdrawn", ContentType: "application/json",
				Payload: "{\"order\": \"57\", \"sum\": 1}", Account: "user1"},
			expected: tests_common.Expected{
				Code:        200,
				ContentType: "text/plain",
				DbDump: memory.Init(
					customer_account.CustomerAccounts{
						"user1": customer_account.CustomerAccount{
							Login:    "user1",
							Password: "pass1",
							AccountBalance: customer_account.AccountBalance{
								Balance:  99,
								Withdraw: 0,
							},
						},
					},
					internal_order.Orders{
						"1": internal_order.Order{
							Order: &order.Order{
								Number:     "1",
								Accural:    100,
								Status:     "PROCESSED",
								UploadedAt: time.Now().Round(time.Second),
							},
							Owner: "user1",
						},
					},
					withdrawn.Withdrawns{
						"57": withdrawn.WithdrawnRecord{
							Withdrawn: &withdrawn.Withdrawn{
								Order: "57",
								Sum:   1,
							},
							ProcessedAt: time.Now().Round(time.Second),
						},
					},
					map[string]*struct{}{},
					map[string][]string{"user1": {"1"}},
					map[string][]string{"user1": {"57"}},
				),
			},
		},
		{
			name: "WITHDRAW INSUFFICIENT FOUNDS",
			input: tests_common.Input{Method: http.MethodPost, Path: "/api/user/balance/withdrawn", ContentType: "application/json",
				Payload: "{\"order\": \"62\", \"sum\": 1000}", Account: "user1"},
			expected: tests_common.Expected{
				Code:        402,
				ContentType: "text/plain",
				DbDump: memory.Init(
					customer_account.CustomerAccounts{
						"user1": customer_account.CustomerAccount{
							Login:    "user1",
							Password: "pass1",
							AccountBalance: customer_account.AccountBalance{
								Balance:  99,
								Withdraw: 0,
							},
						},
					},
					internal_order.Orders{
						"1": internal_order.Order{
							Order: &order.Order{
								Number:     "1",
								Accural:    100,
								Status:     "PROCESSED",
								UploadedAt: time.Now().Round(time.Second),
							},
							Owner: "user1",
						},
					},
					withdrawn.Withdrawns{
						"62": withdrawn.WithdrawnRecord{
							Withdrawn: &withdrawn.Withdrawn{
								Order: "62",
								Sum:   1,
							},
							ProcessedAt: time.Now().Round(time.Second),
						},
					},
					map[string]*struct{}{},
					map[string][]string{"user1": {"1"}},
					map[string][]string{"user1": {"62"}},
				),
			},
		},
		{
			name: "WITHDRAW INCORRECT ID FORMAT",
			input: tests_common.Input{Method: http.MethodPost, Path: "/api/user/balance/withdrawn", ContentType: "application/json",
				Payload: "{\"order\": \"65\", \"sum\": 1}", Account: "user1"},
			expected: tests_common.Expected{
				Code:        422,
				ContentType: "text/plain",
				DbDump: memory.Init(
					customer_account.CustomerAccounts{
						"user1": customer_account.CustomerAccount{
							Login:    "user1",
							Password: "pass1",
							AccountBalance: customer_account.AccountBalance{
								Balance:  99,
								Withdraw: 0,
							},
						},
					},
					internal_order.Orders{
						"1": internal_order.Order{
							Order: &order.Order{
								Number:     "1",
								Accural:    100,
								Status:     "PROCESSED",
								UploadedAt: time.Now().Round(time.Second),
							},
							Owner: "user1",
						},
					},
					withdrawn.Withdrawns{
						"111": withdrawn.WithdrawnRecord{
							Withdrawn: &withdrawn.Withdrawn{
								Order: "111",
								Sum:   1,
							},
							ProcessedAt: time.Now().Round(time.Second),
						},
					},
					map[string]*struct{}{},
					map[string][]string{"user1": {"1"}},
					map[string][]string{"user1": {"111"}},
				),
			},
		},
	}

	Auth := auth.New()
	storage := memory.Init(
		customer_account.CustomerAccounts{
			"user1": customer_account.CustomerAccount{
				Login:    "user1",
				Password: "pass1",
				AccountBalance: customer_account.AccountBalance{
					Balance:  100,
					Withdraw: 0,
				},
			},
		},
		internal_order.Orders{
			"1": internal_order.Order{
				Order: &order.Order{
					Number:     "1",
					Accural:    100,
					Status:     "PROCESSED",
					UploadedAt: time.Now().Round(time.Second),
				},
				Owner: "user1",
			},
		},
		withdrawn.Withdrawns{},
		map[string]*struct{}{},
		map[string][]string{},
		map[string][]string{},
	)

	mux := handler.New(storage, Auth, nil)
	mux.Post("/api/user/balance/withdrawn", mux.AddWithdraw)
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
