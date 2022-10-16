package tests

import (
	"bytes"
	"fmt"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/account"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/auth"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/handler"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/storage/memory"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/storage/memory/orderRecord"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/tests_common"
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
		input    tests_common.Input
		expected tests_common.Expected
	}{
		{
			name: "GET BALANCE CORRECT",
			input: tests_common.Input{Method: http.MethodGet, Path: "/api/user/balance", ContentType: "application/json",
				Payload: "", Account: "user1"},
			expected: tests_common.Expected{
				Code:        200,
				ContentType: "application/json",
				Payload:     "{\"current\":100.5,\"withdrawn\":200.9}",
				DbDump: memory.Init(
					account.CustomerAccounts{
						"user1": account.CustomerAccount{
							Login:    "user1",
							Password: "pass1",
							AccountBalance: account.AccountBalance{
								Current:  100.5,
								Withdraw: 200.9,
							},
						},
					},
					orderRecord.Orders{},
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
		account.CustomerAccounts{
			"user1": account.CustomerAccount{
				Login:    "user1",
				Password: "pass1",
				AccountBalance: account.AccountBalance{
					Current:  100.5,
					Withdraw: 200.9,
				},
			},
		},
		orderRecord.Orders{},
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

			request, err := http.NewRequest(test.input.Method, ts.URL+test.input.Path,
				bytes.NewBuffer([]byte(test.input.Payload)))
			require.NoError(t, err)
			request.Header.Add("Content-Type", test.input.ContentType)
			request.AddCookie(
				&http.Cookie{Name: "username", Value: test.input.Account, MaxAge: 1000})
			res, err := http.DefaultClient.Do(request)
			if err != nil {
				fmt.Println(err)
			}
			defer res.Body.Close()
			require.NoError(t, err)
			assert.Equal(t, test.expected.Code, res.StatusCode)
			if res.StatusCode == http.StatusOK {
				assert.Equal(t, test.expected.ContentType, res.Header.Get("Content-Type"))
				payload, _ := io.ReadAll(res.Body)
				assert.JSONEq(t, test.expected.Payload, string(payload))
			}
		})
	}

}
