package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)

		data := map[string]any{
			"id":   1,
			"name": "John Doe",
			"subData": map[string]any{
				"subKey1": "subValue1",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Printf("Error encoding JSON: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})

	port := 8080
	fmt.Printf("Mock API server starting on http://localhost:%d\n", port)
	fmt.Println("Press Ctrl+C to stop")

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
