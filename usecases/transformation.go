package usecases

import (
	"context"
	"errors"

	"github.com/luckydevil2007/go-lessons/entities"
	"github.com/luckydevil2007/go-lessons/repositories"
)

type TransformService interface {
	Transform(transform entities.ImageTransform, data []byte) ([]byte, error)
}

type ImageUseCase struct {
	imageRepo        *repositories.Repository
	fileStorage      *repositories.FileStorage
	transformService TransformService
}

func NewImageUseCase(imageRepo *repositories.Repository,
	fileStorage *repositories.FileStorage,
	transformService TransformService) *ImageUseCase {
	return &ImageUseCase{
		imageRepo:        imageRepo,
		fileStorage:      fileStorage,
		transformService: transformService,
	}
}

func (uc *ImageUseCase) Upload(ctx context.Context, name string, data []byte, ownerID int) error {
	image := &entities.Image{
		Title: name,
		Path:  name,
		Owner: ownerID,
		Data:  data,
	}
	if err := uc.fileStorage.Save(ctx, image); err != nil {
		return err
	}

	return uc.imageRepo.SaveImage(ctx, image)
}

func (uc *ImageUseCase) Delete(ctx context.Context, image *entities.Image) error {
	if err := uc.fileStorage.Delete(ctx, image); err != nil {
		return err
	}

	return uc.imageRepo.DeleteImage(ctx, image.ID)
}

func (uc *ImageUseCase) Transform(ctx context.Context, image *entities.Image, transform entities.ImageTransform) error {
	err := uc.imageRepo.OpenImage(ctx, image)
	if err != nil {
		return err
	}

	ownerID, ok := ctx.Value("id").(int)
	if !ok || image.Owner != ownerID {
		return errors.New("unauthorized")
	}

	data, err := uc.fileStorage.Open(ctx, image.Path)
	if err != nil {
		return err
	}

	transformedData, err := uc.transformService.Transform(transform, data)
	if err != nil {
		return err
	}
	image.Data = transformedData
	image.Title = transform.Name + image.Title
	return uc.fileStorage.Save(ctx, image)
}
