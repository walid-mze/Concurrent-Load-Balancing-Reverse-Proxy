package cheker

import (
	"ReverseProxy/models"
	"log"
	"net/http"
	"time"
)

func StartHealthCheck(pool *models.ServerPool, interval time.Duration) {
	// Exécuter immédiatement la première vérification
	checkBackends(pool)

	ticker := time.NewTicker(interval)

	go func() {
		for range ticker.C {
			checkBackends(pool)
		}
	}()

}

func checkBackends(pool *models.ServerPool) {

	for _, backend := range pool.Backends {

		alive := isAlive(backend)
		backend.Mux.Lock()
		wasAlive := backend.Alive
		backend.Mux.Unlock()

		pool.SetBackendStatus(backend.URL, alive)
		if alive && !wasAlive {
			log.Printf("Backend %s is UP", backend.URL)
		} else if wasAlive && !alive {
			log.Printf("Backend %s is DOWN", backend.URL)
		}
	}

}
func isAlive(backend *models.Backend) bool {
	client := http.Client{
		Timeout: 2 * time.Second,
	}
	// Check /health for health status
	/*this check is used for slow  backends that may take time to respond ,
	and it it used to test the cancellation of the slow backends requests*/
	healthURL := backend.URL.String() + "/health"
	resp, err := client.Get(healthURL)
	if err != nil {
		// Fallback to root if /health doesn't exist
		resp, err = client.Get(backend.URL.String())
		if err != nil {
			return false
		}
	}
	defer resp.Body.Close()
	if resp.StatusCode < 500 {
		return true
	}
	return false

}
