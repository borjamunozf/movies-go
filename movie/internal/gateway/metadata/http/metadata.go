package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	model "microgomovies/metadata/pkg"
	"microgomovies/movie/internal/gateway"
	"microgomovies/pkg/discovery"
	"net/http"
)

// Gateway defines a movie metadata HTTP gatway
type Gateway struct {
	registry discovery.Register
}

func New(registry discovery.Register) *Gateway {
	return &Gateway{registry: registry}
}

func (g *Gateway) Get(ctx context.Context, id string) (*model.Metadata, error) {

	addrs, err := g.registry.ServiceAddresses(ctx, "metadata")
	if err != nil {
		return nil, err
	}

	url := "http://" + addrs[rand.Intn(len(addrs))] + "/metadata"
	log.Printf("Calling metadata service. Request: GET " + url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("id", id)
	req.URL.RawQuery = values.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, gateway.ErrNotFound
	} else if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("non-2xx response: %v", resp)
	}

	var v *model.Metadata
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}
	return v, nil
}
