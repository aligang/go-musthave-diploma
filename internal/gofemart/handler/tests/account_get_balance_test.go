package tests

import (
	"bytes"
	"fmt"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/auth"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/customer_account"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/handler"
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

func TestGetBalanceAccount(t *testing.T) {
	tests := []struct {
		name     string
		input    input
		expected expected
	}{
		{
			name: "GET BALANCE CORRECT",
			input: input{method: http.MethodGet, path: "/api/user/balance", contentType: "application/json",
				payload: "", account: "user1"},
			expected: expected{
				code:        200,
				contentType: "Application/Json",
				payload:     "{\"balance\":100.5,\"withdrawn\":200.9}",
				dbDump: memory.Init(
					customer_account.CustomerAccounts{
						"user1": customer_account.CustomerAccount{
							Login:    "user1",
							Password: "pass1",
							AccountBalance: customer_account.AccountBalance{
								Balance:  100.5,
								Withdraw: 200.9,
							},
						},
					},
					internal_order.Orders{},
					withdrawn.Withdrawns{},
					map[string]*struct{}{},
					map[string][]string{},
					map[string][]string{},
				),
			},
		},
	}

	//accuralStorage := accuralMemory.Init(message.AccuralMessageMap{})
	Auth := auth.New()
	storage := memory.Init(
		customer_account.CustomerAccounts{
			"user1": customer_account.CustomerAccount{
				Login:    "user1",
				Password: "pass1",
				AccountBalance: customer_account.AccountBalance{
					Balance:  100.5,
					Withdraw: 200.9,
				},
			},
		},
		internal_order.Orders{},
		withdrawn.Withdrawns{},
		map[string]*struct{}{},
		map[string][]string{},
		map[string][]string{},
	)
	mux := handler.New(storage, Auth, nil)
	mux.Get("/api/user/balance", mux.GetAccountBalance)
	ts := httptest.NewServer(mux)
	defer ts.Close()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			request, err := http.NewRequest(test.input.method, ts.URL+test.input.path,
				bytes.NewBuffer([]byte(test.input.payload)))
			require.NoError(t, err)
			request.Header.Add("Content-Type", test.input.contentType)
			request.AddCookie(
				&http.Cookie{Name: "username", Value: test.input.account, MaxAge: 1000})
			res, err := http.DefaultClient.Do(request)
			if err != nil {
				fmt.Println(err)
			}
			defer res.Body.Close()
			require.NoError(t, err)
			assert.Equal(t, test.expected.code, res.StatusCode)
			if res.StatusCode == http.StatusOK {
				assert.Equal(t, test.expected.contentType, res.Header.Get("Content-Type"))
				payload, _ := io.ReadAll(res.Body)
				assert.JSONEq(t, test.expected.payload, string(payload))
			}
		})
	}

}
