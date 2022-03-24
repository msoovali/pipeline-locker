package memory

import (
	"strings"

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

func (r *pipelineRepository) Find(request domain.PipelineIdentifier) *domain.Pipeline {
	key := r.getKey(request.Project, request.Environment)
	pipeline, exists := r.store[key]
	if !exists {
		return nil
	}

	return &pipeline
}

func (r *pipelineRepository) Add(pipeline domain.Pipeline) {
	key := r.getKey(pipeline.Project, pipeline.Environment)
	r.store[key] = pipeline
}

func (r *pipelineRepository) FindLockedPipelines() []domain.Pipeline {
	lockedPipelines := make([]domain.Pipeline, 0)
	for _, p := range r.store {
		if p.LockedBy != "" {
			lockedPipelines = append(lockedPipelines, p)
		}
	}

	return lockedPipelines
}

func (r *pipelineRepository) getKey(project, environment string) string {
	if !r.caseSensitiveKey {
		project = strings.ToLower(project)
		environment = strings.ToLower(environment)
	}
	return project + separator + environment
}
