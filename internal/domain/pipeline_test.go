package domain

import (
	"errors"
	"testing"
)

const project = "area51"
const environment = "production"
const lockedBy = "user"

func getValidIdentifier() PipelineIdentifier {
	return PipelineIdentifier{
		Project:     project,
		Environment: environment,
	}
}

func TestPipelineIdentifier_Validate(t *testing.T) {
	type testCases struct {
		description        string
		pipelineIdentifier PipelineIdentifier
		expectedError      error
	}

	for _, scenario := range []testCases{
		{
			description:        "projectIsEmpty_returnProjectEmptyError",
			pipelineIdentifier: PipelineIdentifier{},
			expectedError:      ErrProjectEmpty,
		},
		{
			description: "environmentIsEmpty_returnEnvironmentEmptyError",
			pipelineIdentifier: PipelineIdentifier{
				Project: project,
			},
			expectedError: ErrEnvironmentEmpty,
		},
		{
			description:        "success",
			pipelineIdentifier: getValidIdentifier(),
			expectedError:      nil,
		},
	} {
		t.Run(scenario.description, func(t *testing.T) {
			err := scenario.pipelineIdentifier.Validate()

			if !errors.Is(err, scenario.expectedError) {
				t.Errorf("Expected %v, received %v", scenario.expectedError, err)
			}
		})
	}
}

func TestPipelineLockRequest_Validate(t *testing.T) {
	type testCases struct {
		description   string
		request       PipelineLockRequest
		expectedError error
	}

	for _, scenario := range []testCases{
		{
			description:   "identifierValidateIsCalled_returnProjectEmptyError",
			expectedError: ErrProjectEmpty,
		},
		{
			description: "lockedByIsEmpty_returnLockedByEmptyError",
			request: PipelineLockRequest{
				PipelineIdentifier: getValidIdentifier(),
			},
			expectedError: ErrLockedByEmpty,
		},
		{
			description: "success",
			request: PipelineLockRequest{
				PipelineIdentifier: getValidIdentifier(),
				PipelineLockedBy: PipelineLockedBy{
					LockedBy: lockedBy,
				},
			},
			expectedError: nil,
		},
	} {
		t.Run(scenario.description, func(t *testing.T) {
			err := scenario.request.Validate()

			if !errors.Is(err, scenario.expectedError) {
				t.Errorf("Expected %v, received %v", scenario.expectedError, err)
			}
		})
	}
}
