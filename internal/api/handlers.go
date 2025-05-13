package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"url-shortener/internal/logic"
)

type ApiResponse struct {
	ShortenedUrl string `json:"shortened_url"`
	Err          string `json:"error,omitempty"`
}

func HandleReq(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path[1:]
		if path == "" {
			http.ServeFile(w, r, "web/static/index.html")
			return
		}

		var url string
		err := db.QueryRow("SELECT url FROM urls WHERE short = ?", path).Scan(&url)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		http.Redirect(w, r, url, http.StatusSeeOther)
	}
}

func HandleApi(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
			return
		}

		var requestData struct {
			Url string `json:"url"`
		}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&requestData)
		if err != nil {
			http.Error(w, "Invalid JSON or missing 'url' field", http.StatusBadRequest)
			return
		}

		u := requestData.Url

		if !logic.ValidateUrl(u) {
			response := ApiResponse{
				Err: "Invalid URL",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(response)
			if err != nil {
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
				return
			}
			return
		}

		isReachable, err := logic.IsRealURL(requestData.Url)
		if err != nil || !isReachable {
			response := ApiResponse{
				Err: "Invalid URL",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(response)
			if err != nil {
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
				return
			}
			return
		}

		if !strings.HasPrefix(u, "http://") && !strings.HasPrefix(u, "https://") {
			u = "http://" + u
		}

		short := logic.AddUrlToDb(u, db)

		response := ApiResponse{
			ShortenedUrl: fmt.Sprintf("http://localhost:8080/%s", short),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}
