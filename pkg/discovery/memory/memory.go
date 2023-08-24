package memory

import (
	"context"
	"errors"
	"microgomovies/pkg/discovery"
	"sync"
	"time"
)

type serviceName string
type instanceID string

type MemoryRegistry struct {
	sync.RWMutex
	serviceAddrs map[serviceName]map[instanceID]*serviceInstance
}

type serviceInstance struct {
	hostPort   string
	lastActive time.Time
}

func NewRegistry() *MemoryRegistry {
	return &MemoryRegistry{serviceAddrs: map[serviceName]map[instanceID]*serviceInstance{}}
}

func (mr *MemoryRegistry) Register(ctx context.Context, instanceID instanceID, serviceName serviceName, hostPort string) error {
	mr.Lock()
	defer mr.Unlock()
	if _, ok := mr.serviceAddrs[serviceName]; !ok {
		mr.serviceAddrs[serviceName] = map[string]*serviceInstance{}
	}

	mr.serviceAddrs[serviceName][instanceID] = &serviceInstance{hostPort: hostPort, lastActive: time.Now()}
	return nil
}

func (mr *MemoryRegistry) Deregister(ctx context.Context, instanceID instanceID, serviceName serviceName) error {
	mr.Lock()
	defer mr.Unlock()
	if _, ok := mr.serviceAddrs[serviceName]; !ok {
		return nil
	}
	delete(mr.serviceAddrs[serviceName], instanceID)
	return nil
}

func (mr *MemoryRegistry) ServiceAddresses(ctx context.Context, serviceName serviceName) ([]string, error) {
	mr.Lock()
	defer mr.Unlock()
	if _, ok := mr.serviceAddrs[serviceName]; !ok {
		return nil, discovery.ErrNotFound
	}
	var res []string
	for _, i := range mr.serviceAddrs[serviceName] {
		if i.lastActive.Before(time.Now().Add(-5 * time.Second)) {
			continue
		}
		res = append(res, i.hostPort)
	}
	return res, nil
}

func (mr *MemoryRegistry) HealthCheck(instanceID instanceID, serviceName serviceName) error {
	mr.Lock()
	defer mr.Unlock()
	if _, ok := mr.serviceAddrs[serviceName]; !ok {
		return errors.New("service is not registered yet")
	}
	if _, ok := mr.serviceAddrs[serviceName][instanceID]; !ok {
		return errors.New("service instance is not registered yet")
	}

	mr.serviceAddrs[serviceName][instanceID].lastActive = time.Now()
	return nil
}
