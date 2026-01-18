package cheker

import (
	"ReverseProxy/models"
	"log"
	"net/http"
	"time"
)
func StartHealthCheck(pool *models.ServerPool, interval time.Duration ){
    ticker := time.NewTicker(interval)
    
    go func() {
        for range ticker.C {
            checkBackends(pool)
        }
    }()

}

func checkBackends(pool *models.ServerPool){

	for _,backend:=range pool.Backends{

		alive:=isAlive(backend)
		backend.Mux.Lock()
		wasAlive:=backend.Alive
		backend.Mux.Unlock()

		pool.SetBackendStatus(backend.URL,alive)
		if alive && !wasAlive{
			log.Printf("Backend %s is UP", backend.URL)
		}else if wasAlive && !alive{
			log.Printf("Backend %s is DOWN", backend.URL)
		}
	}


	
}
func isAlive(backend *models.Backend) bool{	
	client :=http.Client{
		Timeout: 2*time.Second,
	}
	resp,err:=client.Get(backend.URL.String())
	if err!=nil{
		return false
	}
	defer resp.Body.Close()
	if  resp.StatusCode < 500{
		return true
	}
	return false

}
