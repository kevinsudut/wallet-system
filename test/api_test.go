package test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

const ApiUrl = "http://localhost:8000"

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
				step.Expect(t, ctx, &tc, response, step.Result)
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
						return http.NewRequest(http.MethodPost, ApiUrl+"/create_user", bytes.NewBufferString(`{"username":"username"}`))
					},
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
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
						return http.NewRequest(http.MethodPost, ApiUrl+"/create_user", bytes.NewBufferString(`{"username":"username1"}`))
					},
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
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
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
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
						return http.NewRequest(http.MethodPost, ApiUrl+"/create_user", bytes.NewBufferString(`{"username":"username3"}`))
					},
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
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
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
						require.Equal(t, http.StatusNoContent, resp.StatusCode)
					},
				},
				{
					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
						req, err := http.NewRequest(http.MethodPost, ApiUrl+"/balance_topup", bytes.NewBufferString(`{"amount":50000}`))
						req.Header.Set("Authorization", "Bearer "+tc.Steps[0].Result["token"].(string))
						return req, err
					},
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
						require.Equal(t, http.StatusNoContent, resp.StatusCode)
					},
				},
				{
					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
						req, err := http.NewRequest(http.MethodGet, ApiUrl+"/balance_read", nil)
						req.Header.Set("Authorization", "Bearer "+tc.Steps[0].Result["token"].(string))
						return req, err
					},
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
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
						return http.NewRequest(http.MethodPost, ApiUrl+"/create_user", bytes.NewBufferString(`{"username":"username4"}`))
					},
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
						require.Equal(t, http.StatusCreated, resp.StatusCode)
						require.NotEmpty(t, data["token"].(string))
					},
				},
				{
					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
						return http.NewRequest(http.MethodPost, ApiUrl+"/create_user", bytes.NewBufferString(`{"username":"username5"}`))
					},
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
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
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
						require.Equal(t, http.StatusNoContent, resp.StatusCode)
					},
				},
				{
					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
						req, err := http.NewRequest(http.MethodPost, ApiUrl+"/balance_topup", bytes.NewBufferString(`{"amount":50000}`))
						req.Header.Set("Authorization", "Bearer "+tc.Steps[1].Result["token"].(string))
						return req, err
					},
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
						require.Equal(t, http.StatusNoContent, resp.StatusCode)
					},
				},
				{
					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
						req, err := http.NewRequest(http.MethodGet, ApiUrl+"/balance_read", nil)
						req.Header.Set("Authorization", "Bearer "+tc.Steps[0].Result["token"].(string))
						return req, err
					},
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
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
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
						require.Equal(t, http.StatusOK, resp.StatusCode)
						require.Equal(t, float64(50000), data["balance"].(float64))
					},
				},
				{
					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
						req, err := http.NewRequest(http.MethodPost, ApiUrl+"/transfer", bytes.NewBufferString(`{"amount":25000,"to_username":"username4"}`))
						req.Header.Set("Authorization", "Bearer "+tc.Steps[1].Result["token"].(string))
						return req, err
					},
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
						require.Equal(t, http.StatusNoContent, resp.StatusCode)
					},
				},
				{
					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
						req, err := http.NewRequest(http.MethodGet, ApiUrl+"/balance_read", nil)
						req.Header.Set("Authorization", "Bearer "+tc.Steps[0].Result["token"].(string))
						return req, err
					},
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
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
					Expect: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
						require.Equal(t, http.StatusOK, resp.StatusCode)
						require.Equal(t, float64(25000), data["balance"].(float64))
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
type ExpectFunc func(*testing.T, context.Context, *TestCase, *http.Response, map[string]any)

type TestCaseStep struct {
	Request RequestFunc
	Expect  ExpectFunc
	Result  map[string]any
}

func ResponseContains(t *testing.T, resp *http.Response, text string) {
	body, err := io.ReadAll(resp.Body)
	bodyStr := string(body)
	require.NoError(t, err)
	require.Contains(t, bodyStr, text)
}

func ReadJsonResult(t *testing.T, resp *http.Response, step *TestCaseStep) {
	if resp.StatusCode == http.StatusNoContent {
		return
	}

	var result map[string]any
	err := json.NewDecoder(resp.Body).Decode(&result)
	step.Result = result
	require.NoError(t, err)
}

func RequireIsUUID(t *testing.T, value string) {
	_, err := uuid.Parse(value)
	require.NoError(t, err)
}
