package main

import (
	"fmt"
	"log"
	"net/http"
)

type apiConfig struct {
	fileserverhits int
}

var apiCfg = apiConfig{}

func main() {
	const filepathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()
	mux.Handle("/app/*", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("/healthz", handlerReadiness)
	mux.HandleFunc("/metrics", handlerMetrics)
	mux.HandleFunc("/reset", apiCfg.handlerResetMetric)

	corsMux := middlewareCors(middlewareLog(mux))

	httpServer := http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Starting development server at - http://localhost:8080")
	log.Fatal(httpServer.ListenAndServe())
}

func handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %v", apiCfg.fileserverhits)))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverhits++
		next.ServeHTTP(w, r)
	})
}

func middlewareLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
