package usecases

import (
	"bytes"
	"fmt"
	"image"

	"github.com/disintegration/imaging"
	"github.com/luckydevil2007/go-lessons/entities"
)

type TransformAlgorithm struct {
}

func NewTransformAlgorithm() *TransformAlgorithm {
	return &TransformAlgorithm{}
}

func (algo *TransformAlgorithm) Transform(t entities.ImageTransform, data []byte) ([]byte, error) {
	reader := bytes.NewReader(data)
	img, _, err := image.Decode(reader)
	//img, _, err := imaging.Decode(reader)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	width := img.Bounds().Size().X * t.Resize / 100.0
	height := img.Bounds().Size().Y * t.Resize / 100.0

	resized := imaging.Resize(img, width, height, imaging.Lanczos)
	rotated := imaging.Rotate(resized, float64(t.Rotate), nil)

	var w bytes.Buffer //.Writer

	err = imaging.Encode(&w, rotated, imaging.PNG)
	if err != nil {
		return nil, fmt.Errorf("failed to encode image: %v", err)
	}
	return w.Bytes(), nil
}

/*
func (t TransformStruct) Transform(path string) *image.RGBA { //возвращать массив байт
	//лучше найти библиотеку которая работает не с диском а с абстракцией файла или io.Reader
	imaging.Decode()
	img, err := imgio.Open(path)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	w := img.Bounds().Size().X * t.Resize / 100.0
	h := img.Bounds().Size().Y * t.Resize / 100.0

	resized := transform.Resize(img, w, h, transform.Linear)
	rotated := transform.Rotate(resized, float64(t.Rotate), nil)
	return rotated
}
*/
