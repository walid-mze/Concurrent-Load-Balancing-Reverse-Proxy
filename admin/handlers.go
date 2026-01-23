package admin

import (
	"encoding/json"
	"net/http"
	"net/url"

	"ReverseProxy/models"
)
type AdminAPI struct {
    ServerPool *models.ServerPool
}

func (a *AdminAPI) StatusHandler(w http.ResponseWriter, r *http.Request){
	if r.Method!=http.MethodGet{
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return 
	}
	a.ServerPool.Mux.RLock()
	defer a.ServerPool.Mux.RUnlock()
	active:=0
	for _,backend:=range a.ServerPool.Backends{
		backend.Mux.RLock()
		if backend.Alive{
			active++
		}
		backend.Mux.RUnlock()
	}
	response:=map[string]any{
		"total_backends":len(a.ServerPool.Backends),
		"active_backends":active,
		"backends":a.ServerPool.Backends,
	}
	w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}


func (a *AdminAPI) AddBackendHandler(w http.ResponseWriter, r *http.Request){
	if r.Method!=http.MethodPost{
		http.Error(w,"Method not allowed",http.StatusMethodNotAllowed)
		return 
	}

	var request struct{
		URL string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err!=nil{
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
	}
	u,err:=url.Parse(request.URL)
	if err != nil {
        http.Error(w, "Invalid URL", http.StatusBadRequest)
        return
    }
	backend:=&models.Backend{
		URL: u,
		Alive: true,
		CurrentConnections: 0,
	}
	for bc:=range a.ServerPool.Backends{
		if a.ServerPool.Backends[bc].URL.String()==backend.URL.String(){
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
		"message": "Backend already exists",
        "url":u.String(),
	})
		return 
		}
	}
	a.ServerPool.AddBackend(backend)

    w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Backend added successfully",
        "url":u.String(),
	})
}
func (a *AdminAPI )DeleteBackendHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method!=http.MethodDelete{
		http.Error(w,"Method not allowed ",http.StatusMethodNotAllowed)
		return 
	}
	var request struct{
		URL string `json:"url"`
	}
	err:=json.NewDecoder(r.Body).Decode(&request)
	if err!=nil{
		http.Error(w,"Inavlid json format",http.StatusBadRequest)
		return 
	}
	url,err:=url.Parse(request.URL)
		if err!=nil{
		http.Error(w,"Inavlid url",http.StatusBadRequest)
		return 
	}
	backend:=&models.Backend{
		URL: url,
	}
	
	
	for i:=range a.ServerPool.Backends{
		if a.ServerPool.Backends[i].URL.String()==backend.URL.String(){

			a.ServerPool.DeleteBackend(backend)
			
			w.Header().Set("Content-Type","application/json")
			json.NewEncoder(w).Encode(map[string]string{
			"message": "Backend deleted successfully",
			"url": url.String(),
			})

		return 
		}
	}


	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Backend does not exists",
        "url": url.String(),
	})
}



