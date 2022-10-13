package tests

import (
	"bytes"
	"fmt"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/auth"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/customer_account"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/handler"
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
)

func TestRegisterAccount(t *testing.T) {
	tests := []struct {
		name     string
		input    tests_common.Input
		expected tests_common.Expected
	}{
		{
			name: "ACCOUNT REGISTRATION CORRECT",
			input: tests_common.Input{Method: http.MethodPost, Path: "/api/user/register", ContentType: "text/plain",
				Payload: "{\"login\":\"user1\",\"password\":\"pass1\"}"},
			expected: tests_common.Expected{
				Code:        200,
				ContentType: "text/plain",
				DbDump: memory.Init(
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
			name: "ACCOUNT REGISTRATION INCORRECT tests_common.Input ERROR",
			input: tests_common.Input{Method: http.MethodPost, Path: "/api/user/register", ContentType: "text/plain",
				Payload: "{\"login\":\"\",\"password\":\"aa\"}"},
			expected: tests_common.Expected{
				Code:        400,
				ContentType: "text/plain",
				DbDump: memory.Init(
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
			input: tests_common.Input{Method: http.MethodPost, Path: "/api/user/register", ContentType: "text/plain",
				Payload: "{\"login\":\"user1\",\"password\":\"pass1\"}"},
			expected: tests_common.Expected{
				Code:        409,
				ContentType: "text/plain",
				DbDump: memory.Init(
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

			request, err := http.NewRequest(test.input.Method, ts.URL+test.input.Path,
				bytes.NewBuffer([]byte(test.input.Payload)))
			require.NoError(t, err)
			request.Header.Add("Content-Type", test.input.ContentType)
			res, err := http.DefaultClient.Do(request)
			if err != nil {
				fmt.Println(err)
			}
			defer res.Body.Close()
			require.NoError(t, err)
			assert.Equal(t, test.expected.Code, res.StatusCode)
			if res.StatusCode == http.StatusOK {
				assert.Equal(t, test.expected.ContentType, res.Header.Get("Content-Type"))
				assert.Equal(t, reflect.DeepEqual(*test.expected.DbDump, *storage), true)
			}
		})
	}
}
