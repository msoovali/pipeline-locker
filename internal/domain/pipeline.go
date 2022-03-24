package domain

type Pipeline struct {
	PipelineIdentifier
	LockedBy string `json:"locked_by" form:"locked_by"`
}

type PipelineIdentifier struct {
	Project     string `json:"project" form:"project"`
	Environment string `json:"environment" form:"environment"`
}
