package domain

import "time"

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
