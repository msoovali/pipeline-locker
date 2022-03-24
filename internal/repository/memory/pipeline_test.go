package memory

import (
	"testing"

	"github.com/msoovali/pipeline-locker/internal/domain"
)

func TestPipelineRepository(t *testing.T) {
	const projectOne string = "Project"
	const projectTwo string = "project"
	const environmentOne string = "Environment"
	const environmentTwo string = "environment"
	const userOne = "User"
	const userTwo = "user"
	var repository *pipelineRepository
	const pipelineKeyCaseSensitive = true

	t.Run("NewPipelineRepository_repositoryIsInitialized_storeSizeIsZero", func(t *testing.T) {
		repository = NewPipelineRepository(pipelineKeyCaseSensitive)
		if repository.store == nil {
			t.Errorf("Expected store to be initialized, but got nil")
		}
		if len(repository.store) != 0 {
			t.Errorf("Expected store to be empty, but got store size %d instead", len(repository.store))
		}
	})

	type addTestCases struct {
		description       string
		pipeline          domain.Pipeline
		expectedStoreSize int
	}

	for _, scenario := range []addTestCases{
		{
			description: "Add_projectOneAdded_savesToStore",
			pipeline: domain.Pipeline{
				PipelineIdentifier: domain.PipelineIdentifier{
					Project:     projectOne,
					Environment: environmentOne,
				},
				PipelineLockedBy: domain.PipelineLockedBy{
					LockedBy: userOne,
				},
			},
			expectedStoreSize: 1,
		},
		{
			description: "Add_projectOneAddedAnotherTime_overwritesStoreKey",
			pipeline: domain.Pipeline{
				PipelineIdentifier: domain.PipelineIdentifier{
					Project:     projectOne,
					Environment: environmentOne,
				},
				PipelineLockedBy: domain.PipelineLockedBy{
					LockedBy: userTwo,
				},
			},
			expectedStoreSize: 1,
		},
		{
			description: "Add_projectTwoAdded_savesToStore",
			pipeline: domain.Pipeline{
				PipelineIdentifier: domain.PipelineIdentifier{
					Project:     projectTwo,
					Environment: environmentTwo,
				},
				PipelineLockedBy: domain.PipelineLockedBy{
					LockedBy: userTwo,
				},
			},
			expectedStoreSize: 2,
		},
		{
			description: "Add_projectTwoAddeAnotherTimedWithoutLocker_overwritesStoreKey",
			pipeline: domain.Pipeline{
				PipelineIdentifier: domain.PipelineIdentifier{
					Project:     projectTwo,
					Environment: environmentTwo,
				},
			},
			expectedStoreSize: 2,
		},
	} {
		t.Run(scenario.description, func(t *testing.T) {
			repository.Add(scenario.pipeline)

			if len(repository.store) != scenario.expectedStoreSize {
				t.Errorf("Expected store size %d, but got %d", scenario.expectedStoreSize, len(repository.store))
			}
			key := scenario.pipeline.Project + separator + scenario.pipeline.Environment
			value, exists := repository.store[key]
			if !exists {
				t.Errorf("Expected key %s to be added to store, but was not found from store", key)
			}
			if value.LockedBy != scenario.pipeline.LockedBy {
				t.Errorf("Expected locker user %s, but received %s", scenario.pipeline.LockedBy, value.LockedBy)
			}
		})
	}

	t.Run("FindByProjectAndEnvironment_pipelineExists_returnsPipeline", func(t *testing.T) {
		pipeline := repository.Find(domain.PipelineIdentifier{
			Project:     projectOne,
			Environment: environmentOne,
		})
		if pipeline == nil {
			t.Errorf("Expected store to include project, but got nil")
		} else if pipeline.Project != projectOne {
			t.Errorf("Expected store to return pipeline with project %s, but got project %s from store instead", projectOne, pipeline.Project)
		} else if pipeline.Environment != environmentOne {
			t.Errorf("Expected store to return pipeline with environment %s, but got environment %s from store instead", projectOne, pipeline.Environment)
		} else if pipeline.LockedBy != userTwo {
			t.Errorf("Expected store to return pipeline with user %s, but got user %s from store instead", userTwo, pipeline.LockedBy)
		}
	})

	t.Run("FindByProjectAndEnvironment_pipelineNotExists_returnsNil", func(t *testing.T) {
		pipeline := repository.Find(domain.PipelineIdentifier{
			Project:     projectOne,
			Environment: environmentTwo,
		})
		if pipeline != nil {
			t.Errorf("Expected store to not include project, but got %v", pipeline)
		}
	})

	t.Run("FindByProjectAndEnvironment_pipelineNotExists_returnsNil", func(t *testing.T) {
		pipeline := repository.Find(domain.PipelineIdentifier{
			Project:     projectOne,
			Environment: environmentTwo,
		})
		if pipeline != nil {
			t.Errorf("Expected store to not include project, but got %v", pipeline)
		}
	})

	t.Run("FindLockedPipelines_storeHasOnlyOneLockedPipeline_returnsSliceOfOnePipeline", func(t *testing.T) {
		pipelines := repository.FindLockedPipelines()
		if pipelines == nil {
			t.Errorf("Expected store to return slice, but got nil")
		}
		if len(pipelines) != 1 {
			t.Errorf("Expected store to return one pipeline, but got %d", len(pipelines))
		}
	})

	repository = NewPipelineRepository(pipelineKeyCaseSensitive)

	t.Run("FindLockedPipelines_storeHasNoLockedPipelines_returnsEmptySlice", func(t *testing.T) {
		pipelines := repository.FindLockedPipelines()
		if pipelines == nil {
			t.Errorf("Expected store to return slice, but got nil")
		}
		if len(pipelines) != 0 {
			t.Errorf("Expected store to return empty slice, but got %d", len(pipelines))
		}
	})

	repository = NewPipelineRepository(false)
	t.Run("Add_pipelineKeyCaseInSensitive_caseInsensitiveKeyIsAdded", func(t *testing.T) {
		repository.Add(domain.Pipeline{
			PipelineIdentifier: domain.PipelineIdentifier{
				Project:     projectOne,
				Environment: environmentOne,
			},
			PipelineLockedBy: domain.PipelineLockedBy{
				LockedBy: userOne,
			},
		})

		returnedPipeline := repository.Find(domain.PipelineIdentifier{
			Project:     projectTwo,
			Environment: environmentTwo,
		})

		if returnedPipeline == nil {
			t.Errorf("Expected repository to have caseInsensitive keys!")
		}
	})
}
