package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	currentConfig Config
	configMutex   sync.RWMutex
	configPath    string // Set from main
)

func StartWebServer(port string, cfgPath string, initialConfig Config) {
	configPath = cfgPath
	currentConfig = initialConfig

	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", handleDashboard)
	http.HandleFunc("/api/config", handleConfig)
	http.HandleFunc("/api/subscriptions", handleSubscriptions)
	http.HandleFunc("/api/subscriptions/", handleSubscriptionDelete)

	log.Printf("Web Interface started at http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Web server failed: %v", err)
	}
}

func handleDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, htmlContent)
}

func handleConfig(w http.ResponseWriter, r *http.Request) {
	configMutex.Lock()
	defer configMutex.Unlock()

	switch r.Method {
	case "GET":
		json.NewEncoder(w).Encode(currentConfig)
	case "POST":
		var newConfigWrapper struct {
			Global GlobalConfig `json:"global"`
		}
		if err := json.NewDecoder(r.Body).Decode(&newConfigWrapper); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Update only global config part, keep subscriptions
		currentConfig.Global.UpdateInterval = newConfigWrapper.Global.UpdateInterval
		currentConfig.Global.OutputPath = newConfigWrapper.Global.OutputPath
		saveConfig()
		json.NewEncoder(w).Encode(currentConfig)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleSubscriptions(w http.ResponseWriter, r *http.Request) {
	configMutex.Lock()
	defer configMutex.Unlock()

	switch r.Method {
	case "GET":
		json.NewEncoder(w).Encode(currentConfig.Subscriptions)
	case "POST":
		var newSub Subscription
		if err := json.NewDecoder(r.Body).Decode(&newSub); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		
		// Check for duplicate
		found := false
		for i, sub := range currentConfig.Subscriptions {
			if sub.Name == newSub.Name {
				currentConfig.Subscriptions[i] = newSub
				found = true
				break
			}
		}
		if !found {
			currentConfig.Subscriptions = append(currentConfig.Subscriptions, newSub)
		}
		saveConfig()
		json.NewEncoder(w).Encode(currentConfig.Subscriptions)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleSubscriptionDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := r.URL.Path[len("/api/subscriptions/"):]
	if name == "" {
		http.Error(w, "Name required", http.StatusBadRequest)
		return
	}

	configMutex.Lock()
	defer configMutex.Unlock()

	newSubs := []Subscription{}
	for _, sub := range currentConfig.Subscriptions {
		if sub.Name != name {
			newSubs = append(newSubs, sub)
		}
	}
	currentConfig.Subscriptions = newSubs
	saveConfig()
	w.WriteHeader(http.StatusOK)
}

func saveConfig() {
	data, err := yaml.Marshal(&currentConfig)
	if err != nil {
		log.Printf("Failed to marshal config: %v", err)
		return
	}
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		log.Printf("Failed to write config file: %v", err)
	} else {
		log.Println("Configuration saved to disk.")
	}
}
