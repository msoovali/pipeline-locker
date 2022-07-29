package memory

import (
	"github.com/msoovali/pipeline-locker/internal/domain"
)

const separator = ":"

type pipelineRepository struct {
	store            map[string]domain.Pipeline
	caseSensitiveKey bool
}

func NewPipelineRepository(caseSensitiveKey bool) *pipelineRepository {
	return &pipelineRepository{
		store:            make(map[string]domain.Pipeline),
		caseSensitiveKey: caseSensitiveKey,
	}
}

func (r *pipelineRepository) Find(identifier domain.PipelineIdentifier) (*domain.Pipeline, error) {
	key := identifier.GetKey(r.caseSensitiveKey, separator)
	pipeline, exists := r.store[key]
	if !exists {
		return nil, nil
	}

	return &pipeline, nil
}

func (r *pipelineRepository) Add(pipeline domain.Pipeline) error {
	key := pipeline.PipelineIdentifier.GetKey(r.caseSensitiveKey, separator)
	r.store[key] = pipeline

	return nil
}

func (r *pipelineRepository) FindLockedPipelines() ([]domain.Pipeline, error) {
	lockedPipelines := make([]domain.Pipeline, 0)
	for _, p := range r.store {
		if p.LockedBy != "" {
			lockedPipelines = append(lockedPipelines, p)
		}
	}

	return lockedPipelines, nil
}
