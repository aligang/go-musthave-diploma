package tests

import (
	"bytes"
	"fmt"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/account"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/auth"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/handler"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/storage/memory"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/storage/memory/orderrecord"
	"github.com/aligang/go-musthave-diploma/internal/gofemart/testscommon"
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
		input    testscommon.Input
		expected testscommon.Expected
	}{
		{
			name: "ACCOUNT REGISTRATION CORRECT",
			input: testscommon.Input{Method: http.MethodPost, Path: "/api/user/register", ContentType: "text/plain",
				Payload: "{\"login\":\"user1\",\"password\":\"pass1\"}"},
			expected: testscommon.Expected{
				Code:        200,
				ContentType: "text/plain",
				DBDump: memory.Init(
					account.CustomerAccounts{
						"user1": account.CustomerAccount{
							Login:    "user1",
							Password: "pass1",
							AccountBalance: account.AccountBalance{
								Current:  0,
								Withdraw: 0,
							},
						},
					},
					orderrecord.Orders{},
					withdrawn.Withdrawns{},
					map[string]*struct{}{},
					map[string][]string{},
					map[string][]string{},
				),
			},
		},
		{
			name: "ACCOUNT REGISTRATION INCORRECT testscommon.Input ERROR",
			input: testscommon.Input{Method: http.MethodPost, Path: "/api/user/register", ContentType: "text/plain",
				Payload: "{\"login\":\"\",\"password\":\"aa\"}"},
			expected: testscommon.Expected{
				Code:        400,
				ContentType: "text/plain",
				DBDump: memory.Init(
					account.CustomerAccounts{},
					orderrecord.Orders{},
					withdrawn.Withdrawns{},
					map[string]*struct{}{},
					map[string][]string{},
					map[string][]string{},
				),
			},
		},
		{
			name: "DUBLICATE ACCOUNT ERROR",
			input: testscommon.Input{Method: http.MethodPost, Path: "/api/user/register", ContentType: "text/plain",
				Payload: "{\"login\":\"user1\",\"password\":\"pass1\"}"},
			expected: testscommon.Expected{
				Code:        409,
				ContentType: "text/plain",
				DBDump: memory.Init(
					account.CustomerAccounts{
						"user1": account.CustomerAccount{
							Login:    "user1",
							Password: "pass1",
							AccountBalance: account.AccountBalance{
								Current:  0,
								Withdraw: 0,
							},
						},
					},
					orderrecord.Orders{},
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
		account.CustomerAccounts{},
		orderrecord.Orders{},
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
				assert.Equal(t, reflect.DeepEqual(*test.expected.DBDump, *storage), true)
			}
		})
	}
}
