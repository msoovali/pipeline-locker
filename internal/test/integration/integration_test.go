package integration_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/msoovali/pipeline-locker/internal/domain"
	v6 "github.com/msoovali/pipeline-locker/internal/repository/redis/v6"
	"github.com/msoovali/pipeline-locker/internal/service"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	project     = "area51"
	environment = "production"
	user        = "bob"
)

type redisContainer struct {
	testcontainers.Container
	URI string
}

func setupRedis(ctx context.Context) (*redisContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "redis:6",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForListeningPort("6379/tcp"),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, "6379")
	if err != nil {
		return nil, err
	}

	hostIP, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("redis://%s:%s", hostIP, mappedPort.Port())

	return &redisContainer{Container: container, URI: uri}, nil
}

func flushRedis(ctx context.Context, client redis.Client) error {
	return client.FlushAll(ctx).Err()
}

func TestIntegrationLockUnlock(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	ctx := context.Background()

	redisContainer, err := setupRedis(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer redisContainer.Terminate(ctx)

	options, err := redis.ParseURL(redisContainer.URI)
	if err != nil {
		t.Fatal(err)
	}
	client := redis.NewClient(options)
	defer flushRedis(ctx, *client)

	repository := v6.NewPipelineRepository(client, true)
	service := service.NewPipelineService(repository, false)

	pipeline := getPipelineIdentifierMock()
	pipelineLockRequest := getPipelineLockRequestMock()
	// lock pipeline
	err = service.Lock(pipelineLockRequest)
	if err != nil {
		t.Errorf("Failed to lock pipeline: %v", err)
		return
	}
	// check pipeline is locked
	isAllowed, err := service.IsDeployAllowed(pipeline)
	if err != nil {
		t.Errorf("Failed to get deploy allow status: %v", err)
		return
	}
	if isAllowed {
		t.Errorf("Expected pipeline to be locked, but it is not")
		return
	}
	// lock another pipeline
	pipelineLockRequest.Environment = "dev"
	err = service.Lock(pipelineLockRequest)
	if err != nil {
		t.Errorf("Failed to lock pipeline: %v", err)
		return
	}
	// get locked pipelines
	pipelines, err := service.GetLockedPipelines()
	if err != nil {
		t.Errorf("Failed to get locked pipelines: %v", err)
		return
	}
	if len(pipelines) != 2 {
		t.Errorf("Expected 2 locked pipelines, but got %d", len(pipelines))
		return
	}
	// unlock pipeline
	err = service.Unlock(pipeline)
	if err != nil {
		t.Errorf("Failed to unlock pipeline: %v", err)
		return
	}
	// deploy status is allowed
	isAllowed, err = service.IsDeployAllowed(pipeline)
	if err != nil {
		t.Errorf("Failed to get deploy allow status: %v", err)
		return
	}
	if !isAllowed {
		t.Errorf("Expected pipeline to be unlocked, but it is not")
		return
	}
}

func getPipelineIdentifierMock() domain.PipelineIdentifier {
	return domain.PipelineIdentifier{
		Project:     project,
		Environment: environment,
	}
}

func getPipelineLockRequestMock() domain.PipelineLockRequest {
	return domain.PipelineLockRequest{
		PipelineIdentifier: getPipelineIdentifierMock(),
		PipelineLockedBy: domain.PipelineLockedBy{
			LockedBy: user,
		},
	}
}
