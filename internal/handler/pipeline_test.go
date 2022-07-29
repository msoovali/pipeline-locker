package handler

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/msoovali/pipeline-locker/internal/domain"
	"github.com/valyala/fasthttp"
)

type pipelineServiceMock struct {
	domain.PipelineService
	fakeIsDeployAllowed    func(pipeline domain.PipelineIdentifier) (bool, error)
	fakeLock               func(pipeline domain.PipelineLockRequest) error
	fakeUnlock             func(pipeline domain.PipelineIdentifier) error
	fakeGetLockedPipelines func() ([]domain.Pipeline, error)
}

func (m *pipelineServiceMock) IsDeployAllowed(pipeline domain.PipelineIdentifier) (bool, error) {
	if m.fakeIsDeployAllowed != nil {
		return m.fakeIsDeployAllowed(pipeline)
	}

	return true, nil
}

func (m *pipelineServiceMock) Lock(pipeline domain.PipelineLockRequest) error {
	if m.fakeLock != nil {
		return m.fakeLock(pipeline)
	}

	return nil
}

func (m *pipelineServiceMock) Unlock(pipeline domain.PipelineIdentifier) error {
	if m.fakeUnlock != nil {
		return m.fakeUnlock(pipeline)
	}

	return nil
}

func (m *pipelineServiceMock) GetLockedPipelines() ([]domain.Pipeline, error) {
	if m.fakeGetLockedPipelines != nil {
		return m.fakeGetLockedPipelines()
	}

	return nil, nil
}

func getLockRequestBodyMock() string {
	return "{\"project\":\"proj\",\"environment\":\"env\",\"locked_by\":\"user\"}"
}

func getPipelineRequestBodyMock() string {
	return "{\"project\":\"proj\",\"environment\":\"env\"}"
}

func TestPipelineHandler_Lock(t *testing.T) {
	type testCases struct {
		description          string
		requestBody          string
		expectedStatus       int
		expectedResponseBody string
		fakeLockReturnValue  error
		contentTypeHeader    string
	}
	for _, scenario := range []testCases{
		{
			description:          "brokenRequestBody_respondBadRequest",
			requestBody:          "{123",
			expectedStatus:       fiber.StatusBadRequest,
			expectedResponseBody: "Unprocessable Entity",
		},
		{
			description:          "serviceReturnsError_respondConflict",
			requestBody:          getLockRequestBodyMock(),
			expectedStatus:       fiber.StatusConflict,
			expectedResponseBody: domain.ErrPipelineAlreadyLocked.Error(),
			fakeLockReturnValue:  domain.ErrPipelineAlreadyLocked,
			contentTypeHeader:    "application/json",
		},
		{
			description:          "serviceReturnsNil_respondNoContent",
			requestBody:          getLockRequestBodyMock(),
			expectedStatus:       fiber.StatusCreated,
			expectedResponseBody: "Created",
			contentTypeHeader:    "application/json",
		},
	} {
		t.Run(scenario.description, func(t *testing.T) {
			handler := NewPipelineHandlers(&pipelineServiceMock{
				fakeLock: func(pipeline domain.PipelineLockRequest) error {
					return scenario.fakeLockReturnValue
				},
			})
			app := fiber.New()
			c := app.AcquireCtx(&fasthttp.RequestCtx{})
			defer app.ReleaseCtx(c)
			c.Request().Header.Add("content-type", scenario.contentTypeHeader)
			c.Request().AppendBodyString(scenario.requestBody)

			handler.Lock(c)

			if c.Response().StatusCode() != scenario.expectedStatus {
				t.Errorf("Expected status %d, got %d", scenario.expectedStatus, c.Response().StatusCode())
			}
			if string(c.Response().Body()) != scenario.expectedResponseBody {
				t.Errorf("Expected body %s, got %s", scenario.expectedResponseBody, string(c.Response().Body()))
			}
		})
	}
}

func TestPipelineHandler_Unlock(t *testing.T) {
	type testCases struct {
		description           string
		requestBody           string
		expectedStatus        int
		expectedResponseBody  string
		fakeUnlockReturnValue error
		contentTypeHeader     string
	}
	for _, scenario := range []testCases{
		{
			description:          "brokenRequestBody_respondBadRequest",
			requestBody:          "123",
			expectedStatus:       fiber.StatusBadRequest,
			expectedResponseBody: "Unprocessable Entity",
		},
		{
			description:           "serviceReturnsError_respondConflict",
			requestBody:           getPipelineRequestBodyMock(),
			expectedStatus:        fiber.StatusConflict,
			expectedResponseBody:  domain.ErrProjectEmpty.Error(),
			fakeUnlockReturnValue: domain.ErrProjectEmpty,
			contentTypeHeader:     "application/json",
		},
		{
			description:          "serviceReturnsNil_respondNoContent",
			requestBody:          getPipelineRequestBodyMock(),
			expectedStatus:       fiber.StatusNoContent,
			expectedResponseBody: "No Content",
			contentTypeHeader:    "application/json",
		},
	} {
		t.Run(scenario.description, func(t *testing.T) {
			handler := NewPipelineHandlers(&pipelineServiceMock{
				fakeUnlock: func(pipeline domain.PipelineIdentifier) error {
					return scenario.fakeUnlockReturnValue
				},
			})
			app := fiber.New()
			c := app.AcquireCtx(&fasthttp.RequestCtx{})
			defer app.ReleaseCtx(c)
			c.Request().Header.Add("content-type", scenario.contentTypeHeader)
			c.Request().AppendBodyString(scenario.requestBody)

			handler.Unlock(c)

			if c.Response().StatusCode() != scenario.expectedStatus {
				t.Errorf("Expected status %d, got %d", scenario.expectedStatus, c.Response().StatusCode())
			}
			if string(c.Response().Body()) != scenario.expectedResponseBody {
				t.Errorf("Expected body %s, got %s", scenario.expectedResponseBody, string(c.Response().Body()))
			}
		})
	}
}
