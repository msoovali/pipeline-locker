package service

import (
	"errors"
	"time"

	"github.com/msoovali/pipeline-locker/internal/domain"
	"github.com/msoovali/pipeline-locker/internal/repository"
)

var (
	ProjectEmptyError          = errors.New("REQUEST_PROJECT_EMPTY")
	EnvironmentEmptyError      = errors.New("REQUEST_ENVIRONMENT_EMPTY")
	LockedByEmptyError         = errors.New("REQUEST_LOCKED_BY_EMPTY")
	PipelineAlreadyLockedError = errors.New("PIPELINE_ALREADY_LOCKED")
)

type pipelineService struct {
	repository       repository.PipelineRepository
	allowOverLocking bool
}

func NewPipelineService(repository repository.PipelineRepository, allowOverlocking bool) *pipelineService {
	return &pipelineService{
		repository:       repository,
		allowOverLocking: allowOverlocking,
	}
}

func (s *pipelineService) IsDeployAllowed(request domain.PipelineIdentifier) (bool, error) {
	if err := verifyRequest(request); err != nil {
		return false, err
	}
	pipeline := s.repository.Find(request)

	return pipeline == nil || pipeline.LockedBy == "", nil
}

func (s *pipelineService) Lock(pipeline domain.PipelineLockRequest) error {
	if err := verifyLockRequest(pipeline); err != nil {
		return err
	}
	if !s.allowOverLocking {
		existingPipeline := s.repository.Find(pipeline.PipelineIdentifier)
		if existingPipeline != nil && existingPipeline.LockedBy != "" {
			return PipelineAlreadyLockedError
		}
	}
	s.repository.Add(domain.Pipeline{
		PipelineIdentifier: pipeline.PipelineIdentifier,
		PipelineLockedBy:   pipeline.PipelineLockedBy,
		PipelineLockedAt: domain.PipelineLockedAt{
			LockedAt: time.Now(),
		},
	})

	return nil
}

func (s *pipelineService) Unlock(pipeline domain.PipelineIdentifier) error {
	if err := verifyRequest(pipeline); err != nil {
		return err
	}
	s.repository.Add(domain.Pipeline{
		PipelineIdentifier: pipeline,
	})

	return nil
}

func (s *pipelineService) GetLockedPipelines() []domain.Pipeline {
	return s.repository.FindLockedPipelines()
}

func verifyRequest(pipeline domain.PipelineIdentifier) error {
	if pipeline.Project == "" {
		return ProjectEmptyError
	}
	if pipeline.Environment == "" {
		return EnvironmentEmptyError
	}

	return nil
}

func verifyLockRequest(pipeline domain.PipelineLockRequest) error {
	if err := verifyRequest(pipeline.PipelineIdentifier); err != nil {
		return err
	}
	if pipeline.LockedBy == "" {
		return LockedByEmptyError
	}

	return nil
}
