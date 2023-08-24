package consul

import (
	"context"
	"errors"
	"fmt"
	"microgomovies/pkg/discovery"
	"strconv"
	"strings"

	consul "github.com/hashicorp/consul/api"
)

type ConsulRegistry struct {
	client *consul.Client
}

func NewRegistry(addr string) (*ConsulRegistry, error) {
	config := consul.DefaultConfig()
	config.Address = addr
	client, err := consul.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &ConsulRegistry{client: client}, nil
}

func (cr *ConsulRegistry) Register(ctx context.Context, instanceID string, serviceName string, hostPort string) error {
	parts := strings.Split(hostPort, ":")
	if len(parts) != 2 {
		return errors.New("hostPort must be in a form of <host>:<port>, example: localhost:8081")
	}
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}

	return cr.client.Agent().ServiceRegister(&consul.AgentServiceRegistration{
		Address: parts[0],
		ID:      instanceID,
		Name:    serviceName,
		Port:    port,
		Check:   &consul.AgentServiceCheck{CheckID: instanceID, TTL: "5s"},
	})
}

func (cr *ConsulRegistry) Deregister(ctx context.Context, instanceID string, serviceName string) error {
	return cr.client.Agent().ServiceDeregister(instanceID)
}

func (cr *ConsulRegistry) ServiceAddresses(ctx context.Context, serviceName string) ([]string, error) {
	entries, _, err := cr.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, err
	} else if len(entries) == 0 {
		return nil, discovery.ErrNotFound
	}
	var res []string
	for _, e := range entries {
		res = append(res, fmt.Sprintf("%s:%d", e.Service.Address, e.Service.Port))
	}

	return res, nil
}

func (cr *ConsulRegistry) HealthCheck(instanceID string, _ string) error {
	return cr.client.Agent().UpdateTTL(instanceID, "", "pass")
}
