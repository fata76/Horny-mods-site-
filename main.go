package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Developer struct {
	ID   int      `json:"id"`
	Name string   `json:"name"`
	Game string   `json:"game"`
	Mods []string `json:"mods"`
}

var (
	developers []Developer
	nextID     = 1
	mu         sync.Mutex
)

func main() {

	developers = append(developers, Developer{
		ID:   nextID,
		Name: "MorrowindModder",
		Game: "TES III: Morrowind",
		Mods: []string{"Graphic Overhaul", "Rebirth"},
		
	})
	nextID++

	http.HandleFunc("/api/developers", developersHandler)
    
	http.Handle("/", http.FileServer(http.Dir(".")))

	
	fmt.Println("Сервер запущен на http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
}

func developersHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		mu.Lock()
		json.NewEncoder(w).Encode(developers)
		mu.Unlock()

	case http.MethodPost:
		var dev Developer
		err := json.NewDecoder(r.Body).Decode(&dev)
		if err != nil {
			http.Error(w, "Неверный формат данных", http.StatusBadRequest)
			return
		}

		mu.Lock()
		dev.ID = nextID
		nextID++
		developers = append(developers, dev)
		mu.Unlock()

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(dev)

	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}
