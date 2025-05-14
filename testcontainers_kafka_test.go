package testcontainers_go_kafka_2670

import (
	"context"
	"flag"
	"fmt"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/kafka"
)

var imageFlag = flag.String("image", "confluentinc/confluent-local:7.5.0", "image with Kafka")

func TestKafkaContainerStart(t *testing.T) {
	image := *imageFlag
	ctx := context.Background()
	t.Logf("starting container using image: %s", image)
	container, err := kafka.Run(ctx, image,
		withWaitForPortMapping("9093/tcp", time.Minute, time.Second))
	if err != nil {
		t.Fatalf("container start failed: %s", err)
	}
	defer func() {
		t.Logf("terminating container: %s", container.GetContainerID())
		err := container.Terminate(ctx)
		if err != nil {
			t.Errorf("container termination failed: %s", err)
		}
	}()
	t.Logf("successfully started container: %s", container.GetContainerID())
}

func withWaitForPortMapping(port nat.Port, duration time.Duration, interval time.Duration) testcontainers.CustomizeRequestOption {
	return func(req *testcontainers.GenericContainerRequest) error {
		req.LifecycleHooks = append([]testcontainers.ContainerLifecycleHooks{{
			PostStarts: []testcontainers.ContainerHook{
				func(ctx context.Context, c testcontainers.Container) error {
					return waitForPortMapping(ctx, c, port, duration, interval)
				},
			},
		}}, req.LifecycleHooks...)
		return nil
	}
}

func waitForPortMapping(ctx context.Context, container testcontainers.Container, port nat.Port,
	duration time.Duration, interval time.Duration,
) error {
	ctx, cancel := context.WithTimeout(ctx, duration)
	defer cancel()
	_, err := container.MappedPort(ctx, port)
	for i := 0; err != nil; i++ {
		select {
		case <-ctx.Done():
			return fmt.Errorf("mapped port: retries: %d, port: %s, last err: %w, ctx err: %w", i, port, err, ctx.Err())
		case <-time.After(interval):
			_, err = container.MappedPort(ctx, port)
		}
	}
	return nil
}
