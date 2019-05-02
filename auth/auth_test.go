package auth

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/javiercbk/minesweeper/http/response"
	testHelpers "github.com/javiercbk/minesweeper/testing"
	"github.com/labstack/echo"
)

const testErrSearching = "errorSearching"
const testErrBadCredentials = "badCredentials"
const testErrToken = "errorToken"

const testOKToken = "ok"

type mockAPI struct{}

func (m mockAPI) CreateToken(ctx context.Context, jwtSecret string, auth Credentials) (TokenResponse, error) {
	tResponse := TokenResponse{}
	if auth.Name == testErrSearching {
		return tResponse, errors.New("error searching for player")
	}
	if auth.Name == testErrBadCredentials {
		return tResponse, ErrBadCredentials
	}
	if auth.Name == testErrToken {
		return tResponse, errors.New("error creating token")
	}
	tResponse.Token = testOKToken
	return tResponse, nil
}

func compare(expected, given interface{}) error {
	expectedTR, ok := expected.(TokenResponse)
	if !ok {
		return fmt.Errorf("expected is not a TokenResponse")
	}
	givenTR, ok := expected.(TokenResponse)
	if !ok {
		return fmt.Errorf("given is not a TokenResponse")
	}
	if expectedTR.Token != givenTR.Token {
		return fmt.Errorf("expected token to be %s but was %s", expectedTR.Token, givenTR.Token)
	}
	return nil
}

func TestAuthenticate(t *testing.T) {
	tests := []testHelpers.EchoUnitTest{
		{
			Path:        "/api",
			Method:      http.MethodPost,
			ContentType: echo.MIMETextPlain,
			ExpectedResponse: response.ServiceResponse{
				Status: response.Status{
					Error: true,
					Code:  http.StatusBadRequest,
				},
			},
		},
		{
			Path:        "/api",
			Method:      http.MethodPost,
			ContentType: echo.MIMEApplicationJSON,
			Body:        "{}",
			ExpectedResponse: response.ServiceResponse{
				Status: response.Status{
					Error: true,
					Code:  http.StatusBadRequest,
				},
			},
		},
		{
			Path:        "/api",
			Method:      http.MethodPost,
			ContentType: echo.MIMEApplicationJSON,
			Body:        testHelpers.MarshalIgnore(Credentials{Name: "user"}),
			ExpectedResponse: response.ServiceResponse{
				Status: response.Status{
					Error: true,
					Code:  http.StatusBadRequest,
				},
			},
		},
		{
			Path:        "/api",
			Method:      http.MethodPost,
			ContentType: echo.MIMEApplicationJSON,
			Body:        testHelpers.MarshalIgnore(Credentials{Password: "pass"}),
			ExpectedResponse: response.ServiceResponse{
				Status: response.Status{
					Error: true,
					Code:  http.StatusBadRequest,
				},
			},
		},
		{
			Path:        "/api",
			Method:      http.MethodPost,
			ContentType: echo.MIMEApplicationJSON,
			Body:        testHelpers.MarshalIgnore(Credentials{Name: testErrSearching, Password: "pass"}),
			ExpectedResponse: response.ServiceResponse{
				Status: response.Status{
					Error:   true,
					Code:    http.StatusInternalServerError,
					Message: "error searching for player",
				},
			},
		},
		{
			Path:        "/api",
			Method:      http.MethodPost,
			ContentType: echo.MIMEApplicationJSON,
			Body:        testHelpers.MarshalIgnore(Credentials{Name: testErrBadCredentials, Password: "pass"}),
			ExpectedResponse: response.ServiceResponse{
				Status: response.Status{
					Error:   true,
					Code:    http.StatusUnauthorized,
					Message: ErrBadCredentials.Error(),
				},
			},
		},
		{
			Path:        "/api",
			Method:      http.MethodPost,
			ContentType: echo.MIMEApplicationJSON,
			Body:        testHelpers.MarshalIgnore(Credentials{Name: testErrToken, Password: "pass"}),
			ExpectedResponse: response.ServiceResponse{
				Status: response.Status{
					Error:   true,
					Code:    http.StatusInternalServerError,
					Message: "error creating token",
				},
			},
		},
		{
			Path:        "/api",
			Method:      http.MethodPost,
			ContentType: echo.MIMEApplicationJSON,
			Body:        testHelpers.MarshalIgnore(Credentials{Name: "user", Password: "pass"}),
			ExpectedResponse: response.ServiceResponse{
				Status: response.Status{
					Error: false,
					Code:  http.StatusOK,
				},
				Data: TokenResponse{
					Token: testOKToken,
				},
			},
		},
	}
	e := testHelpers.MockEcho()
	apiRouter := e.Group("/api")
	apiFactory = func(logger *log.Logger, db *sql.DB) API {
		return mockAPI{}
	}
	handler := NewHandler(testHelpers.NullLogger(), nil)
	handler.Routes(apiRouter, jwtSecret)
	for i, test := range tests {
		requestText := ""
		if test.Body != "" {
			requestText = test.Body
		}
		req := httptest.NewRequest(test.Method, test.Path, strings.NewReader(requestText))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		_ = handler.AuthenticateFactory(jwtSecret)(c)
		given := response.ServiceResponse{}
		err := json.Unmarshal(rec.Body.Bytes(), &given)
		if err != nil {
			t.Fatalf("Test %d failed: error unmarshalling http response %s", i, err)
		}
		err = testHelpers.AssertEchoResponse(test.ExpectedResponse, given, compare)
		if err != nil {
			t.Fatalf("Test %d failed: %s", i, err)
		}
	}
}
