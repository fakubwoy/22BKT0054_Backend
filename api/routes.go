package api

import (
	"database/sql"
	"net/http"

	"github.com/fakubwoy/go-file-share/internal/auth"
	"github.com/fakubwoy/go-file-share/internal/config"
	"github.com/fakubwoy/go-file-share/internal/handlers"
	"github.com/fakubwoy/go-file-share/internal/storage"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

func SetupRoutes(db *sql.DB, rdb *redis.Client, cfg *config.Config, storage storage.Storage) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/register", handlers.RegisterHandler(db, cfg)).Methods("POST")
	r.HandleFunc("/login", handlers.LoginHandler(db, cfg)).Methods("POST")

	fileRouter := r.PathPrefix("/files").Subrouter()
	fileRouter.Use(auth.AuthMiddleware(cfg))

	fileRouter.HandleFunc("", handlers.ListFilesHandler(db, rdb)).Methods("GET")
	fileRouter.HandleFunc("", handlers.UploadHandler(db, cfg, storage, rdb)).Methods("POST")
	fileRouter.HandleFunc("/search", handlers.SearchFilesHandler(db, rdb)).Methods("GET")
	fileRouter.HandleFunc("/{id}/share", handlers.ShareFileHandler(db, cfg, storage)).Methods("POST")
	fileRouter.HandleFunc("/{id}", handlers.DeleteFileHandler(db, rdb)).Methods("DELETE")

	r.HandleFunc("/share/{token}", handlers.GetSharedFileHandler(db, storage)).Methods("GET")

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	return r
}
