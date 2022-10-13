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

func TestListWithdrawns(t *testing.T) {
	tests := []struct {
		name     string
		input    tests_common.Input
		expected tests_common.Expected
	}{
		{
			name: "WITHDRAW LIST CORRECT",
			input: tests_common.Input{Method: http.MethodGet, Path: "/api/user/balance/withdrawals", ContentType: "text/plain",
				Payload: "1", Account: "user1"},
			expected: tests_common.Expected{
				Payload: "[{\"order\":\"222\",\"sum\":1,\"processed_at\":\"2021-09-19T15:59:42+03:00\"}," +
					"{\"order\":\"111\",\"sum\":1,\"processed_at\":\"2021-09-19T15:59:41+03:00\"}]",
				Code:        200,
				ContentType: "application/json",
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
					Balance:  10.5,
					Withdraw: 2,
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
		},
		withdrawn.Withdrawns{
			"111": withdrawn.WithdrawnRecord{
				Withdrawn: &withdrawn.Withdrawn{
					Order: "111",
					Sum:   1,
				},
				ProcessedAt: tests_common.GenTimeStamps()[0],
			},
			"222": withdrawn.WithdrawnRecord{
				Withdrawn: &withdrawn.Withdrawn{
					Order: "222",
					Sum:   1,
				},
				ProcessedAt: tests_common.GenTimeStamps()[1],
			},
		},
		map[string]*struct{}{},
		map[string][]string{"user1": {"1", "2", "3"}},
		map[string][]string{"user1": {"111", "222"}},
	)

	mux := handler.New(storage, Auth, nil)
	mux.Get("/api/user/balance/withdrawals", mux.ListWithdraws)
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
				assert.Equal(t, test.expected.ContentType, res.Header.Get("Content-Type"))
				assert.JSONEq(t, test.expected.Payload, string(payload))
			}
		})
	}
}

func TestListEmptyWithdrawns(t *testing.T) {
	tests := []struct {
		name     string
		input    tests_common.Input
		expected tests_common.Expected
	}{
		{
			name: "EMPTY WITHDRAW LIST CORRECT",
			input: tests_common.Input{Method: http.MethodGet, Path: "/api/user/balance/withdrawals", ContentType: "text/plain",
				Payload: "1", Account: "user1"},
			expected: tests_common.Expected{
				Payload:     "",
				Code:        204,
				ContentType: "application/json",
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
		},
		withdrawn.Withdrawns{},
		map[string]*struct{}{},
		map[string][]string{"user1": {}},
		map[string][]string{},
	)

	mux := handler.New(storage, Auth, nil)
	mux.Get("/api/user/balance/withdrawals", mux.ListWithdraws)
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
