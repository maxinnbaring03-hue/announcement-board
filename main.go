package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type Announcement struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Message   string `json:"message"`
	AuthorID  string `json:"author_id"`
	CreatedAt string `json:"created_at"`
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type PostRequest struct {
	Title   string `json:"title"`
	Message string `json:"message"`
	Author  string `json:"author_id"`
}

var db *sql.DB

func main() {
	var err error

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("Error: DATABASE_URL environment variable is not set")
	}

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Database connection error:", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Cannot reach Supabase:", err)
	}
	fmt.Println("Successfully connected to Supabase cloud database!")

	http.HandleFunc("/api/announcements", announcementsHandler)
	http.HandleFunc("/api/register", createAdminHandler)
	http.HandleFunc("/api/login", loginHandler)

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func announcementsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// Added DELETE to the allowed methods
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// 1. READING DATA (GET)
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		rows, err := db.Query("SELECT id, title, message, author_id, created_at FROM announcements ORDER BY created_at DESC")
		if err != nil {
			http.Error(w, "Failed to query database", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var list []Announcement
		for rows.Next() {
			var a Announcement
			if err := rows.Scan(&a.ID, &a.Title, &a.Message, &a.AuthorID, &a.CreatedAt); err != nil {
				return
			}
			list = append(list, a)
		}
		if list == nil {
			list = []Announcement{}
		}
		json.NewEncoder(w).Encode(list)
		return
	}

	// 2. WRITING DATA (POST)
	if r.Method == http.MethodPost {
		var req PostRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		_, err := db.Exec("INSERT INTO announcements (title, message, author_id) VALUES ($1, $2, $3)", req.Title, req.Message, req.Author)
		if err != nil {
			http.Error(w, "Failed to save announcement", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Announcement posted successfully"})
		return
	}

	// 3. DELETING DATA (DELETE) - THE NEW ENGINE
	if r.Method == http.MethodDelete {
		// Grab the ID from the URL (e.g., /api/announcements?id=123)
		postID := r.URL.Query().Get("id")
		if postID == "" {
			http.Error(w, "Missing post ID", http.StatusBadRequest)
			return
		}

		// Tell Supabase to wipe it
		_, err := db.Exec("DELETE FROM announcements WHERE id = $1", postID)
		if err != nil {
			http.Error(w, "Failed to delete post", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Post deleted successfully"})
		return
	}
}

func createAdminHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Encryption failed", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("INSERT INTO accounts (username, password_hash, role) VALUES ($1, $2, $3)", req.Username, string(hashedPassword), "admin")
	if err != nil {
		http.Error(w, "Failed to create account", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Admin created"})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	var storedHash string
	err := db.QueryRow("SELECT password_hash FROM accounts WHERE username=$1", req.Username).Scan(&storedHash)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(req.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"})
}
