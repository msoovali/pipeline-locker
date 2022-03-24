package service

import "github.com/msoovali/pipeline-locker/internal/domain"

type PipelineService interface {
	IsDeployAllowed(domain.PipelineIdentifier) (bool, error)
	Lock(domain.Pipeline) error
	Unlock(domain.PipelineIdentifier) error
	GetLockedPipelines() []domain.Pipeline
}
