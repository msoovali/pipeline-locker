package domain

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrProjectEmpty          = errors.New("REQUEST_PROJECT_EMPTY")
	ErrEnvironmentEmpty      = errors.New("REQUEST_ENVIRONMENT_EMPTY")
	ErrLockedByEmpty         = errors.New("REQUEST_LOCKED_BY_EMPTY")
	ErrPipelineAlreadyLocked = errors.New("PIPELINE_ALREADY_LOCKED")
)

type Pipeline struct {
	PipelineIdentifier
	PipelineLockedBy
	PipelineLockedAt
}

type PipelineIdentifier struct {
	Project     string `json:"project" form:"project"`
	Environment string `json:"environment" form:"environment"`
}

type PipelineLockedBy struct {
	LockedBy string `json:"locked_by" form:"locked_by"`
}

type PipelineLockedAt struct {
	LockedAt time.Time `json:"locked_at"`
}

type PipelineLockRequest struct {
	PipelineIdentifier
	PipelineLockedBy
}

func (p *PipelineIdentifier) Validate() error {
	if p.Project == "" {
		return ErrProjectEmpty
	}
	if p.Environment == "" {
		return ErrEnvironmentEmpty
	}

	return nil
}

func (p *PipelineIdentifier) GetKey(caseSensitiveKey bool, separator string) string {
	project := p.Project
	environment := p.Environment
	if !caseSensitiveKey {
		project = strings.ToLower(project)
		environment = strings.ToLower(environment)
	}
	return project + separator + environment
}

func (p *PipelineLockRequest) Validate() error {
	if err := p.PipelineIdentifier.Validate(); err != nil {
		return err
	}
	if p.LockedBy == "" {
		return ErrLockedByEmpty
	}

	return nil
}

type PipelineRepository interface {
	Find(pipeline PipelineIdentifier) (*Pipeline, error)
	Add(pipeline Pipeline) error
	FindLockedPipelines() ([]Pipeline, error)
}

type PipelineService interface {
	IsDeployAllowed(PipelineIdentifier) (bool, error)
	Lock(PipelineLockRequest) error
	Unlock(PipelineIdentifier) error
	GetLockedPipelines() ([]Pipeline, error)
}
