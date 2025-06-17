package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/wbartholomay/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db *database.Queries
	platform string
}

func main() {
	const port = "8080"
	godotenv.Load()

	dbUrl := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)

	cfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db : dbQueries,
		platform: os.Getenv("PLATFORM"),
	}

	mux := http.NewServeMux()

	fileServerHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))

	mux.Handle("/app/", cfg.middlewareMetricsInc(fileServerHandler))
	mux.HandleFunc("GET /api/healthz", makeHandler(ReadinessHandler))
	mux.HandleFunc("POST /api/chirps", makeHandler(cfg.CreateChirpHandler))
	mux.HandleFunc("GET /api/chirps", makeHandler(cfg.GetAllChirpsHandler))
	mux.HandleFunc("GET /api/chirps/{chirp_id}", makeHandler(cfg.GetChirpByIdHandler))
	mux.HandleFunc("POST /api/users", makeHandler(cfg.CreateUserHandler))
	mux.HandleFunc("GET /admin/metrics", makeHandler(cfg.MetricsHandler))
	mux.HandleFunc("POST /admin/reset", makeHandler(cfg.ResetHandler))

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


func ReadinessHandler(w http.ResponseWriter, req *http.Request) error{
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)

	w.Write(([]byte)("OK"))
	return nil
}

func (cfg *apiConfig) MetricsHandler(w http.ResponseWriter, req *http.Request) error{
	hitsStr := fmt.Sprintf(`<html>
		<body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
		</body>
		</html>`, cfg.fileserverHits.Load())

	w.Write(([]byte)(hitsStr))
	return nil
}

func (cfg *apiConfig) ResetHandler(w http.ResponseWriter, req *http.Request) error{
	if cfg.platform != "dev" {
		return APIError{
			Status: 403,
			Msg: "Forbidden",
		}
	}

	cfg.fileserverHits.Store(0)
	cfg.db.DeleteAllUsers(req.Context())
	return nil
}
