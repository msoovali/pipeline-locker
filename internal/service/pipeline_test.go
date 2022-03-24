package service

import (
	"errors"
	"testing"

	"github.com/msoovali/pipeline-locker/internal/domain"
	"github.com/msoovali/pipeline-locker/internal/repository"
)

const (
	project     string = "project"
	environment string = "environment"
	user        string = "user"
)

type pipelineRepositoryMock struct {
	repository.PipelineRepository
	fakeAdd                 func(pipeline domain.Pipeline)
	fakeFind                func(pipeline domain.PipelineIdentifier) *domain.Pipeline
	fakeFindLockedPipelines func() []domain.Pipeline
}

func (r *pipelineRepositoryMock) Find(pipeline domain.PipelineIdentifier) *domain.Pipeline {
	if r.fakeFind != nil {
		return r.fakeFind(pipeline)
	}

	return nil
}

func (r *pipelineRepositoryMock) Add(pipeline domain.Pipeline) {
	if r.fakeAdd != nil {
		r.fakeAdd(pipeline)
	}
}

func (r *pipelineRepositoryMock) FindLockedPipelines() []domain.Pipeline {
	if r.fakeFindLockedPipelines != nil {
		return r.fakeFindLockedPipelines()
	}

	return make([]domain.Pipeline, 0)
}

func getPipelineMock() *domain.Pipeline {
	return &domain.Pipeline{
		PipelineIdentifier: domain.PipelineIdentifier{
			Project:     project,
			Environment: environment,
		},
	}
}

func getPipelineMockWithLockedByValue() *domain.Pipeline {
	pipeline := getPipelineMock()
	lockedBy := user
	pipeline.LockedBy = lockedBy

	return pipeline
}

func TestPipelineService_Lock(t *testing.T) {
	type testCases struct {
		description             string
		input                   domain.Pipeline
		serviceAllowOverLocking bool
		expectedError           error
		expectedAddCalls        int
		fakeFindReturnValue     *domain.Pipeline
	}

	for _, scenario := range []testCases{
		{
			description:   "projectIsEmpty_returnError",
			input:         domain.Pipeline{},
			expectedError: ProjectEmptyError,
		},
		{
			description: "environmentIsEmpty_returnError",
			input: domain.Pipeline{
				PipelineIdentifier: domain.PipelineIdentifier{
					Project: project,
				},
			},
			expectedError: EnvironmentEmptyError,
		},
		{
			description:   "lockedByIsNil_returnError",
			input:         *getPipelineMock(),
			expectedError: LockedByEmptyError,
		},
		{
			description: "lockedByIsEmptyString_returnError",
			input: domain.Pipeline{
				PipelineIdentifier: getPipelineMock().PipelineIdentifier,
				LockedBy:           "",
			},
			expectedError: LockedByEmptyError,
		},
		{
			description:         "lockAlreadyExistsAndOverLockingNotAllowed_returnError",
			input:               *getPipelineMockWithLockedByValue(),
			expectedError:       PipelineAlreadyLockedError,
			fakeFindReturnValue: getPipelineMockWithLockedByValue(),
		},
		{
			description:      "pipelineNotExistsAndOverLockingNotAllowed_callsAdd",
			input:            *getPipelineMockWithLockedByValue(),
			expectedAddCalls: 1,
		},
		{
			description:         "pipelineAlreadyExistsAndOverLockingNotAllowed_callsAdd",
			input:               *getPipelineMockWithLockedByValue(),
			expectedAddCalls:    1,
			fakeFindReturnValue: getPipelineMock(),
		},
		{
			description:             "lockAlreadyExistsButOverLockingIsAllowed_callsAdd",
			input:                   *getPipelineMockWithLockedByValue(),
			serviceAllowOverLocking: true,
			expectedAddCalls:        1,
			fakeFindReturnValue:     getPipelineMockWithLockedByValue(),
		},
	} {
		t.Run(scenario.description, func(t *testing.T) {
			var addCallsCount int
			repository := &pipelineRepositoryMock{
				fakeFind: func(pipeline domain.PipelineIdentifier) *domain.Pipeline {
					return scenario.fakeFindReturnValue
				},
				fakeAdd: func(pipeline domain.Pipeline) {
					addCallsCount++
				},
			}
			service := NewPipelineService(repository, scenario.serviceAllowOverLocking)

			err := service.Lock(scenario.input)

			if !errors.Is(err, scenario.expectedError) {
				t.Errorf("Expected error %s, but received %s", scenario.expectedError, err)
			}
			if addCallsCount != scenario.expectedAddCalls {
				t.Errorf("Expected repository Add method calls %d, but Add was called %d times", scenario.expectedAddCalls, addCallsCount)
			}
		})
	}
}

