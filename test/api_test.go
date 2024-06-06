package test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	ApiUrl         = "http://localhost:8000"
	PrefixUsername = "1."
)

func TestApi(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip API tests")
	}

	testcases := getTestCases()
	ctx := context.Background()
	client := &http.Client{}

	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			for idx := range tc.Steps {
				step := &tc.Steps[idx]
				request, err := step.Request(t, ctx, &tc)
				request.Header.Set("Content-Type", "application/json")
				request.Header.Set("Accept", "application/json")
				require.NoError(t, err)

				// Send request
				response, err := client.Do(request)

				require.NoError(t, err)
				defer response.Body.Close()

				// Check response
				ReadJsonResult(t, response, step)
				step.Expect(t, ctx, &tc, response, step.Result, step.ResultArray)
			}
		})
	}
}

func getTestCases() []TestCase {
	return []TestCase{
		{
			Name: "Test 1",
			Steps: []TestCaseStep{
				{
					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
						return http.NewRequest(http.MethodPost, ApiUrl+"/create_user", bytes.NewBufferString(`{"username":"`+PrefixUsername+`username"}`))
					},
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any, data2 []map[string]any) {
						require.Equal(t, http.StatusCreated, resp.StatusCode)
						require.NotEmpty(t, data["token"].(string))
					},
				},
			},
		},
		{
			Name: "Test 2",
			Steps: []TestCaseStep{
				{
					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
						return http.NewRequest(http.MethodPost, ApiUrl+"/create_user", bytes.NewBufferString(`{"username":"`+PrefixUsername+`username1"}`))
					},
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any, data2 []map[string]any) {
						require.Equal(t, http.StatusCreated, resp.StatusCode)
						require.NotEmpty(t, data["token"].(string))
					},
				},
				{
					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
						req, err := http.NewRequest(http.MethodGet, ApiUrl+"/balance_read", nil)
						req.Header.Set("Authorization", "Bearer "+tc.Steps[0].Result["token"].(string))
						return req, err
					},
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any, data2 []map[string]any) {
						require.Equal(t, http.StatusOK, resp.StatusCode)
						require.Equal(t, float64(0), data["balance"].(float64))
					},
				},
			},
		},
		{
			Name: "Test 3",
			Steps: []TestCaseStep{
				{
					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
						return http.NewRequest(http.MethodPost, ApiUrl+"/create_user", bytes.NewBufferString(`{"username":"`+PrefixUsername+`username3"}`))
					},
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any, data2 []map[string]any) {
						require.Equal(t, http.StatusCreated, resp.StatusCode)
						require.NotEmpty(t, data["token"].(string))
					},
				},
				{
					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
						req, err := http.NewRequest(http.MethodPost, ApiUrl+"/balance_topup", bytes.NewBufferString(`{"amount":100000}`))
						req.Header.Set("Authorization", "Bearer "+tc.Steps[0].Result["token"].(string))
						return req, err
					},
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any, data2 []map[string]any) {
						require.Equal(t, http.StatusNoContent, resp.StatusCode)
					},
				},
				{
					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
						req, err := http.NewRequest(http.MethodPost, ApiUrl+"/balance_topup", bytes.NewBufferString(`{"amount":50000}`))
						req.Header.Set("Authorization", "Bearer "+tc.Steps[0].Result["token"].(string))
						return req, err
					},
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any, data2 []map[string]any) {
						require.Equal(t, http.StatusNoContent, resp.StatusCode)
					},
				},
				{
					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
						req, err := http.NewRequest(http.MethodGet, ApiUrl+"/balance_read", nil)
						req.Header.Set("Authorization", "Bearer "+tc.Steps[0].Result["token"].(string))
						return req, err
					},
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any, data2 []map[string]any) {
						require.Equal(t, http.StatusOK, resp.StatusCode)
						require.Equal(t, float64(150000), data["balance"].(float64))
					},
				},
			},
		},
		{
			Name: "Test 4",
			Steps: []TestCaseStep{
				{
					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
						return http.NewRequest(http.MethodPost, ApiUrl+"/create_user", bytes.NewBufferString(`{"username":"`+PrefixUsername+`username4"}`))
					},
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any, data2 []map[string]any) {
						require.Equal(t, http.StatusCreated, resp.StatusCode)
						require.NotEmpty(t, data["token"].(string))
					},
				},
				{
					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
						return http.NewRequest(http.MethodPost, ApiUrl+"/create_user", bytes.NewBufferString(`{"username":"`+PrefixUsername+`username5"}`))
					},
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any, data2 []map[string]any) {
						require.Equal(t, http.StatusCreated, resp.StatusCode)
						require.NotEmpty(t, data["token"].(string))
					},
				},
				{
					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
						req, err := http.NewRequest(http.MethodPost, ApiUrl+"/balance_topup", bytes.NewBufferString(`{"amount":100000}`))
						req.Header.Set("Authorization", "Bearer "+tc.Steps[0].Result["token"].(string))
						return req, err
					},
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any, data2 []map[string]any) {
						require.Equal(t, http.StatusNoContent, resp.StatusCode)
					},
				},
				{
					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
						req, err := http.NewRequest(http.MethodPost, ApiUrl+"/balance_topup", bytes.NewBufferString(`{"amount":50000}`))
						req.Header.Set("Authorization", "Bearer "+tc.Steps[1].Result["token"].(string))
						return req, err
					},
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any, data2 []map[string]any) {
						require.Equal(t, http.StatusNoContent, resp.StatusCode)
					},
				},
				{
					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
						req, err := http.NewRequest(http.MethodGet, ApiUrl+"/balance_read", nil)
						req.Header.Set("Authorization", "Bearer "+tc.Steps[0].Result["token"].(string))
						return req, err
					},
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any, data2 []map[string]any) {
						require.Equal(t, http.StatusOK, resp.StatusCode)
						require.Equal(t, float64(100000), data["balance"].(float64))
					},
				},
				{
					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
						req, err := http.NewRequest(http.MethodGet, ApiUrl+"/balance_read", nil)
						req.Header.Set("Authorization", "Bearer "+tc.Steps[1].Result["token"].(string))
						return req, err
					},
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any, data2 []map[string]any) {
						require.Equal(t, http.StatusOK, resp.StatusCode)
						require.Equal(t, float64(50000), data["balance"].(float64))
					},
				},
				{
					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
						req, err := http.NewRequest(http.MethodPost, ApiUrl+"/transfer", bytes.NewBufferString(`{"amount":25000,"to_username":"`+PrefixUsername+`username4"}`))
						req.Header.Set("Authorization", "Bearer "+tc.Steps[1].Result["token"].(string))
						return req, err
					},
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any, data2 []map[string]any) {
						require.Equal(t, http.StatusNoContent, resp.StatusCode)
					},
				},
				{
					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
						req, err := http.NewRequest(http.MethodGet, ApiUrl+"/balance_read", nil)
						req.Header.Set("Authorization", "Bearer "+tc.Steps[0].Result["token"].(string))
						return req, err
					},
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any, data2 []map[string]any) {
						require.Equal(t, http.StatusOK, resp.StatusCode)
						require.Equal(t, float64(125000), data["balance"].(float64))
					},
				},
				{
					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
						req, err := http.NewRequest(http.MethodGet, ApiUrl+"/balance_read", nil)
						req.Header.Set("Authorization", "Bearer "+tc.Steps[1].Result["token"].(string))
						return req, err
					},
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any, data2 []map[string]any) {
						require.Equal(t, http.StatusOK, resp.StatusCode)
						require.Equal(t, float64(25000), data["balance"].(float64))
					},
				},
			},
		},
		{
			Name: "Test 5",
			Steps: []TestCaseStep{
				{
					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
						return http.NewRequest(http.MethodPost, ApiUrl+"/create_user", bytes.NewBufferString(`{"username":"`+PrefixUsername+`username6"}`))
					},
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any, data2 []map[string]any) {
						require.Equal(t, http.StatusCreated, resp.StatusCode)
						require.NotEmpty(t, data["token"].(string))
					},
				},
				{
					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
						return http.NewRequest(http.MethodPost, ApiUrl+"/create_user", bytes.NewBufferString(`{"username":"`+PrefixUsername+`username7"}`))
					},
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any, data2 []map[string]any) {
						require.Equal(t, http.StatusCreated, resp.StatusCode)
						require.NotEmpty(t, data["token"].(string))
					},
				},
				{
					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
						req, err := http.NewRequest(http.MethodPost, ApiUrl+"/balance_topup", bytes.NewBufferString(`{"amount":100000}`))
						req.Header.Set("Authorization", "Bearer "+tc.Steps[0].Result["token"].(string))
						return req, err
					},
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any, data2 []map[string]any) {
						require.Equal(t, http.StatusNoContent, resp.StatusCode)
					},
				},
				{
					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
						req, err := http.NewRequest(http.MethodPost, ApiUrl+"/balance_topup", bytes.NewBufferString(`{"amount":100000}`))
						req.Header.Set("Authorization", "Bearer "+tc.Steps[0].Result["token"].(string))
						return req, err
					},
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any, data2 []map[string]any) {
						require.Equal(t, http.StatusNoContent, resp.StatusCode)
					},
				},
				{
					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
						req, err := http.NewRequest(http.MethodPost, ApiUrl+"/balance_topup", bytes.NewBufferString(`{"amount":50000}`))
						req.Header.Set("Authorization", "Bearer "+tc.Steps[1].Result["token"].(string))
						return req, err
					},
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any, data2 []map[string]any) {
						require.Equal(t, http.StatusNoContent, resp.StatusCode)
					},
				},
				{
					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
						req, err := http.NewRequest(http.MethodGet, ApiUrl+"/balance_read", nil)
						req.Header.Set("Authorization", "Bearer "+tc.Steps[0].Result["token"].(string))
						return req, err
					},
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any, data2 []map[string]any) {
						require.Equal(t, http.StatusOK, resp.StatusCode)
						require.Equal(t, float64(200000), data["balance"].(float64))
					},
				},
				{
					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
						req, err := http.NewRequest(http.MethodGet, ApiUrl+"/balance_read", nil)
						req.Header.Set("Authorization", "Bearer "+tc.Steps[1].Result["token"].(string))
						return req, err
					},
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any, data2 []map[string]any) {
						require.Equal(t, http.StatusOK, resp.StatusCode)
						require.Equal(t, float64(50000), data["balance"].(float64))
					},
				},
				{
					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
						req, err := http.NewRequest(http.MethodPost, ApiUrl+"/transfer", bytes.NewBufferString(`{"amount":100000,"to_username":"`+PrefixUsername+`username7"}`))
						req.Header.Set("Authorization", "Bearer "+tc.Steps[0].Result["token"].(string))
						return req, err
					},
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any, data2 []map[string]any) {
						require.Equal(t, http.StatusNoContent, resp.StatusCode)
					},
				},
				{
					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
						req, err := http.NewRequest(http.MethodPost, ApiUrl+"/transfer", bytes.NewBufferString(`{"amount":25000,"to_username":"`+PrefixUsername+`username6"}`))
						req.Header.Set("Authorization", "Bearer "+tc.Steps[1].Result["token"].(string))
						return req, err
					},
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any, data2 []map[string]any) {
						require.Equal(t, http.StatusNoContent, resp.StatusCode)
					},
				},
				{
					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
						req, err := http.NewRequest(http.MethodGet, ApiUrl+"/balance_read", nil)
						req.Header.Set("Authorization", "Bearer "+tc.Steps[0].Result["token"].(string))
						return req, err
					},
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any, data2 []map[string]any) {
						require.Equal(t, http.StatusOK, resp.StatusCode)
						require.Equal(t, float64(125000), data["balance"].(float64))
					},
				},
				{
					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
						req, err := http.NewRequest(http.MethodGet, ApiUrl+"/balance_read", nil)
						req.Header.Set("Authorization", "Bearer "+tc.Steps[1].Result["token"].(string))
						return req, err
					},
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any, data2 []map[string]any) {
						require.Equal(t, http.StatusOK, resp.StatusCode)
						require.Equal(t, float64(125000), data["balance"].(float64))
					},
				},
				{
					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
						req, err := http.NewRequest(http.MethodGet, ApiUrl+"/top_transaction_per_user", nil)
						req.Header.Set("Authorization", "Bearer "+tc.Steps[0].Result["token"].(string))
						return req, err
					},
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any, data2 []map[string]any) {
						require.Equal(t, http.StatusOK, resp.StatusCode)
						require.Equal(t, int(4), len(data2))

						require.Equal(t, float64(100000), data2[0]["amount"].(float64))
						require.Equal(t, PrefixUsername+"username6", data2[0]["username"].(string))

						require.Equal(t, float64(100000), data2[1]["amount"].(float64))
						require.Equal(t, PrefixUsername+"username6", data2[1]["username"].(string))

						require.Equal(t, float64(25000), data2[2]["amount"].(float64))
						require.Equal(t, PrefixUsername+"username7", data2[2]["username"].(string))

						require.Equal(t, float64(-100000), data2[3]["amount"].(float64))
						require.Equal(t, PrefixUsername+"username7", data2[3]["username"].(string))
					},
				},
				{
					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
						req, err := http.NewRequest(http.MethodGet, ApiUrl+"/top_transaction_per_user", nil)
						req.Header.Set("Authorization", "Bearer "+tc.Steps[1].Result["token"].(string))
						return req, err
					},
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any, data2 []map[string]any) {
						require.Equal(t, http.StatusOK, resp.StatusCode)
						require.Equal(t, int(3), len(data2))

						require.Equal(t, float64(100000), data2[0]["amount"].(float64))
						require.Equal(t, PrefixUsername+"username6", data2[0]["username"].(string))

						require.Equal(t, float64(50000), data2[1]["amount"].(float64))
						require.Equal(t, PrefixUsername+"username7", data2[1]["username"].(string))

						require.Equal(t, float64(-25000), data2[2]["amount"].(float64))
						require.Equal(t, PrefixUsername+"username6", data2[2]["username"].(string))
					},
				},
				{
					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
						req, err := http.NewRequest(http.MethodGet, ApiUrl+"/top_users", nil)
						req.Header.Set("Authorization", "Bearer "+tc.Steps[0].Result["token"].(string))
						return req, err
					},
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any, data2 []map[string]any) {
						require.Equal(t, http.StatusOK, resp.StatusCode)
						require.Equal(t, int(1), len(data2))
						require.Equal(t, float64(100000), data2[0]["transacted_value"].(float64))
						require.Equal(t, PrefixUsername+"username7", data2[0]["username"].(string))
					},
				},
				{
					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
						req, err := http.NewRequest(http.MethodGet, ApiUrl+"/top_users", nil)
						req.Header.Set("Authorization", "Bearer "+tc.Steps[1].Result["token"].(string))
						return req, err
					},
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any, data2 []map[string]any) {
						require.Equal(t, http.StatusOK, resp.StatusCode)
						require.Equal(t, int(1), len(data2))
						require.Equal(t, float64(25000), data2[0]["transacted_value"].(float64))
						require.Equal(t, PrefixUsername+"username6", data2[0]["username"].(string))
					},
				},
			},
		},
	}
}

type TestCase struct {
	Name  string
	Steps []TestCaseStep
}

type RequestFunc func(*testing.T, context.Context, *TestCase) (*http.Request, error)
type ExpectFunc func(*testing.T, context.Context, *TestCase, *http.Response, map[string]any, []map[string]any)

type TestCaseStep struct {
	Request     RequestFunc
	Expect      ExpectFunc
	Result      map[string]any
	ResultArray []map[string]any
}

func ReadJsonResult(t *testing.T, resp *http.Response, step *TestCaseStep) {
	if resp.StatusCode == http.StatusNoContent {
		return
	}

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	err = json.Unmarshal(body, &step.ResultArray)
	if err != nil {
		err = json.Unmarshal(body, &step.Result)
	}
	require.NoError(t, err)
}
