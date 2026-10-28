package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func main() {
	dbPath := env("DB_PATH", "sportshub.db")
	addr := env("ADDR", ":8080")
	allowed := strings.Split(env("CORS_ORIGINS", "http://localhost:5173"), ",")

	db, err := openDB(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := seed(db); err != nil {
		log.Fatal("seed: ", err)
	}

	s := &Server{db: db}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("GET /api/facets", s.facets)
	mux.HandleFunc("GET /api/products", s.listProducts)
	mux.HandleFunc("GET /api/products/{id}", s.getProduct)
	mux.HandleFunc("GET /api/events", s.listEvents)
	mux.HandleFunc("GET /api/events/{id}", s.getEvent)
	mux.HandleFunc("POST /api/events/{id}/registrations", s.register)
	mux.HandleFunc("POST /api/orders", s.createOrder)

	srv := &http.Server{
		Addr:              addr,
		Handler:           logging(cors(allowed, mux)),
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Printf("SportsHub API listening on %s (db: %s)", addr, dbPath)
	log.Fatal(srv.ListenAndServe())
}

// cors is only needed if the browser calls the API directly; the SvelteKit app
// calls it server-side, so by default this just allows the Vite dev origin.
func cors(origins []string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		for _, o := range origins {
			if strings.TrimSpace(o) == origin || o == "*" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
				break
			}
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) { r.status = code; r.ResponseWriter.WriteHeader(code) }

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{w, 200}
		next.ServeHTTP(rec, r)
		log.Printf("%s %s %d %s", r.Method, r.URL.RequestURI(), rec.status, time.Since(start).Round(time.Microsecond))
	})
}
