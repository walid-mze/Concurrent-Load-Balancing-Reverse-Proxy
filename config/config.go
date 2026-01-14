package config

import (
	"encoding/json"
	"log"
	"net/url"
	"os"

	"ReverseProxy/models"
)

func LoadConfig(filename string) (*models.ServerPool, error) {

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var config struct {
		Backends []string `json:"backends"`
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	pool := &models.ServerPool{
		Backends: make([]*models.Backend, 0),
	}

	for _, urlStr := range config.Backends {
		u, err := url.Parse(urlStr)
		if err != nil {
			log.Printf("Error parsing URL %s: %v", urlStr, err)
			continue
		}
		backend := &models.Backend{
			URL:   u,
			Alive: true,
		}
		pool.AddBackend(backend)
	}
	return pool, nil
}
