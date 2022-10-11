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
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListWithdrawns(t *testing.T) {
	tests := []struct {
		name     string
		input    input
		expected expected
	}{
		{
			name: "WITHDRAW LIST CORRECT",
			input: input{method: http.MethodGet, path: "/api/user/balance/withdrawals", contentType: "text/plain",
				payload: "1", account: "user1"},
			expected: expected{
				payload: "[{\"order\":\"222\",\"sum\":1,\"processed_at\":\"2021-09-19T15:59:42+03:00\"}," +
					"{\"order\":\"111\",\"sum\":1,\"processed_at\":\"2021-09-19T15:59:41+03:00\"}]",
				code:        200,
				contentType: "Application/Json",
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
					UploadedAt: genTimeStamps()[0],
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
				ProcessedAt: genTimeStamps()[0],
			},
			"222": withdrawn.WithdrawnRecord{
				Withdrawn: &withdrawn.Withdrawn{
					Order: "222",
					Sum:   1,
				},
				ProcessedAt: genTimeStamps()[1],
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
			if res.StatusCode == http.StatusOK {
				payload, _ := io.ReadAll(res.Body)
				assert.Equal(t, test.expected.contentType, res.Header.Get("Content-Type"))
				assert.JSONEq(t, test.expected.payload, string(payload))
			}
		})
	}
}

func TestListEmptyWithdrawns(t *testing.T) {
	tests := []struct {
		name     string
		input    input
		expected expected
	}{
		{
			name: "EMPTY WITHDRAW LIST CORRECT",
			input: input{method: http.MethodGet, path: "/api/user/balance/withdrawals", contentType: "text/plain",
				payload: "1", account: "user1"},
			expected: expected{
				payload:     "",
				code:        204,
				contentType: "Application/Json",
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
					UploadedAt: genTimeStamps()[0],
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
		})
	}
}
