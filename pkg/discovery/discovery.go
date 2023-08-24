package discovery

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type Register interface {
	// Register creates a service instance record in the registry
	Register(ctx context.Context, instanceID string, serviceName string, hostPort string) error
	// Deregister removes a service instance record in the registry
	Deregister(ctx context.Context, instanceID string, serviceName string) error
	// ServiceAddresses returns a service instance address list from the registry
	ServiceAddresses(ctx context.Context, serviceID string) ([]string, error)
	HealthCheck(instanceID string, serviceName string) error
}

var ErrNotFound = errors.New("no service address found")

func GenerateInstanceID(serviceName string) string {
	return fmt.Sprintf("%s-%d", serviceName, rand.New(rand.NewSource(time.Now().UnixNano())).Int63())
}
