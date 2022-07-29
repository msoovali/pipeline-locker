package v7

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/go-redis/redis/v9"
	"github.com/msoovali/pipeline-locker/internal/domain"
)

const separator = ":"

type pipelineRepository struct {
	redisClient      *redis.Client
	caseSensitiveKey bool
}

func NewPipelineRepository(redisClient *redis.Client, caseSensitiveKey bool) *pipelineRepository {
	return &pipelineRepository{
		redisClient:      redisClient,
		caseSensitiveKey: caseSensitiveKey,
	}
}

func (r *pipelineRepository) Find(identifier domain.PipelineIdentifier) (*domain.Pipeline, error) {
	key := identifier.GetKey(r.caseSensitiveKey, separator)

	return r.findByKey(key)
}

func (r *pipelineRepository) findByKey(key string) (*domain.Pipeline, error) {
	value, err := r.redisClient.Get(context.Background(), key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}
	var pipeline domain.Pipeline
	if err = json.Unmarshal([]byte(value), &pipeline); err != nil {
		return nil, err
	}

	return &pipeline, nil
}

func (r *pipelineRepository) Add(pipeline domain.Pipeline) error {
	key := pipeline.PipelineIdentifier.GetKey(r.caseSensitiveKey, separator)
	marshaledPipeline, err := json.Marshal(pipeline)
	if err != nil {
		return err
	}
	err = r.redisClient.Set(context.Background(), key, string(marshaledPipeline), 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *pipelineRepository) FindLockedPipelines() ([]domain.Pipeline, error) {
	keys := make([]string, 0)
	ctx := context.Background()
	iter := r.redisClient.Scan(context.Background(), 0, "*", 0).Iterator()
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}
	lockedPipelines := make([]domain.Pipeline, 0)
	for _, key := range keys {
		p, err := r.findByKey(key)
		if err != nil {
			return nil, err
		}
		if p.LockedBy != "" {
			lockedPipelines = append(lockedPipelines, *p)
		}
	}

	return lockedPipelines, nil
}
