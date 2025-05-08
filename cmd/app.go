package cmd

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"github.com/luckydevil2007/go-lessons/controllers"
	"github.com/luckydevil2007/go-lessons/repositories"
	"github.com/luckydevil2007/go-lessons/usecases"
	"github.com/pressly/goose"
	"gopkg.in/yaml.v3"
)

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
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
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

func Hash(str, key string) string {

	return str
}

//controllers
//----http
//----другой способ дергания юзкейсов

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

// move to another package
// по сути это драйвер БД
func openBD(dsn string) *sql.DB {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		fmt.Println("Failed to open database:", err)
		return nil
	}

	if err := goose.Up(db, "migrations"); err != nil {
		fmt.Println("Couldn't apply a migration: ", err)
		return nil
	}
	fmt.Println("Migrations applied successfully!")

	return db
}

func Run() {
	var config, _ = LoadConfig("db.yaml")
	dsn := formatDSN(config.DbConfig)
	fmt.Println(dsn)
	db := openBD(dsn)

	r := chi.NewRouter()
	repo := repositories.NewUserRepository(db)
	imageRepo := repositories.NewRepository(db)
	fileStorage := repositories.NewFileStorage(db, "storage")
	authUsecase := usecases.NewAuthUseCase(repo)
	algo := usecases.NewTransformAlgorithm()
	imageUsecase := usecases.NewImageUseCase(imageRepo, fileStorage, algo)

	c := controllers.NewImageController(authUsecase, imageUsecase, imageRepo)

	//TODO middleware for autorization
	r.Use(c.AuthMiddleware)
	r.Get("/", mainHandler)
	//	r.Get("/download/{image}", c.downloadHandler)

	r.Put("/update/json", c.TransformHandler)
	r.Post("/upload", c.UploadHandler)
	r.Delete("/delete/{image}", c.DeleteHandler)

	fmt.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Println("Couldn't establish connection on :8080...")
	}
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello")
}
