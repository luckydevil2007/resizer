package cmd

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
 _  "github.com/lib/pq"
	"github.com/go-chi/chi/v5"
	"github.com/pressly/goose"
	"gopkg.in/yaml.v3"
)

type TransformStruct struct {
	Name   string `json:"name"`
	Rotate int    `json:"rotate"`
	Resize int    `json:"resize"`
}

type DbConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type Config struct {
	DbConfig DbConfig `yaml:"database"`
}

func formatDSN(dbCponfig DbConfig) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		dbCponfig.Host, dbCponfig.Port, dbCponfig.User, dbCponfig.Password, dbCponfig.DBName)
}

func LoadConfig(Path string) (Config, error) {

	var config Config

	content, err := os.ReadFile(Path)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func openBD(dsn string) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		fmt.Println("Failed to open database:", err)
	}
	defer db.Close()

	if err := goose.Up(db, "migrations"); err != nil {
		fmt.Println("Couldn't apply a migration: ", err)
	}

	fmt.Println("Migrations applied successfully!")
}

func Run() {
	var config, _ = LoadConfig("db.yaml")
	dsn := formatDSN(config.DbConfig)
	openBD(dsn)

	r := chi.NewRouter()

	//TODO middleware for autorization
	r.Get("/", mainHandler)
	r.Get("/download/{image}", downloadHandler)

	r.Put("/update/json", resizeHandler)
	r.Post("/upload", uploadHandler)
	r.Delete("/delete/{image}", deleteHandler)

	fmt.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Println("Coulan't establish connection on :8080...")
	}
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello")
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Cannot upload image", http.StatusExpectationFailed)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)

	if os.WriteFile(header.Filename, data, 0644) != nil {
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File uploaded successfully")
}

func resizeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Resize")
	var transform TransformStruct
	err := json.NewDecoder(r.Body).Decode(&transform)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "%s", "Image "+transform.Name+" formatted")
	w.WriteHeader(http.StatusOK)

}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "image")
	_, err := os.Stat(name)
	if err != nil {
		http.Error(w, "File not found", http.StatusBadRequest)
		return
	}

	if os.Remove(name) != nil {
		http.Error(w, "File not found", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Image successfully deleted")
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "image")
	_, err := os.Stat(filename)
	if err != nil {
		http.Error(w, "File not found", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", "application/octet-stream")
	//
	http.ServeFile(w, r, filename)
}
