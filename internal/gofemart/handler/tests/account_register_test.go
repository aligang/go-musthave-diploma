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
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestRegisterAccount(t *testing.T) {
	tests := []struct {
		name     string
		input    input
		expected expected
	}{
		{
			name: "ACCOUNT REGISTRATION CORRECT",
			input: input{method: http.MethodPost, path: "/api/user/register", contentType: "text/plain",
				payload: "{\"login\":\"user1\",\"password\":\"pass1\"}"},
			expected: expected{
				code:        200,
				contentType: "text/plain",
				dbDump: memory.Init(
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
					internal_order.Orders{},
					withdrawn.Withdrawns{},
					map[string]*struct{}{},
					map[string][]string{},
					map[string][]string{},
				),
			},
		},
		{
			name: "ACCOUNT REGISTRATION INCORRECT INPUT ERROR",
			input: input{method: http.MethodPost, path: "/api/user/register", contentType: "text/plain",
				payload: "{\"login\":\"\",\"password\":\"aa\"}"},
			expected: expected{
				code:        400,
				contentType: "text/plain",
				dbDump: memory.Init(
					customer_account.CustomerAccounts{},
					internal_order.Orders{},
					withdrawn.Withdrawns{},
					map[string]*struct{}{},
					map[string][]string{},
					map[string][]string{},
				),
			},
		},
		{
			name: "DUBLICATE ACCOUNT ERROR",
			input: input{method: http.MethodPost, path: "/api/user/register", contentType: "text/plain",
				payload: "{\"login\":\"user1\",\"password\":\"pass1\"}"},
			expected: expected{
				code:        409,
				contentType: "text/plain",
				dbDump: memory.Init(
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
		customer_account.CustomerAccounts{},
		internal_order.Orders{},
		withdrawn.Withdrawns{},
		map[string]*struct{}{},
		map[string][]string{},
		map[string][]string{},
	)
	mux := handler.New(storage, Auth, nil)
	mux.Post("/api/user/register", mux.RegisterCustomerAccount)
	ts := httptest.NewServer(mux)
	defer ts.Close()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			request, err := http.NewRequest(test.input.method, ts.URL+test.input.path,
				bytes.NewBuffer([]byte(test.input.payload)))
			require.NoError(t, err)
			request.Header.Add("Content-Type", test.input.contentType)
			res, err := http.DefaultClient.Do(request)
			if err != nil {
				fmt.Println(err)
			}
			defer res.Body.Close()
			require.NoError(t, err)
			assert.Equal(t, test.expected.code, res.StatusCode)
			if res.StatusCode == http.StatusOK {
				assert.Equal(t, test.expected.contentType, res.Header.Get("Content-Type"))
				assert.Equal(t, reflect.DeepEqual(*test.expected.dbDump, *storage), true)
			}
		})
	}
}
