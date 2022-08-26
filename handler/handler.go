package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Status string

const (
	UNCHECKED = "UNCHECKED"
	UP        = "UP"
	DOWN      = "DOWN"
)

// We are using global variable because if we create it inside ReqHandler, a new Map with be created with each GET/POST request,
// Thus the url list updated by POST will not be accessible for GET
var m = map[string]Status{}

type Urls struct {
	Websites []string `json:"websites"`
}

type StatusChecker interface {
	Check(ctx context.Context, name string) (status bool, err error)
}
type httpChecker struct {
}

func (h httpChecker) Check(ctx context.Context, name string) (status bool) {
	resp, err := http.Get(name)
	if err == nil && resp.StatusCode == http.StatusOK {
		return true
	} else if err != nil {
		return false
	}
	return
}

// Request Handler function
func ReqHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		updateMap(w, r, GetMap())
	case "GET":
		getStatus(w, r, GetMap())
	default:
		fmt.Println("Unexpected command")
	}
}

// Function to check and update status of URLs in every minute
func UpdateStatus(m *map[string]Status) {
	for {
		for key := range *m {
			go UpdateStatusUtil(key, m)
		}
		fmt.Println("Status is checked and updated")
		time.Sleep(60 * time.Second)
	}
}

// Utility function to check the status of URL
func UpdateStatusUtil(key string, m *map[string]Status) {
	H := httpChecker{}
	status := H.Check(context.Background(), key)
	if status {
		(*m)[key] = UP
	} else {
		(*m)[key] = DOWN
	}
}

func updateMap(w http.ResponseWriter, r *http.Request, m *map[string]Status) {
	urls := Urls{}
	err := json.NewDecoder(r.Body).Decode(&urls)
	if err != nil {
		log.Println("Unable to decode JSON request body:", err)
	}
	for _, val := range urls.Websites {
		if _, ok := (*m)[val]; !ok {
			(*m)[val] = UNCHECKED
		}
	}
	fmt.Fprint(w, "Map is updated")
}

func getStatus(w http.ResponseWriter, r *http.Request, m *map[string]Status) {
	if len(*m) == 0 {
		fmt.Fprint(w, "Website list is empty, please add URLs using POST method")
		return
	}
	name := r.URL.Query().Get("name")
	if name != "" {
		if _, ok := (*m)[name]; !ok {
			fmt.Fprint(w, "The URL '", name, "' is not present in the watch list, please add it using POST method")
		} else {
			w.Header().Set("Content-Type", "application/json")
			p := map[string]Status{name: (*m)[name]}
			json.NewEncoder(w).Encode(p)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(*m)
}

func GetMap() *map[string]Status {
	return &m
}
