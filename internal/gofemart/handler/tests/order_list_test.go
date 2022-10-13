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
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListOrders(t *testing.T) {
	tests := []struct {
		name     string
		input    tests_common.Input
		expected tests_common.Expected
	}{
		{
			name: "ORDER LIST CORRECT",
			input: tests_common.Input{Method: http.MethodGet, Path: "/api/user/orders", ContentType: "text/plain",
				Payload: "1", Account: "user1"},
			expected: tests_common.Expected{
				Payload: "[" +
					"{\"number\":\"3\",\"status\":\"NEW\",\"accural\":20.5,\"uploaded_at\":\"2021-09-19T15:59:43+03:00\"}," +
					"{\"number\":\"2\",\"status\":\"PROCESSING\",\"accural\":25.5,\"uploaded_at\":\"2021-09-19T15:59:42+03:00\"}," +
					"{\"number\":\"1\",\"status\":\"PROCESSED\",\"accural\":10.5,\"uploaded_at\":\"2021-09-19T15:59:41+03:00\"}]",
				Code:        200,
				ContentType: "application/json",
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
								UploadedAt: tests_common.GenTimeStamps()[0],
							},
							Owner: "user1",
						},
						"2": internal_order.Order{
							Order: &order.Order{
								Number:     "2",
								Accural:    10.5,
								Status:     "PROCESSING",
								UploadedAt: tests_common.GenTimeStamps()[1],
							},
							Owner: "user1",
						},
						"3": internal_order.Order{
							Order: &order.Order{
								Number:     "3",
								Accural:    10.5,
								Status:     "NEW",
								UploadedAt: tests_common.GenTimeStamps()[2],
							},
							Owner: "user1",
						},
					},
					withdrawn.Withdrawns{},
					map[string]*struct{}{},
					map[string][]string{"user1": {"1", "2", "3"}},
					map[string][]string{},
				),
			},
		},
	}

	//Starting Test ApplicationServer
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
		internal_order.Orders{
			"1": internal_order.Order{
				Order: &order.Order{
					Number:     "1",
					Accural:    10.5,
					Status:     "PROCESSED",
					UploadedAt: tests_common.GenTimeStamps()[0],
				},
				Owner: "user1",
			},
			"2": internal_order.Order{
				Order: &order.Order{
					Number:     "2",
					Accural:    25.5,
					Status:     "PROCESSING",
					UploadedAt: tests_common.GenTimeStamps()[1],
				},
				Owner: "user1",
			},
			"3": internal_order.Order{
				Order: &order.Order{
					Number:     "3",
					Accural:    20.5,
					Status:     "NEW",
					UploadedAt: tests_common.GenTimeStamps()[2],
				},
				Owner: "user1",
			},
		},
		withdrawn.Withdrawns{},
		map[string]*struct{}{},
		map[string][]string{"user1": {"1", "2", "3"}},
		map[string][]string{},
	)

	mux := handler.New(storage, Auth, nil)
	mux.Get("/api/user/orders", mux.ListOrders)
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
			if res.StatusCode == http.StatusOK {
				payload, _ := io.ReadAll(res.Body)
				fmt.Println(string(payload))
				assert.Equal(t, test.expected.ContentType, res.Header.Get("Content-Type"))
				assert.JSONEq(t, test.expected.Payload, string(payload))
			}
		})
	}
}

func TestListEmptyOrders(t *testing.T) {
	tests := []struct {
		name     string
		input    tests_common.Input
		expected tests_common.Expected
	}{
		{
			name: "EMPTY ORDER LIST CORRECT",
			input: tests_common.Input{Method: http.MethodGet, Path: "/api/user/orders", ContentType: "text/plain",
				Payload: "1", Account: "user1"},
			expected: tests_common.Expected{
				Payload:     "",
				Code:        204,
				ContentType: "application/json",
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
					internal_order.Orders{},
					withdrawn.Withdrawns{},
					map[string]*struct{}{},
					map[string][]string{"user1": {"1", "2", "3"}},
					map[string][]string{},
				),
			},
		},
	}

	//Starting Test ApplicationServer
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
		map[string][]string{"user1": {}},
		map[string][]string{},
	)

	mux := handler.New(storage, Auth, nil)
	mux.Get("/api/user/orders", mux.ListOrders)
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
		})
	}
}
