package controllers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/luckydevil2007/go-lessons/entities"
	"github.com/luckydevil2007/go-lessons/repositories"
	"github.com/luckydevil2007/go-lessons/usecases"
)

type ImageController struct {
	checkAuth *usecases.AuthUseCase
	image     *usecases.ImageUseCase
	repo      *repositories.Repository
}

type HttpController struct { //->responsible only for REST
	db        *sql.DB
	checkAuth *usecases.AuthUseCase
	//db *sql.DB //хранить не ДБ а интерфейс юзкейса
}

func NewHttpController(db *sql.DB) HttpController {
	return HttpController{db: db}
}

func NewImageController(authUC *usecases.AuthUseCase,
	imageUC *usecases.ImageUseCase,
	repo *repositories.Repository) *ImageController {
	return &ImageController{
		checkAuth: authUC,
		image:     imageUC,
		repo:      repo,
	}
}

func (c *ImageController) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			http.Error(w, "401", http.StatusUnauthorized)
			return
		}

		userID, err := c.checkAuth.Authenticate(r.Context(), username, password)
		if err != nil {
			http.Error(w, "401", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "id", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (c *ImageController) UploadHandler(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Cannot upload image", http.StatusBadRequest)
		return
	}
	defer file.Close()

	userID := r.Context().Value("id").(int)
	var data []byte
	data, err = io.ReadAll(file)

	if err != nil {
		http.Error(w, "Cannot readf image", http.StatusInternalServerError)
		return
	}
	if err := c.image.Upload(r.Context(), header.Filename, data, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (c *ImageController) TransformHandler(w http.ResponseWriter, r *http.Request) {
	var transform entities.ImageTransform
	if err := json.NewDecoder(r.Body).Decode(&transform); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var img entities.Image
	img.ID = transform.ID
	c.image.Transform(r.Context(), &img, transform)
	fmt.Fprintf(w, "%s", "Image "+transform.Name+" formatted")

	//w.WriteHeader(http.StatusOK)
	//json.NewEncoder(w).Encode(map[string]string{"status": "transformed"})
}

func (c *ImageController) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	var image entities.Image
	image.ID, _ = strconv.Atoi(r.FormValue("id"))
	if c.repo.OpenImage(r.Context(), &image) == nil {
		w.WriteHeader(http.StatusNotFound)
	}
	if c.image.Delete(r.Context(), &image) != nil {
		w.WriteHeader(http.StatusOK)
	}
	w.WriteHeader(http.StatusInternalServerError)
}

/*
func (c *ImageController) resizeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Resize")
	var transform entities.ImageTransform
	err := json.NewDecoder(r.Body).Decode(&transform)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	sqlStatement := `SELECT image_path from images WHERE (id) ($1)`

	var path string
	_ = c.db.QueryRow(sqlStatement, r.Context().Value("id")).Scan(&path)
	img = io.Reader.Read()
	image.TransformStruct.Transform(transform, path)
	fmt.Fprintf(w, "%s", "Image "+transform.Name+" formatted")
	w.WriteHeader(http.StatusOK)

}

func (c Controller) downloadHandler(w http.ResponseWriter, r *http.Request) {
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
*/
