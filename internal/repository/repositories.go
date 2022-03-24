package repository

import "github.com/msoovali/pipeline-locker/internal/domain"

type PipelineRepository interface {
	Find(pipeline domain.PipelineIdentifier) *domain.Pipeline
	Add(pipeline domain.Pipeline)
	FindLockedPipelines() []domain.Pipeline
}
