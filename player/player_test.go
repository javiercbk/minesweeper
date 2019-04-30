package player

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

const testErrHash = "errorHashing"
const testErrPlayerExists = "playerExists"
const testErrInsert = "errorInsert"
const testPlayerID = 64

var testErrExists = response.HTTPError{
	Code:    http.StatusConflict,
	Message: fmt.Sprintf("player playerExists already exists", testErrPlayerExists),
}

type mockAPI struct{}

func (m mockAPI) CreatePlayer(ctx context.Context, pPlayer *ProspectPlayer) error {
	hashPassword, err := HashPassword(pPlayer.Password)
	if pPlayer.Name == testErrHash {
		return errors.New("error hashing password")
	}
	if pPlayer.Name == testErrPlayerExists {
		return testErrExists
	}
	if pPlayer.Name == testErrInsert {
		return errors.New("error inserting player")
	}
	pPlayer.ID = testPlayerID
	return nil
}

func compare(expected, given interface{}) error {
	expectedTR, ok := expected.(ProspectPlayer)
	if !ok {
		return fmt.Errorf("expected is not a ProspectPlayer")
	}
	givenTR, ok := expected.(ProspectPlayer)
	if !ok {
		return fmt.Errorf("given is not a ProspectPlayer")
	}
	return nil
}

func TestCreatePlayer(t *testing.T) {
	// Setup
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
			Body:        testHelpers.MarshalIgnore(ProspectPlayer{Name: "user"}),
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
			Body:        testHelpers.MarshalIgnore(ProspectPlayer{Password: "pass"}),
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
			Body:        testHelpers.MarshalIgnore(ProspectPlayer{Name: testErrInsert, Password: "pass"}),
			ExpectedResponse: response.ServiceResponse{
				Status: response.Status{
					Error:   true,
					Code:    http.StatusInternalServerError,
					Message: "error inserting player",
				},
			},
		},
		{
			Path:        "/api",
			Method:      http.MethodPost,
			ContentType: echo.MIMEApplicationJSON,
			Body:        testHelpers.MarshalIgnore(ProspectPlayer{Name: testErrPlayerExists, Password: "pass"}),
			ExpectedResponse: response.ServiceResponse{
				Status: response.Status{
					Error:   true,
					Code:    testErrExists.Code,
					Message: testErrExists.Message,
				},
			},
		},
		{
			Path:        "/api",
			Method:      http.MethodPost,
			ContentType: echo.MIMEApplicationJSON,
			Body:        testHelpers.MarshalIgnore(ProspectPlayer{Name: testErrHash, Password: "pass"}),
			ExpectedResponse: response.ServiceResponse{
				Status: response.Status{
					Error:   true,
					Code:    http.StatusInternalServerError,
					Message: "error hashing password",
				},
			},
		},
		{
			Path:        "/api",
			Method:      http.MethodPost,
			ContentType: echo.MIMEApplicationJSON,
			Body:        testHelpers.MarshalIgnore(ProspectPlayer{Name: "user", Password: "pass"}),
			ExpectedResponse: response.ServiceResponse{
				Status: response.Status{
					Error: false,
					Code:  http.StatusOK,
				},
				Data: ProspectPlayer{ID: testPlayerID, Name: "user", Password: "pass"},
			},
		},
	}
	e := testHelpers.MockEcho()
	apiRouter := e.Group("/api")
	apiFactory = func(logger *log.Logger, db *sql.DB) API {
		return mockAPI{}
	}
	handler := NewHandler(testHelpers.NullLogger(), nil)
	handler.Routes(apiRouter)
	for i, test := range tests {
		requestText := ""
		if test.Body != "" {
			requestText = test.Body
		}
		req := httptest.NewRequest(test.Method, test.Path, strings.NewReader(requestText))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		handler.AuthenticateFactory(jwtSecret)(c)
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
