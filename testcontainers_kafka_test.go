package testcontainers_go_kafka_2670

import (
	"context"
	"flag"
	"testing"

	"github.com/testcontainers/testcontainers-go/modules/kafka"
)

var imageFlag = flag.String("image", "confluentinc/confluent-local:7.5.0", "image with Kafka")

func TestKafkaContainerStart(t *testing.T) {
	image := *imageFlag
	ctx := context.Background()
	t.Logf("starting container using image: %s", image)
	container, err := kafka.Run(ctx, image)
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
