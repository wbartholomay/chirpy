package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const port = "8080"

	cfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	mux := http.NewServeMux()

	fileServerHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))

	mux.Handle("/app/", cfg.middlewareMetricsInc(fileServerHandler))
	mux.HandleFunc("/healthz", ReadinessHandler)
	mux.HandleFunc("/metrics", cfg.MetricsHandler)
	mux.HandleFunc("/reset", cfg.ResetMetricsHandler)

	server := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())

}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {

	return http.HandlerFunc(func (w http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, req)
	})
}


func ReadinessHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)

	w.Write(([]byte)("OK"))
}

func (cfg *apiConfig) MetricsHandler(w http.ResponseWriter, req *http.Request) {
	hitsStr := fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())

	w.Write(([]byte)(hitsStr))
}

func (cfg *apiConfig) ResetMetricsHandler(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)
}