func TestPipelineService_Unlock(t *testing.T) {
	type testCases struct {
		description      string
		input            domain.PipelineIdentifier
		expectedError    error
		expectedAddCalls int
	}

	for _, scenario := range []testCases{
		{
			description:   "projectIsEmpty_returnError",
			input:         domain.PipelineIdentifier{},
			expectedError: ProjectEmptyError,
		},
		{
			description: "environmentIsEmpty_returnError",
			input: domain.PipelineIdentifier{
				Project: project,
			},
			expectedError: EnvironmentEmptyError,
		},
		{
			description:      "inputIsOK_callsAdd",
			input:            getPipelineMock().PipelineIdentifier,
			expectedAddCalls: 1,
		},
	} {
		t.Run(scenario.description, func(t *testing.T) {
			var addCallsCount int
			var lockedByValue string
			repository := &pipelineRepositoryMock{
				fakeAdd: func(pipeline domain.Pipeline) {
					addCallsCount++
					lockedByValue = pipeline.LockedBy
				},
			}
			service := NewPipelineService(repository, false)

			err := service.Unlock(scenario.input)

			if !errors.Is(err, scenario.expectedError) {
				t.Errorf("Expected error %s, but received %s", scenario.expectedError, err)
			}
			if addCallsCount != scenario.expectedAddCalls {
				t.Errorf("Expected repository Add method calls %d, but Add was called %d times", scenario.expectedAddCalls, addCallsCount)
			}
			if lockedByValue != "" {
				t.Errorf("Expected unlock to call Add without LockedBy value, but received %s", lockedByValue)
			}
		})
	}
}

func TestPipelineService_IsDeployAllowed(t *testing.T) {
	type testCases struct {
		description         string
		input               domain.PipelineIdentifier
		expectedError       error
		expectedValue       bool
		fakeFindReturnValue *domain.Pipeline
	}

	for _, scenario := range []testCases{
		{
			description:   "projectIsEmpty_returnError",
			input:         domain.PipelineIdentifier{},
			expectedError: ProjectEmptyError,
		},
		{
			description: "environmentIsEmpty_returnError",
			input: domain.PipelineIdentifier{
				Project: project,
			},
			expectedError: EnvironmentEmptyError,
		},
		{
			description:   "pipelineIsNotFoundFromStore_returnTrue",
			input:         getPipelineMock().PipelineIdentifier,
			expectedValue: true,
		},
		{
			description:         "pipelineIsLocked_returnFalse",
			input:               getPipelineMock().PipelineIdentifier,
			fakeFindReturnValue: getPipelineMockWithLockedByValue(),
		},
		{
			description:         "pipelineIsNotLocked_returnTrue",
			input:               getPipelineMock().PipelineIdentifier,
			fakeFindReturnValue: getPipelineMock(),
			expectedValue:       true,
		},
	} {
		t.Run(scenario.description, func(t *testing.T) {
			repository := &pipelineRepositoryMock{
				fakeFind: func(pipeline domain.PipelineIdentifier) *domain.Pipeline {
					return scenario.fakeFindReturnValue
				},
			}
			service := NewPipelineService(repository, false)

			isAllowed, err := service.IsDeployAllowed(scenario.input)

			if !errors.Is(err, scenario.expectedError) {
				t.Errorf("Expected error %s, but received %s", scenario.expectedError, err)
			}
			if isAllowed != scenario.expectedValue {
				t.Errorf("Expected return value %t, got %t", scenario.expectedValue, isAllowed)
			}
		})
	}
}

func TestPipelineService_GetLockedPipelines(t *testing.T) {
	t.Run("repositoryFindLockedPipelinesIsCalled_proxiesValue", func(t *testing.T) {
		var findLockedPipelinesCalls int
		repository := &pipelineRepositoryMock{
			fakeFindLockedPipelines: func() []domain.Pipeline {
				findLockedPipelinesCalls++
				return make([]domain.Pipeline, 0)
			},
		}
		service := NewPipelineService(repository, false)

		lockedPipelines := service.GetLockedPipelines()

		if findLockedPipelinesCalls != 1 {
			t.Errorf("Expected repository findLockedPipelines() to be called once, got %d", findLockedPipelinesCalls)
		}
		if lockedPipelines == nil || len(lockedPipelines) != 0 {
			t.Errorf("Expected empty slice to be returned, but got %v", lockedPipelines)
		}
	})
}
