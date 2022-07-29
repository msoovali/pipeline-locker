package service

import (
	"time"

	"github.com/msoovali/pipeline-locker/internal/domain"
)

type pipelineService struct {
	repository       domain.PipelineRepository
	allowOverLocking bool
}

func NewPipelineService(repository domain.PipelineRepository, allowOverlocking bool) *pipelineService {
	return &pipelineService{
		repository:       repository,
		allowOverLocking: allowOverlocking,
	}
}

func (s *pipelineService) IsDeployAllowed(request domain.PipelineIdentifier) (bool, error) {
	if err := request.Validate(); err != nil {
		return false, err
	}
	pipeline, err := s.repository.Find(request)
	if err != nil {
		return false, err
	}

	return pipeline == nil || pipeline.LockedBy == "", nil
}

func (s *pipelineService) Lock(pipeline domain.PipelineLockRequest) error {
	if err := pipeline.Validate(); err != nil {
		return err
	}
	if !s.allowOverLocking {
		existingPipeline, err := s.repository.Find(pipeline.PipelineIdentifier)
		if err != nil {
			return err
		}
		if existingPipeline != nil && existingPipeline.LockedBy != "" {
			return domain.ErrPipelineAlreadyLocked
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
	if err := pipeline.Validate(); err != nil {
		return err
	}
	s.repository.Add(domain.Pipeline{
		PipelineIdentifier: pipeline,
	})

	return nil
}

func (s *pipelineService) GetLockedPipelines() ([]domain.Pipeline, error) {
	return s.repository.FindLockedPipelines()
}
