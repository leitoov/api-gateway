package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	// URLs de nuestros microservicios
	aiServiceURL, _ := url.Parse("http://localhost:8000") // ai-service
	authServiceURL, _ := url.Parse("http://localhost:8080") // auth-service

	// Creamos los proxies inversos
	aiProxy := httputil.NewSingleHostReverseProxy(aiServiceURL)
	authProxy := httputil.NewSingleHostReverseProxy(authServiceURL)

	// Configuramos el enrutador
	mux := http.NewServeMux()

	// Enrutamos /api/v1/chat y /api/v1/tools hacia el ai-service
	mux.HandleFunc("/api/v1/chat", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Redirigiendo a AI Service: %s", r.URL.Path)
		// Configuraciones CORS irían aquí
		aiProxy.ServeHTTP(w, r)
	})
	
	mux.HandleFunc("/api/v1/tools", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Redirigiendo a AI Service (Tools): %s", r.URL.Path)
		aiProxy.ServeHTTP(w, r)
	})

	// Enrutamos /api/auth hacia el auth-service
	mux.HandleFunc("/api/auth/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Redirigiendo a Auth Service: %s", r.URL.Path)
		authProxy.ServeHTTP(w, r)
	})

	// Iniciamos el servidor en el puerto 5000
	port := ":5000"
	log.Printf("API Gateway iniciado en http://localhost%s", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
