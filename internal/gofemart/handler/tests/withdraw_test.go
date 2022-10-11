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
		input    input
		expected expected
	}{
		{
			name: "WITHDRAW REGISTRATION",
			input: input{method: http.MethodPost, path: "/api/user/balance/withdrawn", contentType: "application/json",
				payload: "{\"order\": \"57\", \"sum\": 1}", account: "user1"},
			expected: expected{
				code:        200,
				contentType: "text/plain",
				dbDump: memory.Init(
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
			input: input{method: http.MethodPost, path: "/api/user/balance/withdrawn", contentType: "application/json",
				payload: "{\"order\": \"62\", \"sum\": 1000}", account: "user1"},
			expected: expected{
				code:        402,
				contentType: "text/plain",
				dbDump: memory.Init(
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
			input: input{method: http.MethodPost, path: "/api/user/balance/withdrawn", contentType: "application/json",
				payload: "{\"order\": \"65\", \"sum\": 1}", account: "user1"},
			expected: expected{
				code:        422,
				contentType: "text/plain",
				dbDump: memory.Init(
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

			request, err := http.NewRequest(test.input.method, ts.URL+test.input.path,
				bytes.NewBuffer([]byte(test.input.payload)))
			require.NoError(t, err)
			request.AddCookie(&http.Cookie{Name: "username", Value: test.input.account, MaxAge: 1000})
			request.Header.Add("Content-Type", test.input.contentType)
			res, err := http.DefaultClient.Do(request)
			if err != nil {
				fmt.Println(err)
			}
			defer res.Body.Close()
			require.NoError(t, err)
			assert.Equal(t, test.expected.code, res.StatusCode)
			if res.StatusCode == http.StatusAccepted {
				assert.Equal(t, test.expected.contentType, res.Header.Get("Content-Type"))
				//fmt.Printf("%+v\n", test.expected.dbDump)
				//fmt.Printf("%+v\n", storage)
				assert.Equal(t, reflect.DeepEqual(*test.expected.dbDump, *storage), true)
			}
		})
	}
}
