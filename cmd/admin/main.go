package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	// Get the absolute path to the web directory
	webDir, err := filepath.Abs("web")
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("Serving files from: %s\n", webDir)
	
	// Serve static files from the web directory
	fs := http.FileServer(http.Dir(webDir))
	http.Handle("/", fs)
	
	// Specific handler for admin page
	http.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		
		adminPath := filepath.Join(webDir, "admin.html")
		fmt.Printf("Serving admin.html from: %s\n", adminPath)
		http.ServeFile(w, r, adminPath)
	})
	
	// Proxy API calls to Cloud Run to bypass CORS
	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		// Remove /api prefix and forward to Cloud Run
		cloudRunURL := "https://lookie-727276629029.us-central1.run.app" + r.URL.Path[4:]
		
		// Create new request
		req, err := http.NewRequest(r.Method, cloudRunURL, r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		// Copy headers
		for name, values := range r.Header {
			for _, value := range values {
				req.Header.Add(name, value)
			}
		}
		
		// Make request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()
		
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		
		// Copy response headers
		for name, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(name, value)
			}
		}
		
		w.WriteHeader(resp.StatusCode)
		
		// Copy response body
		_, err = io.Copy(w, resp.Body)
		if err != nil {
			fmt.Printf("Error copying response: %v\n", err)
		}
	})
	
	fmt.Println("🎛️  Lookie Admin Panel Server")
	fmt.Println("📊 Dashboard: http://localhost:8081/admin.html")
	fmt.Println("🔧 Admin Panel: http://localhost:8081/admin")
	fmt.Println("🚀 Starting server on port 8081...")
	
	log.Fatal(http.ListenAndServe(":8081", nil))
}