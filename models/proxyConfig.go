package models 




type ProxyConfig struct{
	Port 				int `json:"port"`
	Strategy 			string `json:"strategy"`
	HealthChekerFreq 	string `json:"health_check_frequency"`
}


